package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/ereminiu/pvz/internal/background"
	"github.com/ereminiu/pvz/internal/cache"
	"github.com/ereminiu/pvz/internal/config/application"
	auconfig "github.com/ereminiu/pvz/internal/config/auditlog"
	kafkaconfig "github.com/ereminiu/pvz/internal/config/kafka"
	"github.com/ereminiu/pvz/internal/config/postgres"
	redisconf "github.com/ereminiu/pvz/internal/config/redis"
	"github.com/ereminiu/pvz/internal/pkg/auditlog"
	dbadapter "github.com/ereminiu/pvz/internal/pkg/auditlog/db_adapter"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/filter"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/stdinadapter"
	rediscache "github.com/ereminiu/pvz/internal/pkg/cache"
	"github.com/ereminiu/pvz/internal/pkg/db"
	pkglog "github.com/ereminiu/pvz/internal/pkg/logger"
	kafkaoutbox "github.com/ereminiu/pvz/internal/pkg/outbox"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	rep "github.com/ereminiu/pvz/internal/repository"
	"github.com/ereminiu/pvz/internal/tracing"
	"github.com/ereminiu/pvz/internal/transport/grpc"
	grpcorder "github.com/ereminiu/pvz/internal/transport/grpc/handler/order"
	grpcpvz "github.com/ereminiu/pvz/internal/transport/grpc/handler/pvz"
	grpcuser "github.com/ereminiu/pvz/internal/transport/grpc/handler/user"
	"github.com/ereminiu/pvz/internal/transport/rest"
	hnd "github.com/ereminiu/pvz/internal/transport/rest/handler"
	"github.com/ereminiu/pvz/internal/transport/rest/monitoring"
	restRouter "github.com/ereminiu/pvz/internal/transport/rest/router"
	uc "github.com/ereminiu/pvz/internal/usecases"
)

func getConnection(ctx context.Context, config postgres.Config) (*db.Database, error) {
	URL := config.URL()
	database, err := db.New(ctx, URL)
	if err != nil {
		return nil, err
	}

	if err = database.GetPool().Ping(ctx); err != nil {
		return nil, err
	}

	return database, nil
}

func getRedisConnection(ctx context.Context, config redisconf.Config) (*redis.Client, error) {
	address := config.GetAddress()
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       config.DB,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func getAuditLog(auditlogConfig auconfig.Config, repos *rep.Repository, outbox *kafkaoutbox.Outbox) *auditlog.AuditLog {
	var auditlogger *auditlog.AuditLog

	if auditlogConfig.Adapter == "default" {
		auditlogger = auditlog.New(stdinadapter.New(), filter.NewEmpty(), outbox)
	} else if auditlogConfig.Adapter == "db" {
		adapter := dbadapter.New(repos)

		auditlogger = auditlog.New(adapter, filter.NewAction(auditlogConfig.Filter), outbox)
	} else {
		panic("AAAAAAAAAAAA")
	}

	return auditlogger
}

func loadDBConfig() postgres.Config {
	config, err := postgres.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadRedisConfig() redisconf.Config {
	config, err := redisconf.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadAppConfig() application.Config {
	config, err := application.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadAuditLogConfig() auconfig.Config {
	config, err := auconfig.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadKafkaConfig() kafkaconfig.Config {
	config, err := kafkaconfig.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func initLogger() {
	logger, err := pkglog.NewLogger("local")
	if err != nil {
		log.Fatalln(err)
	}

	slog.SetDefault(logger)
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initLogger()

	dbConfig := loadDBConfig()
	appConfig := loadAppConfig()
	auditLogConfig := loadAuditLogConfig()
	redisConfig := loadRedisConfig()
	kafkaConfig := loadKafkaConfig()

	conn, err := getConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return
	}

	defer conn.GetPool().Close()

	redisClient, err := getRedisConnection(ctx, redisConfig)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return
	}
	pkgCache := rediscache.New(redisClient)

	txManager := txmanager.New(conn)
	repository := rep.New(txManager)

	cache := cache.New(pkgCache)

	usecases := uc.New(repository, cache)

	outobx, err := kafkaoutbox.New(repository, kafkaConfig.GetAddress())
	if err != nil {
		slog.Error("error during init outbox", slog.Any("err", err))
		return
	}
	audit := getAuditLog(auditLogConfig, repository, outobx)

	h := hnd.New(ctx, usecases, audit)

	router := restRouter.New(h)
	restServer := rest.New(rest.Deps{
		Config: appConfig,
		Router: router,
	})

	grpcServer := grpc.New(grpc.Deps{
		OrderHandler: grpcorder.New(usecases),
		UserHandler:  grpcuser.New(usecases),
		PVZHandler:   grpcpvz.New(usecases),
	})

	monServer := monitoring.New(appConfig.MonitoringPort)

	backgrounder := background.New(usecases, cache)

	_, closer := tracing.InitTracer(tracing.TracerConfig{
		ServiceName: appConfig.ServiceName,
		Host:        appConfig.Host,
		Port:        appConfig.TracerPort,
	})
	defer closer.Close()

	go backgrounder.Run(ctx, appConfig.CacheUpateTimeout, func(ctx context.Context) error {
		return backgrounder.FillHistory(ctx)
	})

	go backgrounder.Run(ctx, appConfig.CacheUpateTimeout, func(ctx context.Context) error {
		return backgrounder.FillRefunds(ctx)
	})

	go audit.Run(ctx)

	go func() {
		slog.InfoContext(ctx, "Rest-server is started")
		if err = restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error during rest-server running: ", slog.Any("err", err))
		}
	}()

	go func() {
		slog.InfoContext(ctx, "gRPC-server is started")
		if err = grpcServer.ListenAndServe(appConfig.GRPCPort); err != nil {
			slog.Error("error during gRPC-server running: ", slog.Any("err", err))
		}
	}()

	go func() {
		slog.InfoContext(ctx, "Monitoring-server is started")
		if err = monServer.ListenAndServe(); err != nil {
			slog.Error("error during monitoring-server running: ", slog.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	audit.Stop()

	slog.InfoContext(ctx, "Server is shutting down")

	if err = restServer.Shutdown(ctx); err != nil {
		slog.Error("error during rest-server shutting down: ", slog.Any("err", err))
	}

	grpcServer.Stop()

	if err = monServer.Shutdown(ctx); err != nil {
		slog.Error("error during monitoring-server shutting down: ", slog.Any("err", err))
	}

	conn.GetPool().Close()

	slog.InfoContext(ctx, "Server is closed")
}
