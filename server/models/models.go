package models

type Ratelimit struct {
	APIPath           string `json:"api_path"`
	RequestsPerSecond int    `json:"requests_per_second"`
	RequestsPerMinute int    `json:"requests_per_minute"`
	RequestsPerHour   int    `json:"requests_per_hour"`
	RequestsPerDay    int    `json:"requests_per_day"`
}

type CheckRequest struct {
	APIPath  string `json:"api_path"`
	ClientID string `json:"client_id"`
}

type CheckResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}
