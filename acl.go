package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Mapping struct {
	Team    uuid.UUID `json:"team"`
	Project uuid.UUID `json:"project"`
}

type AccessManagementService struct {
	client *Client
}

func (am AccessManagementService) Map(ctx context.Context, team uuid.UUID, project uuid.UUID) (m Mapping, err error) {
	mapping := Mapping{Team: team, Project: project}
	req, err := am.client.newRequest(ctx, http.MethodPut, "/api/v1/acl/mapping", withBody(mapping))
	if err != nil {
		return
	}
	_, err = am.client.doRequest(req, &m)
	return
}

func (am AccessManagementService) Delete(ctx context.Context, team string, project string) (err error) {
	req, err := am.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/acl/mapping/team/%s/project/%s", team, project))
	if err != nil {
		return
	}

	_, err = am.client.doRequest(req, nil)
	return
}

func (am AccessManagementService) GetProjects(ctx context.Context, team string, excludeInactive, onlyRoot bool) (p []string, err error) {
	excludeInactiveStr := "false"
	onlyRootStr := "false"
	if excludeInactive {
		excludeInactiveStr = "true"
	}
	if onlyRoot {
		onlyRootStr = "true"
	}
	params := map[string]string{
		"excludeInactive": excludeInactiveStr,
		"onlyRoot":        onlyRootStr,
	}

	req, err := am.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/acl/mapping/team/%s", team), withParams(params))
	if err != nil {
		return
	}

	_, err = am.client.doRequest(req, &p)
	return
}
