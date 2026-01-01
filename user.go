package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type UserService struct {
	client *Client
}

type ManagedUser struct {
	Username            string       `json:"username"`
	LastPasswordChange  int          `json:"lastPasswordChange"`
	Fullname            string       `json:"fullname,omitempty"`
	Email               string       `json:"email,omitempty"`
	Suspended           bool         `json:"suspended,omitempty"`
	ForcePasswordChange bool         `json:"forcePasswordChange,omitempty"`
	NonExpiryPassword   bool         `json:"nonExpiryPassword,omitempty"`
	Teams               []Team       `json:"teams,omitempty"`
	Permissions         []Permission `json:"permissions,omitempty"`
	NewPassword         string       `json:"newPassword,omitempty"`
	ConfirmPassword     string       `json:"confirmPassword,omitempty"`
}

type UserPrincipal struct {
	Teams       []Team       `json:"teams"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Id          int64        `json:"id,omitempty"`
	Permissions []Permission `json:"permissions"`
	Name        string       `json:"name"`
}

type IdentifiableObject struct {
	UUID uuid.UUID `json:"uuid"`
}

func (us UserService) Login(ctx context.Context, username, password string) (token string, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	body := url.Values{}
	body.Set("username", username)
	body.Set("password", password)

	req, err := us.client.newRequest(ctx, http.MethodPost, "api/v1/user/login", withBody(body))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")

	_, err = us.client.doRequest(req, &token)
	return
}

func (us UserService) ForceChangePassword(ctx context.Context, username, password, newPassword string) (err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	body := url.Values{}
	body.Set("username", username)
	body.Set("password", password)
	body.Set("newPassword", newPassword)
	body.Set("confirmPassword", newPassword)

	req, err := us.client.newRequest(ctx, http.MethodPost, "api/v1/user/forceChangePassword", withBody(body))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")

	_, err = us.client.doRequest(req, nil)
	return
}

func (us UserService) GetAllManaged(ctx context.Context, po PageOptions) (p Page[ManagedUser], err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodGet, "api/v1/user/managed", withPageOptions(po))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &p.Items)
	return
}

func (us UserService) CreateManaged(ctx context.Context, usr ManagedUser) (user ManagedUser, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodPut, "api/v1/user/managed", withBody(usr))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) UpdateManaged(ctx context.Context, usr ManagedUser) (user ManagedUser, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodPost, "api/v1/user/managed", withBody(usr))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) DeleteManaged(ctx context.Context, user ManagedUser) (err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodDelete, "api/v1/user/managed", withBody(user))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, nil)
	return
}

func (us UserService) AddTeamToUser(ctx context.Context, username string, team uuid.UUID) (user UserPrincipal, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v1/user/%s/membership", username), withBody(IdentifiableObject{
		UUID: team,
	}))
	if err != nil {
		return
	}

	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) RemoveTeamFromUser(ctx context.Context, username string, team uuid.UUID) (user UserPrincipal, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/user/%s/membership", username), withBody(IdentifiableObject{
		UUID: team,
	}))
	if err != nil {
		return
	}

	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) GetSelf(ctx context.Context) (user UserPrincipal, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodGet, "api/v1/user/self")
	if err != nil {
		return
	}

	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) UpdateSelf(ctx context.Context, userReq ManagedUser) (userRes ManagedUser, err error) {
	err = us.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := us.client.newRequest(ctx, http.MethodPost, "api/v1/user/self", withBody(userReq))
	if err != nil {
		return
	}

	_, err = us.client.doRequest(req, &userRes)
	return
}
