package main

import (
	"log"
	"net/http"

	"github.com/alinowrouzii/service-health-check/models"
	"github.com/alinowrouzii/service-health-check/routers"
	"github.com/alinowrouzii/service-health-check/token"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	jwt    *token.JWTMaker
}

func (a *App) connectDB(user, password, dbName string) {

}

func (a *App) Initialize(user, password, dbname, secretKey string) {
	dbConn, err := models.InitModels()
	if err != nil {
		log.Fatal(err)
	}
	a.DB = dbConn

	a.jwt, err = token.NewJWTMaker(secretKey, a.DB)
	if err != nil {
		log.Fatal(err)
	}

	router := routers.InitRouter(a.DB, a.jwt)
	a.Router = router
}

func (a *App) Run(addr string) {
	log.Println("starting on", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
