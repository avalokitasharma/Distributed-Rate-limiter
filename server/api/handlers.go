package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rate-limiter/server/models"
	"rate-limiter/server/redis"
	"time"
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
	var checkReq models.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	limit, err := h.redisClient.GetRateLimitRule(ctx, checkReq.APIPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check each time window
	windows := map[string]struct {
		limit    int
		duration time.Duration
	}{
		"second": {limit.RequestsPerSecond, time.Second},
		"minute": {limit.RequestsPerMinute, time.Minute},
		"hour":   {limit.RequestsPerHour, time.Hour},
		"day":    {limit.RequestsPerDay, 24 * time.Hour},
	}
	for window, config := range windows {
		key := fmt.Sprintf("%s:%s:%s", checkReq.APIPath, checkReq.ClientID, window)
		allowed, err := h.redisClient.IncrementAndCheck(ctx, key, config.limit, config.duration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allowed {
			json.NewEncoder(w).Encode(models.CheckResponse{
				Allowed: false,
				Reason:  fmt.Sprintf("Rate limit exceeded for %s window", window),
			})
			return
		}
	}
	json.NewEncoder(w).Encode(models.CheckResponse{Allowed: true})

}
