package routers

import (
	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/token"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, jwt *token.JWTMaker) *mux.Router {
	router := mux.NewRouter()
	cfg := &controllers.Config{
		DB:  db,
		JWT: jwt,
	}

	InitAuthRouter(router, cfg)
	InitLinkRouter(router, cfg)

	return router
}
