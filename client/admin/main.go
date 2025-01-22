package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"rate-limiter/server/models"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "Rate limiter server URL")
	flag.Parse()

	limit := models.Ratelimit{
		APIPath:           "/api/example.com",
		RequestsPerSecond: 10,
		RequestsPerMinute: 100,
		RequestsPerHour:   1000,
		RequestsPerDay:    100000,
	}
	data, err := json.Marshal(limit)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(*serverURL+"/api/limit", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Rate limit set: %d\n", resp.StatusCode)
}
