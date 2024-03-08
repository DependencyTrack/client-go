package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type User struct {
	Username          string       `json:"username,omitempty"`
	SubjectIdentifier string       `json:"subjectIdentifier,omitempty"`
	Email             string       `json:"email,omitempty"`
	Teams             []Team       `json:"teams,omitempty"`
	Permissions       []Permission `json:"permissions,omitempty"`
}

type UserService struct {
	client *Client
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

func (us UserService) GetAll(ctx context.Context) (p Page[User], err error) {
	req, err := us.client.newRequest(ctx, http.MethodGet, "/api/v1/user/oidc")
	if err != nil {
		return
	}

	res, err := us.client.doRequest(req, &p.Items)

	p.TotalCount = res.TotalCount

	return
}

func (us UserService) Membership(ctx context.Context, user *User, team *Team) (u User, err error) {
	body := make(map[string]string)
	body["uuid"] = team.UUID.String()

	req, err := us.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/user/%s/membership", user.Username), withBody(body))
	if err != nil {
		return
	}
	_, err = us.client.doRequest(req, &u)
	return
}
