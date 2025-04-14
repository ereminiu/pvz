package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	totalRequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request counter",
		Help: "total amount of requests",
	})

	totalErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total error counter",
		Help: "total amount of requests",
	})

	totalOkResponseCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "OK counter",
		Help: "total amount of ok response",
	})

	responseTimeSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "response_time_summary",
		Help: "Summary of response times",
	})

	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request duration seconds",
		Help:    "Request duration in seconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	}, []string{"request_name"})

	errorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "error counter",
		Help: "amount of error occured during request running",
	}, []string{"request_name"})

	okCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ok response counter",
		Help: "amount of ok reponse",
	}, []string{"request_name"})

	rpsCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rps",
		Help: "amount of requests",
	}, []string{"request_name"})

	cacheCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cached request counter",
		Help: "amount of cached requests",
	}, []string{"request_name"})
)

func SetResponseTimeSummary(resposeTime float64) {
	responseTimeSummary.Observe(resposeTime)
}

func SetRequestDuration(requestName string, duration float64) {
	requestDuration.WithLabelValues(requestName).Observe(duration)
}

func SetTotalErrorCounter() {
	totalErrorCounter.Inc()
}

func SetTotalRequestCounter() {
	totalRequestCounter.Inc()
}

func SetTotalOkCounter() {
	totalOkResponseCounter.Inc()
}

func SetCacheCounter(requestName string) {
	cacheCounter.WithLabelValues(requestName).Inc()
}

func SetErrorCounter(requestName string) {
	errorCounter.WithLabelValues(requestName).Inc()
}

func SetOkCounter(requestName string) {
	okCounter.WithLabelValues(requestName).Inc()
}

func SetRPSCounter(requestName string) {
	rpsCounter.WithLabelValues(requestName).Inc()
}

func init() {
	prometheus.MustRegister(
		totalRequestCounter,
		totalErrorCounter,
		totalOkResponseCounter,
		responseTimeSummary,
		requestDuration,
		errorCounter,
		okCounter,
		rpsCounter,
		cacheCounter,
	)
}
