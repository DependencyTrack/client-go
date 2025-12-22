package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

type OIDCService struct {
	client *Client
}

type OIDCGroup struct {
	Name string    `json:"name,omitempty"`
	UUID uuid.UUID `json:"uuid,omitempty"`
}

type OIDCMappingRequest struct {
	Team  uuid.UUID `json:"team"`
	Group uuid.UUID `json:"group"`
}

type OIDCMapping struct {
	Group OIDCGroup `json:"group"`
	UUID  uuid.UUID `json:"uuid"`
}

type OIDCUser struct {
	Username          string       `json:"username"`
	SubjectIdentifier string       `json:"subjectIdentifier"`
	Email             string       `json:"email"`
	Teams             []Team       `json:"teams"`
	Permissions       []Permission `json:"permissions"`
}

type OIDCTokens struct {
	ID     string `json:"idToken"`
	Access string `json:"accessToken,omitempty"`
}

func (s OIDCService) Available(ctx context.Context) (available bool, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, "api/v1/oidc/available", withAcceptContentType("text/plain"))
	if err != nil {
		return
	}

	var value string

	_, err = s.client.doRequest(req, &value)
	if err != nil {
		return
	}
	available, err = strconv.ParseBool(value)
	return
}

func (s OIDCService) GetAllGroups(ctx context.Context) (groups []OIDCGroup, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, "api/v1/oidc/group")
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &groups)
	return

}

func (s OIDCService) CreateGroup(ctx context.Context, name string) (g OIDCGroup, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodPut, "api/v1/oidc/group", withBody(OIDCGroup{Name: name}))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &g)
	return
}
func (s OIDCService) UpdateGroup(ctx context.Context, group OIDCGroup) (g OIDCGroup, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "api/v1/oidc/group", withBody(group))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &g)
	return
}

func (s OIDCService) DeleteGroup(ctx context.Context, groupUUID uuid.UUID) (err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/oidc/group/%s", groupUUID.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}

func (s OIDCService) GetAllTeamsOf(ctx context.Context, group OIDCGroup) (teams []Team, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v1/oidc/group/%s/team", group.UUID.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &teams)
	return
}

func (s OIDCService) AddTeamMapping(ctx context.Context, mapping OIDCMappingRequest) (m OIDCMapping, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodPut, "api/v1/oidc/mapping", withBody(mapping))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &m)
	return
}

func (s OIDCService) RemoveTeamMapping(ctx context.Context, mappingID uuid.UUID) (err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/oidc/mapping/%s", mappingID.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}

func (s OIDCService) RemoveTeamMapping2(ctx context.Context, groupID, teamID uuid.UUID) (err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/oidc/group/%s/team/%s/mapping", groupID.String(), teamID.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}

func (s OIDCService) GetAllUsers(ctx context.Context) (p Page[OIDCUser], err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, "api/v1/user/oidc")
	if err != nil {
		return
	}

	res, err := s.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (s OIDCService) CreateUser(ctx context.Context, userReq OIDCUser) (userRes OIDCUser, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodPut, "api/v1/user/oidc", withBody(userReq))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &userRes)
	return
}

func (s OIDCService) DeleteUser(ctx context.Context, user OIDCUser) (err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete, "api/v1/user/oidc", withBody(user))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}

func (s OIDCService) Login(ctx context.Context, tokens OIDCTokens) (token string, err error) {
	err = s.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	body := url.Values{}
	body.Set("idToken", tokens.ID)
	body.Set("accessToken", tokens.Access)

	req, err := s.client.newRequest(ctx, http.MethodPost, "api/v1/user/oidc/login", withBody(body))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &token)
	return
}
