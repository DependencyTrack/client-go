package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type LDAPService struct {
	client *Client
}

type LdapUser struct {
	Username          string       `json:"username,omitempty"`
	DistinguishedName string       `json:"dn,omitempty"`
	Teams             []Team       `json:"teams,omitempty"`
	Email             string       `json:"email,omitempty"`
	Permissions       []Permission `json:"permissions,omitempty"`
}

type MappedLdapGroupRequest struct {
	Team              string `json:"team"`
	DistinguishedName string `json:"group"`
}

type MappedLdapGroup struct {
	DistinguishedName string    `json:"dn,omitempty"`
	UUID              uuid.UUID `json:"uuid"`
}

func (s LDAPService) AddMapping(ctx context.Context, mapping MappedLdapGroupRequest) (g MappedLdapGroup, err error) {
	req, err := s.client.newRequest(ctx, http.MethodPut, "/api/v1/ldap/mapping", withBody(mapping))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &g)
	return
}

func (s LDAPService) RemoveMapping(ctx context.Context, mappingId uuid.UUID) (err error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/ldap/mapping/%s", mappingId.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}

func (s LDAPService) GetAllAccessibleGroups(ctx context.Context, po PageOptions) (gs Page[string], err error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/api/v1/ldap/groups", withPageOptions(po))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &gs)
	return
}

func (s LDAPService) GetTeamMappings(ctx context.Context, teamUUID uuid.UUID) (gs []MappedLdapGroup, err error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/ldap/team/%s", teamUUID.String()))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &gs)
	return
}

func (s LDAPService) GetUsers(ctx context.Context, po PageOptions) (us Page[LdapUser], err error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/api/v1/user/ldap", withPageOptions(po))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &us)
	return
}

func (s LDAPService) CreateUser(ctx context.Context, user LdapUser) (userOut LdapUser, err error) {
	req, err := s.client.newRequest(ctx, http.MethodPut, "/api/v1/user/ldap", withBody(user))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, &userOut)
	return
}

func (s LDAPService) DeleteUser(ctx context.Context, user LdapUser) (err error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "/api/v1/user/ldap", withBody(user))
	if err != nil {
		return
	}

	_, err = s.client.doRequest(req, nil)
	return
}
