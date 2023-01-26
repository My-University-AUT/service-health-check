// main.go

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	AppDbUsername string `env:"APP_DB_USERNAME,file"`
	AppDbPassword string `env:"APP_DB_PASSWORD,file"`
	AppDbName     string `env:"APP_DB_NAME,file"`
}

func main() {
	log.Println("here")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	runnerIntervalMs, err := strconv.Atoi(os.Getenv("RUNNER_INTERVAL_MS"))

	a := App{}

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("SECRET_KEY"),
		runnerIntervalMs)

	a.Run(":8010")
}
