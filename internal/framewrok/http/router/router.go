package router

import (
	"fmt"
	"net/http"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/adapter/controller"
	"github.com/murasame29/hackathon-util/internal/framewrok/http/middleware"
)

type Router struct {
	mux *http.ServeMux

	handler *controller.Handler
}

func NewRoute(handler *controller.Handler) http.Handler {
	router := &Router{
		mux:     http.NewServeMux(),
		handler: handler,
	}

	router.common()

	{
		router.DiscordOps()
	}

	return router.mux
}

func (r *Router) common() {
	// health check
	r.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("health ok! environment: %s", config.Config.Application.Env)))
	})

	// version check
	r.mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(config.Config.Application.Version))
	})
}

func (r *Router) DiscordOps() {
	// discord channel control
	r.mux.Handle("POST /discord/channel", middleware.BuildChain(
		http.HandlerFunc(r.handler.Channel),
		middleware.LoggerInContext,
		middleware.AccessLog,
	))
	// discord role control
	r.mux.Handle("POST /discord/role", middleware.BuildChain(
		http.HandlerFunc(r.handler.Role),
		middleware.LoggerInContext,
		middleware.AccessLog,
	))
	// sync
	r.mux.Handle("POST /discord/sync", middleware.BuildChain(
		http.HandlerFunc(r.handler.Sync),
		middleware.LoggerInContext,
		middleware.AccessLog,
	))
}
