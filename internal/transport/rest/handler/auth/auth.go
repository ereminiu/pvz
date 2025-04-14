package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	sdecoder "github.com/bytedance/sonic/decoder"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/auditlog"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

type usecases interface {
	SignIn(ctx context.Context, username, password string) (string, error)
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

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := sdecoder.NewStreamDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("error during json decoding", slog.Any("err", err))
		return
	}

	token, err := h.uc.SignIn(ctx, input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		slog.Error("error during user creation", slog.Any("err", err))
		return
	}

	h.audit.Send(models.Log{
		Action:      "/sign-in",
		Description: fmt.Sprintf("user %s signed-in", input.Username),
		Timestamp:   time.Now(),
	})

	w.Header().Set("Content-type", "text/raw")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(token)); err != nil {
		http.Error(w, myerrors.ErrDuringWritingResponse.Error(), http.StatusInternalServerError)
	}
}
