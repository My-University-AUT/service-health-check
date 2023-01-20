package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alinowrouzii/service-health-check/token"
	"gorm.io/gorm"
)

type Config struct {
	DB  *gorm.DB
	JWT *token.JWTMaker
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// fmt.Println(payload)
	response, _ := json.Marshal(payload)
	// fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
