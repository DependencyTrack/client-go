package dtrack

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type About struct {
	UUID        uuid.UUID      `json:"uuid"`
	SystemUUID  uuid.UUID      `json:"systemUuid"`
	Application string         `json:"application"`
	Version     string         `json:"version"`
	Timestamp   string         `json:"timestamp"`
	Framework   AboutFramework `json:"framework"`
}

type AboutFramework struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Timestamp string    `json:"timestamp"`
}

type HealthCheck struct {
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type Health struct {
	Status string        `json:"status"`
	Checks []HealthCheck `json:"checks"`
}

type AboutService struct {
	client *Client
}

func (as AboutService) Get(ctx context.Context) (a About, err error) {
	req, err := as.client.newRequest(ctx, http.MethodGet, "/api/version", withoutAuth())
	if err != nil {
		return
	}

	_, err = as.client.doRequest(req, &a)
	return
}

func (as AboutService) Health(ctx context.Context) (h Health, err error) {
	req, err := as.client.newRequest(ctx, http.MethodGet, "/health", withoutAuth())
	if err != nil {
		return
	}

	_, err = as.client.doRequest(req, &h)
	return
}
