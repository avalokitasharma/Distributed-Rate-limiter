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

	checkReq := models.CheckRequest{
		APIPath:  "/api/example",
		ClientID: "test-client",
	}

	data, err := json.Marshal(checkReq)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(*serverURL+"/api/check", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var checkResp models.CheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request allowed: %v\nReason: %s\n", checkResp.Allowed, checkResp.Reason)
}
