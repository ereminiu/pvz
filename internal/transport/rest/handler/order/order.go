package order

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	sdecoder "github.com/bytedance/sonic/decoder"
	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/auditlog"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

type usecases interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	RemoveOrder(ctx context.Context, id int) error
}

type Handler struct {
	uc    usecases
	audit *auditlog.AuditLog
}

func New(uc usecases, audit *auditlog.AuditLog) *Handler {
	return &Handler{
		uc:    uc,
		audit: audit,
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input *entities.Order

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))
		return
	}

	logEvent := models.Log{
		UserID:    input.UserID,
		Action:    "/add",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	if err := h.uc.AddOrder(ctx, input); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, myerrors.ErrOrderAlreadyCreated) || errors.Is(err, myerrors.ErrInvalidOrderPackingType) {
			code = http.StatusBadRequest
		}

		http.Error(w, err.Error(), code)
		slog.Error("error during adding order", slog.Any("err", err))

		logEvent.Description = "error during adding order"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "text/raw")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("order has been added")); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		ID int `json:"order_id"`
	}

	logEvent := models.Log{
		Action:    "/remove",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during deconding", slog.Any("err", err))
		return
	}

	if err := h.uc.RemoveOrder(ctx, input.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during deleting order", slog.Any("err", err))

		logEvent.Description = "error during deleting order"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "text/raw")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("order has been removed")); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}
