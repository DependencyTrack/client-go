package dtrack

import (
	"context"
	"net/http"
)

type HealthCheck struct {
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type Health struct {
	Status string        `json:"status"`
	Checks []HealthCheck `json:"checks"`
}

type HealthService struct {
	client *Client
}

func (hs HealthService) Get(ctx context.Context) (h Health, err error) {
	req, err := hs.client.newRequest(ctx, http.MethodGet, "health", withoutAuth())
	if err != nil {
		return
	}

	_, err = hs.client.doRequest(req, &h)
	return
}
