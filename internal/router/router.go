package router

import (
	"fmt"
	"net/http"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/handler"
	"github.com/murasame29/hackathon-util/internal/router/middleware"
)

type Router struct {
	mux *http.ServeMux

	handler *handler.Handler
}

func NewRoute(handler *handler.Handler) http.Handler {
	router := &Router{
		mux:     http.NewServeMux(),
		handler: handler,
	}

	router.common()

	{
		router.DiscordOps()
		router.DiscordBreakoutRoom()
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

	// discord bot INTERACTIONS ENDPOINT URL
	r.mux.HandleFunc("POST /interactions", func(w http.ResponseWriter, r *http.Request) {
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

func (r *Router) DiscordBreakoutRoom() {
	// health check
	r.mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("health ok! environment: %s", config.Config.Application.Env)))
	})
	// discord role control
	r.mux.Handle("POST /discord/role", middleware.BuildChain(
		http.HandlerFunc(r.handler.Role),
		middleware.LoggerInContext,
		middleware.AccessLog,
	))
}
