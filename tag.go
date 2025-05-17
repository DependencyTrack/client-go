package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Tag struct {
	Name string `json:"name"`
}

type TagService struct {
	client *Client
}

type TagListResponseItem struct {
	Name                  string `json:"name,omitempty"`
	ProjectCount          int64  `json:"projectCount,omitempty"`
	PolicyCount           int64  `json:"policyCount,omitempty"`
	NotificationRuleCount int64  `json:"notificationRuleCount,omitempty"`
}

func (ts TagService) Create(ctx context.Context, names []string) (err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPut, "/api/v1/tag", withBody(names))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) GetAll(ctx context.Context, po PageOptions, so SortOptions) (p Page[TagListResponseItem], err error) {
	req, err := ts.client.newRequest(ctx, http.MethodGet, "/api/v1/tag", withPageOptions(po), withSortOptions(so))
	if err != nil {
		return
	}

	res, err := ts.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (ts TagService) Delete(ctx context.Context, names []string) (err error) {
	req, err := ts.client.newRequest(ctx, http.MethodDelete, "/api/v1/tag", withBody(names))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) TagProjects(ctx context.Context, tag string, projects []uuid.UUID) (err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/%s/project", tag), withBody(Map(projects, uuid.UUID.String)))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}
