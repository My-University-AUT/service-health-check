package routers

import (
	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/token"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitRouter(r *mux.Router, db *gorm.DB, jwt *token.JWTMaker) {
	cfg := &controllers.Config{
		DB:  db,
		JWT: jwt,
	}

	InitAuthRouter(r, cfg)
	InitLinkRouter(r, cfg)
}
