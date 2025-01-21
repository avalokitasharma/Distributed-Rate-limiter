package api

import "github.com/gorilla/mux"

func SetupRoutes(handler *Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/limit", handler.SetRateLimitRule).Methods("POST")
	r.HandleFunc("/api/heck", handler.CheckRateMLimit).Methods("POST")
	return r
}
