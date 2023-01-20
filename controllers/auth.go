package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alinowrouzii/service-health-check/models"
	"github.com/go-playground/validator/v10"
)

func (cfg *Config) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var userToLogin models.UserToLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userToLogin); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(userToLogin); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := models.GetUserByEmailPassword(cfg.DB, userToLogin.Email, userToLogin.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "no credential was found")
		return
	}

	tokenPayload, jwt, err := cfg.JWT.CreateToken(user.Email, 1000*time.Second)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"payload": tokenPayload, "token": jwt})
}

func (cfg *Config) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var userToRegister models.UserToRegister
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userToRegister); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(userToRegister); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := &models.User{
		Email:          userToRegister.Email,
		HashedPassword: userToRegister.Password,
	}
	err := user.CreateUser(cfg.DB)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"createdUser": userResponse})
}

func (cfg *Config) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*models.UserResponse)
	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"user": user})
}
