package rest

import (
	"fmt"
	"net/http"

	"github.com/ereminiu/pvz/internal/config/application"
)

type Deps struct {
	Config application.Config
	Router *http.ServeMux
}

func New(deps Deps) *http.Server {
	return &http.Server{
		Handler:      deps.Router,
		Addr:         fmt.Sprintf(":%d", deps.Config.RestPort),
		WriteTimeout: deps.Config.WriteTimeout,
		ReadTimeout:  deps.Config.ReadTimeout,
	}
}
