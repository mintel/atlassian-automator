package confluence

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

// SpaceService handles Spaces for the Confluence instance / API.
//
// Confluence API docs: https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space/#api-group-space
type SpaceService service

type Space struct {
	ID          string       `json:"id,omitempty" structs:"id,omitempty"`
	Key         string       `json:"key,omitempty" structs:"key,omitempty"`
	Name        string       `json:"name,omitempty" structs:"name,omitempty"`
	Type        string       `json:"type,omitempty" structs:"type,omitempty"`
	Status      string       `json:"status,omitempty" structs:"status,omitempty"`
	HomepageID  string       `json:"homepageId,omitempty" structs:"homepageId,omitempty"`
	Description *Description `json:"description,omitempty" structs:"description,omitempty"`
}

type Description struct {
	Plain *BodyType `json:"plain,omitempty" structs:"plain,omitempty"`
	View  *BodyType `json:"view,omitempty" structs:"view,omitempty"`
}

type Spaces struct {
	Results []*Space `json:"results,omitempty" structs:"results,omitempty"`
	Links   *Links   `json:"_links,omitempty" structs:"_links,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty" structs:"next,omitempty"`
}

type GetSpacesOptions struct {
	IDs                   []string `url:"ids,omitempty"`
	Keys                  []string `url:"keys,omitempty"`
	Type                  string   `url:"type,omitempty"`
	Status                string   `url:"status,omitempty"`
	Labels                []string `url:"labels,omitempty"`
	Sort                  string   `url:"sort,omitempty"`
	DescriptionFormat     string   `url:"description-format,omitempty"`
	Cursor                string   `url:"cursor,omitempty"`
	Limit                 int      `url:"limit,omitempty"`
	SerializeIDsAsStrings bool     `url:"serialize-ids-as-strings,omitempty"`
}

func (s *SpaceService) GetSpaces(ctx context.Context, options *GetSpacesOptions) (*Spaces, *Response, error) {
	apiEndpoint := "wiki/api/v2/spaces"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	spaces := new(Spaces)
	resp, err := s.client.Do(req, spaces)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	return spaces, resp, nil
}
