package dtrack

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type License struct {
	UUID                uuid.UUID `json:"uuid"`
	Name                string    `json:"name"`
	Text                string    `json:"licenseText,omitempty"`
	Template            string    `json:"standardLicenseTemplate,omitempty"`
	Header              string    `json:"standardLicenseHeader,omitempty"`
	Comment             string    `json:"licenseComments,omitempty"`
	LicenseID           string    `json:"licenseId"`
	OSIApproved         bool      `json:"isOsiApproved"`
	FSFLibre            bool      `json:"isFsfLibre"`
	DeprecatedLicenseID bool      `json:"isDeprecatedLicenseId"`
	SeeAlso             []string  `json:"seeAlso,omitempty"`
}

type LicenseService struct {
	client *Client
}

func (l LicenseService) GetAll(ctx context.Context, po PageOptions) (p Page[License], err error) {
	req, err := l.client.newRequest(ctx, http.MethodGet, "api/v1/license", withPageOptions(po))
	if err != nil {
		return
	}

	res, err := l.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}
