package controllers

import (
	"net/http"

	"github.com/alinowrouzii/service-health-check/models"
)

func (cfg *Config) GetWarningsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.UserResponse)

	warnings, err := models.GetWarning(cfg.DB, user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"warnings": warnings})
}
