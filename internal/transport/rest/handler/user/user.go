package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	sdecoder "github.com/bytedance/sonic/decoder"
	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/auditlog"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	"github.com/ereminiu/pvz/internal/transport/rest/handler/validate"
)

type usecases interface {
	RefundOrders(ctx context.Context, userID int, orderIDs []int) error
	ReturnOrders(ctx context.Context, userID int, orderIDs []int) error
	GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error)
}

type Handler struct {
	uc usecases

	audit *auditlog.AuditLog
}

func New(uc usecases, audit *auditlog.AuditLog) *Handler {
	return &Handler{
		uc:    uc,
		audit: audit,
	}
}

func (h *Handler) RefundOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}
	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))
		return
	}

	logEvent := models.Log{
		UserID:    input.UserID,
		Action:    "/refund",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	if err := h.uc.RefundOrders(ctx, input.UserID, input.OrderIDs); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, myerrors.ErrInvalidOrderInput) {
			code = http.StatusBadRequest
		}

		http.Error(w, err.Error(), code)
		slog.Error("error during refunding orders", slog.Any("err", err))

		logEvent.Description = "error during refunding orders"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "text/raw")
	w.WriteHeader(http.StatusResetContent)
	if _, err := w.Write([]byte("orders have been refunded")); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}

func (h *Handler) ReturnOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))
		return
	}

	logEvent := models.Log{
		UserID:    input.UserID,
		Action:    "/return",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	if err := h.uc.ReturnOrders(ctx, input.UserID, input.OrderIDs); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, myerrors.ErrInvalidOrderInput) {
			code = http.StatusBadRequest
		}

		http.Error(w, err.Error(), code)
		slog.Error("error during returning orders", slog.Any("err", err))

		logEvent.Description = "error during returning orders"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "text/raw")
	w.WriteHeader(http.StatusResetContent)
	if _, err := w.Write([]byte("orders have been returned")); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		UserID  int               `json:"user_id"`
		LastN   int               `json:"last_n"`
		Located bool              `json:"located"`
		Pattern map[string]string `json:"pattern"`
	}

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))
		return
	}

	logEvent := models.Log{
		UserID:    input.UserID,
		Action:    "/list",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	if !validate.CheckPattern(input.Pattern) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pages, err := h.uc.GetList(ctx, input.UserID, input.LastN, input.Located, input.Pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during getting list", slog.Any("err", err))

		logEvent.Description = "error during getting list"
		logEvent.Error = err.Error()

		return
	}

	response, err := sonic.Marshal(pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during json unmarshaling", slog.Any("err", err))
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}
