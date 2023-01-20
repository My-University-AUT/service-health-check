package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alinowrouzii/service-health-check/models"
	"github.com/go-playground/validator/v10"
)

func (cfg *Config) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.UserResponse)

	var link models.Link
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&link); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(link); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	link.UserID = user.ID

	err := link.CreateLink(cfg.DB)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"link": link})
}

func (cfg *Config) GetLinksHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.UserResponse)

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"links": user.Links})
}
