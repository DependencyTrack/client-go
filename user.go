package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

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

type UserDefinitions struct {
	Username            string `json:"username"`
	NewPassword         string `json:"newPassword"`
	ConfirmPassword     string `json:"confirmPassword"`
	FullName            string `json:"fullname"`
	Email               string `json:"email"`
	Suspended           bool   `json:"suspended,omitempty"`
	ForcePasswordChange bool   `json:"forcePasswordChange,omitempty"`
	NonExpiryPassword   bool   `json:"nonExpiryPassword,omitempty"`

	//TODO: LastPasswordChange, Teams, Permissions
}

func (us UserService) CreateUser(ctx context.Context, user *UserDefinitions, token string) (err error) {

	req, err := us.client.newRequest(ctx, http.MethodPut, "/api/v1/user/managed", withBody(user))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	_, err = us.client.doRequest(req, nil)
	return
}

func (us UserService) DeleteUser(ctx context.Context, user *UserDefinitions, token string) (err error) {

	req, err := us.client.newRequest(ctx, http.MethodDelete, "/api/v1/user/managed", withBody(user))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	_, err = us.client.doRequest(req, nil)
	return
}
