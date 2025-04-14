package pvz

import (
	"context"
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
	GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
	GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
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

func (h *Handler) GetRefunds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		Page    int               `json:"page"`
		Limit   int               `json:"limit"`
		OrderBy string            `json:"order_by"`
		Pattern map[string]string `json:"pattern"`
	}

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))
		return
	}

	logEvent := models.Log{
		Action:    "/refund-list",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	if !orderByIsValid(input.OrderBy) {
		http.Error(w, "Invalid order_by", http.StatusBadRequest)

		logEvent.Description = "error during order refund"
		logEvent.Error = "invalid order_by"

		return
	}

	if !validate.CheckPattern(input.Pattern) {
		http.Error(w, "invalid search pattern", http.StatusBadRequest)

		logEvent.Description = "error during order refund"
		logEvent.Error = "invalid search pattern"

		return
	}

	refunds, err := h.uc.GetRefunds(ctx, input.Page, input.Limit, input.OrderBy, input.Pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during getting refunds", slog.Any("err", err))

		logEvent.Description = "error during order refund"
		logEvent.Error = err.Error()

		return
	}

	response, err := sonic.Marshal(refunds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during json unmarshaling", slog.Any("err", err))

		logEvent.Description = "error during order refund"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)

		logEvent.Description = "error during order refund"
		logEvent.Error = myerrors.ErrDuringWritingResponse.Error()
	}

	logEvent.Description = "success response"
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		Page    int               `json:"page"`
		Limit   int               `json:"limit"`
		OrderBy string            `json:"order_by"`
		Pattern map[string]string `json:"pattern"`
	}

	logEvent := models.Log{
		Action:    "/history",
		Timestamp: time.Now(),
	}

	defer func() {
		h.audit.Send(logEvent)
	}()

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during decoding", slog.Any("err", err))

		logEvent.Description = "error during decoding"
		logEvent.Error = err.Error()

		return
	}

	if !orderByIsValid(input.OrderBy) {
		http.Error(w, "Invalid order_by", http.StatusBadRequest)

		logEvent.Description = "error during validation"
		logEvent.Error = "invalid order_by"

		return
	}

	if !validate.CheckPattern(input.Pattern) {
		http.Error(w, "invalid search pattern", http.StatusBadRequest)

		logEvent.Description = "error during validation"
		logEvent.Error = "invalid search pattern"

		return
	}

	refunds, err := h.uc.GetHistory(ctx, input.Page, input.Limit, input.OrderBy, input.Pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during getting history", slog.Any("err", err))

		logEvent.Description = "error during getting history"
		logEvent.Error = err.Error()

		return
	}

	response, err := sonic.Marshal(refunds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("error during json unmarshaling", slog.Any("err", err))

		logEvent.Description = "error during json unmarshaling"
		logEvent.Error = err.Error()

		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}

	logEvent.Description = "success response"
}
