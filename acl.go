package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type AclService struct {
	client *Client
}

type AclMappingRequest struct {
	Team    uuid.UUID `json:"team"`
	Project uuid.UUID `json:"project"`
}

func (as AclService) AddProjectMapping(ctx context.Context, mapping AclMappingRequest) (err error) {
	req, err := as.client.newRequest(ctx, http.MethodPut, "/api/v1/acl/mapping", withBody(mapping))
	if err != nil {
		return
	}
	_, err = as.client.doRequest(req, nil)
	return
}

func (as AclService) RemoveProjectMapping(ctx context.Context, team, project uuid.UUID) (err error) {
	req, err := as.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/acl/mapping/team/%s/project/%s", team, project))
	if err != nil {
		return
	}
	_, err = as.client.doRequest(req, nil)
	return
}

func (as AclService) GetAllProjects(ctx context.Context, team uuid.UUID, po PageOptions) (p Page[Project], err error) {
	req, err := as.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/acl/team/%s", team), withPageOptions(po))
	if err != nil {
		return
	}
	res, err := as.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}
