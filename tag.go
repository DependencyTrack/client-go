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

type TaggedProjectListResponseItem struct {
	UUID    uuid.UUID `json:"uuid,omitempty"`
	Name    string    `json:"name,omitempty"`
	Version string    `json:"version,omitempty"`
}

type TaggedPolicyListResponseItem struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
	Name string    `json:"name,omitempty"`
}

func (ts TagService) Create(ctx context.Context, names []string) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.13.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodPut, "/api/v1/tag", withBody(names))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) GetAll(ctx context.Context, po PageOptions, so SortOptions) (p Page[TagListResponseItem], err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

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
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodDelete, "/api/v1/tag", withBody(names))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) TagProjects(ctx context.Context, tag string, projects []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tag/%s/project", tag), withBody(projects))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) UntagProjects(ctx context.Context, tag string, projects []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/tag/%s/project", tag), withBody(projects))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) GetProjects(ctx context.Context, tag string, po PageOptions, so SortOptions) (p Page[TaggedProjectListResponseItem], err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tag/%s/project", tag), withPageOptions(po), withSortOptions(so))
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

func (ts TagService) TagPolicies(ctx context.Context, tag string, policies []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tag/%s/policy", tag), withBody(policies))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) UntagPolicies(ctx context.Context, tag string, policies []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/tag/%s/policy", tag), withBody(policies))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) GetPolicies(ctx context.Context, tag string, po PageOptions, so SortOptions) (p Page[TaggedPolicyListResponseItem], err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tag/%s/policy", tag), withPageOptions(po), withSortOptions(so))
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

func (ts TagService) TagNotificationRules(ctx context.Context, tag string, rules []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tag/%s/notificationRule", tag), withBody(rules))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) UntagNotificationRules(ctx context.Context, tag string, rules []uuid.UUID) (err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/tag/%s/notificationRule", tag), withBody(rules))
	if err != nil {
		return
	}
	_, err = ts.client.doRequest(req, nil)
	return
}

func (ts TagService) GetNotificationRules(ctx context.Context, tag string, po PageOptions, so SortOptions) (p Page[TaggedPolicyListResponseItem], err error) {
	err = ts.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tag/%s/notificationRule", tag), withPageOptions(po), withSortOptions(so))
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

func (ts TagService) GetTagsForPolicy(ctx context.Context, policy uuid.UUID, po PageOptions, so SortOptions) (p Page[Tag], err error) {
	var req *http.Request
	if ts.client.isServerVersionAtLeast("4.12.0") {
		req, err = ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tag/policy/%s", policy), withPageOptions(po), withSortOptions(so))
	} else {
		err = ts.client.assertServerVersionAtLeast("4.6.0")
		if err != nil {
			return
		}

		req, err = ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tag/%s", policy), withPageOptions(po), withSortOptions(so))
	}
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
