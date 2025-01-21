package api

import (
	"encoding/json"
	"net/http"
	"rate-limiter/server/models"
	"rate-limiter/server/redis"
)

type Handler struct {
	redisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{redisClient: redisClient}
}

func (h *Handler) SetRateLimitRule(w http.ResponseWriter, r *http.Request) {
	var limit models.Ratelimit
	if err := json.NewDecoder(r.Body).Decode(&limit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := h.redisClient.SetRateLimitRule(ctx, limit.APIPath, &limit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CheckRateMLimit(w http.ResponseWriter, r *http.Request) {

}
