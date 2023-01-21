package routers

import (
	"fmt"
	"net/http"

	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/middleware"
	"github.com/gorilla/mux"
)

func InitLinkRouter(r *mux.Router, cfg *controllers.Config) {
	fmt.Println("Initialize link route...")

	r.PathPrefix("/link").Subrouter().HandleFunc("/", middleware.TokenMiddleware(http.HandlerFunc(cfg.CreateLinkHandler), cfg.JWT)).Methods("POST")
	r.PathPrefix("/link").Subrouter().HandleFunc("/", middleware.TokenMiddleware(http.HandlerFunc(cfg.GetLinksHandler), cfg.JWT)).Methods("GET")
	r.PathPrefix("/link").Subrouter().HandleFunc("/getStat", middleware.TokenMiddleware(http.HandlerFunc(cfg.GetLinksStatHandler), cfg.JWT)).Methods("GET")
	r.PathPrefix("/link").Subrouter().HandleFunc("/getStat/{id}", middleware.TokenMiddleware(http.HandlerFunc(cfg.GetLinksStatByIDHandler), cfg.JWT)).Methods("GET")
}
