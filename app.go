package main

import (
	"database/sql"
	"fmt"
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

func (a *App) dropAndCreateDatabase(connectionString, dbName string) {

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE final_db")
	if err != nil {
		panic(err)
	}
	createDBStmt := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARACTER SET = 'utf8mb4'", dbName)
	_, err = db.Exec(createDBStmt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database created successfully")
}

func (a *App) connectDB(user, password, dbName string) {
	dbConn, err := models.InitModels()
	if err != nil {
		log.Fatal(err)
	}
	a.DB = dbConn
}

func (a *App) Initialize(user, password, dbname, secretKey string) {
	a.connectDB(user, password, dbname)

	var err error
	a.jwt, err = token.NewJWTMaker(secretKey, a.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	routers.InitRouter(a.Router, a.DB, a.jwt)
}

func (a *App) Run(addr string) {
	log.Println("starting on", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
