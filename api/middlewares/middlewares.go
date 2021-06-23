package middlewares

import (
	"errors"
	"net/http"

	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next(res, req)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := auth.TokenValid(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		next(res, req)
	}
}
