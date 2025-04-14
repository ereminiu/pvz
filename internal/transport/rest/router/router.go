package router

import (
	"net/http"

	"github.com/ereminiu/pvz/internal/transport/rest/handler"
	mid "github.com/ereminiu/pvz/internal/transport/rest/middleware"
)

func New(h *handler.Handler) *http.ServeMux {
	apiRouter := http.NewServeMux()

	var (
		order = h.OrderHandler
		user  = h.UserHandler
		pvz   = h.PVZHandler
		auth  = h.AuthHandler
	)

	midChain := mid.Chain(mid.Auth, mid.Log)

	apiRouter.HandleFunc("POST /add", midChain(order.Add))
	apiRouter.HandleFunc("DELETE /remove", midChain(order.Remove))

	apiRouter.HandleFunc("POST /refund", midChain(user.RefundOrders))
	apiRouter.HandleFunc("POST /return", midChain(user.ReturnOrders))
	apiRouter.HandleFunc("POST /list", midChain(user.GetList))

	apiRouter.HandleFunc("POST /refund-list", midChain(pvz.GetRefunds))
	apiRouter.HandleFunc("POST /history", midChain(pvz.GetHistory))

	apiRouter.HandleFunc("POST /sign-in", mid.Log(auth.SignIn))

	r := http.NewServeMux()
	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}
