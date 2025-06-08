package dtrack

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Component struct {
	UUID               uuid.UUID                `json:"uuid,omitempty"`
	Author             string                   `json:"author,omitempty"`
	Publisher          string                   `json:"publisher,omitempty"`
	Group              string                   `json:"group,omitempty"`
	Name               string                   `json:"name"`
	Version            string                   `json:"version"`
	Classifier         string                   `json:"classifier,omitempty"`
	FileName           string                   `json:"filename,omitempty"`
	Extension          string                   `json:"extension,omitempty"`
	MD5                string                   `json:"md5,omitempty"`
	SHA1               string                   `json:"sha1,omitempty"`
	SHA256             string                   `json:"sha256,omitempty"`
	SHA384             string                   `json:"sha384,omitempty"`
	SHA512             string                   `json:"sha512,omitempty"`
	SHA3_256           string                   `json:"sha3_256,omitempty"`
	SHA3_384           string                   `json:"sha3_384,omitempty"`
	SHA3_512           string                   `json:"sha3_512,omitempty"`
	BLAKE2b_256        string                   `json:"blake2b_256,omitempty"`
	BLAKE2b_384        string                   `json:"blake2b_384,omitempty"`
	BLAKE2b_512        string                   `json:"blake2b_512,omitempty"`
	BLAKE3             string                   `json:"blake3,omitempty"`
	CPE                string                   `json:"cpe,omitempty"`
	PURL               string                   `json:"purl,omitempty"`
	SWIDTagID          string                   `json:"swidTagId,omitempty"`
	Internal           bool                     `json:"isInternal,omitempty"`
	Description        string                   `json:"description,omitempty"`
	Copyright          string                   `json:"copyright,omitempty"`
	License            string                   `json:"license,omitempty"`
	ResolvedLicense    *License                 `json:"resolvedLicense,omitempty"`
	DirectDependencies string                   `json:"directDependencies,omitempty"`
	Notes              string                   `json:"notes,omitempty"`
	ExternalReferences []ExternalReference      `json:"externalReferences,omitempty"`
	Project            *Project                 `json:"project,omitempty"`
	RepositoryMeta     *RepositoryMetaComponent `json:"repositoryMeta,omitempty"`
}

type ExternalReference struct {
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type ComponentService struct {
	client *Client
}

type ComponentFilterOptions struct {
	OnlyOutdated bool
	OnlyDirect   bool
}

type ComponentProperty struct {
	Group       string    `json:"groupName,omitempty"`
	Name        string    `json:"propertyName,omitempty"`
	Value       string    `json:"propertyValue,omitempty"`
	Type        string    `json:"propertyType"`
	Description string    `json:"description,omitempty"`
	UUID        uuid.UUID `json:"uuid"`
}

type ComponentIdentityQueryOptions struct {
	Group     string
	Name      string
	Version   string
	PURL      string
	CPE       string
	SWIDTagID string
	Project   uuid.UUID
}

func (cs ComponentService) Get(ctx context.Context, componentUUID uuid.UUID) (c Component, err error) {
	err = cs.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/%s", componentUUID))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, &c)
	return
}

func (cs ComponentService) GetAll(ctx context.Context, projectUUID uuid.UUID, po PageOptions, filterOptions ComponentFilterOptions) (p Page[Component], err error) {
	err = cs.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/project/%s", projectUUID), withPageOptions(po), withComponentFilterOptions(filterOptions))
	if err != nil {
		return
	}

	res, err := cs.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func withComponentFilterOptions(filterOptions ComponentFilterOptions) requestOption {
	return func(req *http.Request) error {
		query := req.URL.Query()
		if filterOptions.OnlyDirect {
			query.Set("onlyDirect", "true")
		}
		if filterOptions.OnlyOutdated {
			query.Set("onlyOutdated", "true")
		}
		req.URL.RawQuery = query.Encode()
		return nil
	}
}

func (cs ComponentService) Create(ctx context.Context, projectUUID uuid.UUID, component Component) (c Component, err error) {
	err = cs.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/component/project/%s", projectUUID), withBody(component))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, &c)
	return
}

func (cs ComponentService) Update(ctx context.Context, component Component) (c Component, err error) {
	err = cs.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodPost, "/api/v1/component", withBody(component))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, &c)
	return
}

func (cs ComponentService) Delete(ctx context.Context, componentUUID uuid.UUID) (err error) {
	err = cs.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/component/%s", componentUUID))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, nil)
	return
}

func (cs ComponentService) GetProperties(ctx context.Context, componentUUID uuid.UUID) (ps []ComponentProperty, err error) {
	err = cs.client.assertServerVersionAtLeast("4.11.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/%s/property", componentUUID))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, &ps)
	return
}

func (cs ComponentService) CreateProperty(ctx context.Context, componentUUID uuid.UUID, property ComponentProperty) (p ComponentProperty, err error) {
	err = cs.client.assertServerVersionAtLeast("4.11.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/component/%s/property", componentUUID), withBody(property))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, &p)
	return
}

func (cs ComponentService) DeleteProperty(ctx context.Context, componentUUID, propertyUUID uuid.UUID) (err error) {
	err = cs.client.assertServerVersionAtLeast("4.11.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/component/%s/property/%s", componentUUID, propertyUUID))
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, nil)
	return
}

func (cs ComponentService) GetByHash(ctx context.Context, hash string, po PageOptions, so SortOptions) (p Page[Component], err error) {
	err = cs.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/hash/%s", hash), withPageOptions(po), withSortOptions(so))
	if err != nil {
		return
	}

	res, err := cs.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	if cs.client.isServerVersionAtLeast("4.0.0") {
		p.TotalCount = res.TotalCount
	} else {
		p.TotalCount = len(p.Items)
	}
	return
}

func (cs ComponentService) GetByIdentity(ctx context.Context, po PageOptions, so SortOptions, io ComponentIdentityQueryOptions) (p Page[Component], err error) {
	err = cs.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, "/api/v1/component/identity", withPageOptions(po), withSortOptions(so), withComponentIdentityOptions(io))
	if err != nil {
		return
	}

	res, err := cs.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func withComponentIdentityOptions(identityOptions ComponentIdentityQueryOptions) requestOption {
	return func(req *http.Request) error {
		query := req.URL.Query()
		if len(identityOptions.Group) > 0 {
			query.Set("group", identityOptions.Group)
		}
		if len(identityOptions.Name) > 0 {
			query.Set("name", identityOptions.Name)
		}
		if len(identityOptions.Version) > 0 {
			query.Set("version", identityOptions.Version)
		}
		if len(identityOptions.PURL) > 0 {
			query.Set("purl", identityOptions.PURL)
		}
		if len(identityOptions.CPE) > 0 {
			query.Set("cpe", identityOptions.CPE)
		}
		if len(identityOptions.SWIDTagID) > 0 {
			query.Set("swidTagId", identityOptions.SWIDTagID)
		}
		if identityOptions.Project != uuid.Nil {
			query.Set("project", identityOptions.Project.String())
		}
		req.URL.RawQuery = query.Encode()
		return nil
	}
}

func (cs ComponentService) IdentifyInternal(ctx context.Context) (err error) {
	err = cs.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := cs.client.newRequest(ctx, http.MethodGet, "/api/v1/component/internal/identify")
	if err != nil {
		return
	}

	_, err = cs.client.doRequest(req, nil)
	return
}
