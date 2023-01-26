package routers

import (
	"fmt"
	"net/http"

	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/middleware"
	"github.com/gorilla/mux"
)

func InitWarningRouter(r *mux.Router, cfg *controllers.Config) {
	fmt.Println("Initialize link route...")
	r.PathPrefix("/warning").Subrouter().HandleFunc("/", middleware.TokenMiddleware(http.HandlerFunc(cfg.GetWarningsHandler), cfg.JWT)).Methods("GET")
}
