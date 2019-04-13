package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func RegisterRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/opentv", openTv)
	router.HandleFunc("/closetv", closeTv)
	router.Use(loggingMiddleware)
	return router
}
