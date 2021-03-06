package server

import (
	"covid19-update-service/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func toUInt(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err

	}
	return uint(u64), nil
}

func (ws *Covid19UpdateWebServer) defaultNotFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeHTTPResponse(model.NewError("Resource not found."), http.StatusNotFound, w, r)
	})
}

func (ws *Covid19UpdateWebServer) createNotAllowedHandler(r *mux.Router) http.HandlerFunc {
	allowedMethods := getAllMethodsForRouter(r)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (ws *Covid19UpdateWebServer) optionHandler(r *mux.Router) http.HandlerFunc {
	allowedMethods := getAllMethodsForRouter(r)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allowed", strings.Join(allowedMethods, ", "))
		writeHTTPResponse(nil, http.StatusOK, w, r)
	})
}

/*
if r.Method == http.MethodOptions {
		w.Header().Set("Allowed", strings.Join(c.allowedMethods, ", "))
	}
*/
