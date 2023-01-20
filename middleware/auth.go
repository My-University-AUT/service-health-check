package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alinowrouzii/service-health-check/controllers"
	"github.com/alinowrouzii/service-health-check/models"
	"github.com/alinowrouzii/service-health-check/token"
)

type middleware func(http.HandlerFunc, *token.JWTMaker) http.HandlerFunc

func ChainMiddleware(f http.HandlerFunc, jwt *token.JWTMaker, m ...middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}
	currentMiddleWare := m[0]
	return currentMiddleWare(ChainMiddleware(f, jwt, m[1:cap(m)]...), jwt)
}

func TokenMiddleware(next http.HandlerFunc, jwt *token.JWTMaker) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		payload, err := jwt.VerifyToken(authHeader[1])

		if err != nil {
			controllers.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// fmt.Println("payload inside middleware", payload)
		user := &models.User{
			Email: payload.Email,
		}
		user, err = models.GetUserByEmail(jwt.DB, payload.Email)
		if err != nil {
			controllers.RespondWithError(w, http.StatusUnauthorized, "Internal server error")
			return
		}

		userResponse := &models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Links: user.Links,
		}

		rcopy := r.WithContext(context.WithValue(r.Context(), "user", userResponse))

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, rcopy)
	})
}
