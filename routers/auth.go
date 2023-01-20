package routers

import (
	"fmt"
	"net/http"

	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/middleware"
	"github.com/gorilla/mux"
)

func InitAuthRouter(r *mux.Router, cfg *controllers.Config) {
	fmt.Println("Initialize auth route...")

	r.PathPrefix("/auth").Subrouter().HandleFunc("/login", cfg.LoginHandler).Methods("POST")
	r.PathPrefix("/auth").Subrouter().HandleFunc("/register", cfg.RegisterHandler).Methods("POST")
	r.PathPrefix("/auth").Subrouter().HandleFunc("/", middleware.TokenMiddleware(http.HandlerFunc(cfg.GetUserHandler), cfg.JWT)).Methods("GET")

}
