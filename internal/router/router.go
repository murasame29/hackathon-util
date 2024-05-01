package router

import (
	"fmt"
	"net/http"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/handler"
)

type Router struct {
	mux *http.ServeMux
}

func NewRoute() http.Handler {
	router := &Router{mux: http.NewServeMux()}

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
	// create channel
	r.mux.HandleFunc("POST /discord/channel", handler.CreateChannel)
	// delete channel
	r.mux.HandleFunc("DELETE /discord/channel", handler.DeleteChannel)
	// create role
	r.mux.HandleFunc("POST /discord/role", handler.CreateRole)
	// delete role
	r.mux.HandleFunc("DELETE /discord/channel", handler.DeleteRole)
	// sync
	r.mux.HandleFunc("POST /discord/sync", handler.Sync)
}
