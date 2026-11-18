package dtrack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ExtensionService struct {
	client *Client
}

type ExtensionPoint struct {
	Name string `json:"name"`
}

type Extension struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Configurable bool   `json:"configurable"`
	Testable     bool   `json:"testable"`
}

type ExtensionTestCheckStatus string

const (
	ExtensionTestCheckStatusPassed  ExtensionTestCheckStatus = "PASSED"
	ExtensionTestCheckStatusFailed  ExtensionTestCheckStatus = "FAILED"
	ExtensionTestCheckStatusSkipped ExtensionTestCheckStatus = "SKIPPED"
)

type ExtensionTestCheck struct {
	Name    string                   `json:"name"`
	Status  ExtensionTestCheckStatus `json:"status"`
	Message string                   `json:"message,omitempty"`
}

type ExtensionTestResult struct {
	Checks []ExtensionTestCheck `json:"checks"`
}

// Failed reports whether at least one check failed.
func (r ExtensionTestResult) Failed() bool {
	for _, check := range r.Checks {
		if check.Status == ExtensionTestCheckStatusFailed {
			return true
		}
	}
	return false
}

func (es ExtensionService) GetAllExtensionPoints(ctx context.Context) (points []ExtensionPoint, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, "api/v2/extension-points")
	if err != nil {
		return
	}

	var response struct {
		Items []ExtensionPoint `json:"items"`
	}
	_, err = es.client.doRequest(req, &response)
	return response.Items, err
}

func (es ExtensionService) GetAllExtensions(ctx context.Context, extensionPoint string) (extensions []Extension, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	path := fmt.Sprintf("api/v2/extension-points/%s/extensions", url.PathEscape(extensionPoint))
	req, err := es.client.newRequest(ctx, http.MethodGet, path)
	if err != nil {
		return
	}

	var response struct {
		Items []Extension `json:"items"`
	}
	_, err = es.client.doRequest(req, &response)
	return response.Items, err
}

func extensionPath(extensionPoint, extension, resource string) string {
	return fmt.Sprintf("api/v2/extension-points/%s/extensions/%s/%s",
		url.PathEscape(extensionPoint), url.PathEscape(extension), resource)
}

// GetConfig returns the configuration object without the API's config envelope.
func (es ExtensionService) GetConfig(ctx context.Context, extensionPoint, extension string) (config json.RawMessage, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, extensionPath(extensionPoint, extension, "config"))
	if err != nil {
		return
	}

	var response struct {
		Config json.RawMessage `json:"config"`
	}
	_, err = es.client.doRequest(req, &response)
	return response.Config, err
}

// UpdateConfig supplies a configuration object. Fields marked x-secret-ref in the
// schema must contain managed secret names, not clear text credentials.
// An unchanged configuration (HTTP 304) is treated as success.
func (es ExtensionService) UpdateConfig(ctx context.Context, extensionPoint, extension string, config json.RawMessage) (err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	body := struct {
		Config json.RawMessage `json:"config"`
	}{Config: config}
	req, err := es.client.newRequest(ctx, http.MethodPut, extensionPath(extensionPoint, extension, "config"), withBody(body))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, nil)
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotModified {
		return nil
	}
	return
}

// GetConfigSchema returns the JSON Schema, or nil if the extension has no schema.
func (es ExtensionService) GetConfigSchema(ctx context.Context, extensionPoint, extension string) (schema json.RawMessage, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, extensionPath(extensionPoint, extension, "config-schema"))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &schema)
	return
}

// Test requires a configuration for configurable extensions. For other extensions,
// config may be nil. Secret fields must contain managed secret names.
func (es ExtensionService) Test(ctx context.Context, extensionPoint, extension string, config json.RawMessage) (result ExtensionTestResult, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	body := struct {
		Config json.RawMessage `json:"config,omitempty"`
	}{Config: config}
	req, err := es.client.newRequest(ctx, http.MethodPost, extensionPath(extensionPoint, extension, "test"), withBody(body))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &result)
	return
}
