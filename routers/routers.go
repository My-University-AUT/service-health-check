package routers

import (
	"database/sql"

	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/token"
	"github.com/gorilla/mux"
)

func InitRouter(r *mux.Router, db *sql.DB, jwt *token.JWTMaker) {
	_ = &controllers.Config{
		DB:  db,
		JWT: jwt,
	}

}
