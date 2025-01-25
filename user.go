package dtrack

import (
	"context"
	"net/http"
	"net/url"
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

func (us UserService) Login(ctx context.Context, username, password string) (token string, err error) {
	body := url.Values{}
	body.Set("username", username)
	body.Set("password", password)

	req, err := us.client.newRequest(ctx, http.MethodPost, "/api/v1/user/login", withBody(body))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")

	_, err = us.client.doRequest(req, &token)
	return
}

func (us UserService) ForceChangePassword(ctx context.Context, username, password, newPassword string) (err error) {
	body := url.Values{}
	body.Set("username", username)
	body.Set("password", password)
	body.Set("newPassword", newPassword)
	body.Set("confirmPassword", newPassword)

	req, err := us.client.newRequest(ctx, http.MethodPost, "/api/v1/user/forceChangePassword", withBody(body))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")

	_, err = us.client.doRequest(req, nil)
	return
}

func (us UserService) GetAllManaged(ctx context.Context, po PageOptions) (p Page[ManagedUser], err error) {
	req, err := us.client.newRequest(ctx, http.MethodGet, "/api/v1/user/managed", withPageOptions(po))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &p.Items)
	return
}

func (us UserService) CreateManaged(ctx context.Context, usr ManagedUser) (user ManagedUser, err error) {
	req, err := us.client.newRequest(ctx, http.MethodPut, "/api/v1/user/managed", withBody(usr))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) UpdateManaged(ctx context.Context, usr ManagedUser) (user ManagedUser, err error) {
	req, err := us.client.newRequest(ctx, http.MethodPost, "/api/v1/user/managed", withBody(usr))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &user)
	return
}

func (us UserService) DeleteManaged(ctx context.Context, user ManagedUser) (err error) {
	req, err := us.client.newRequest(ctx, http.MethodDelete, "/api/v1/user/managed", withBody(user))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, nil)
	return
}
