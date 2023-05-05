package confluence

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// PageService handles Pages for the Confluence instance / API.
//
// Confluence API docs: https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-page/#api-group-page
type PageService service

// Page represents a Confluence page.
type Page struct {
	ID        string   `json:"id,omitempty" structs:"id,omitempty"`
	Status    string   `json:"status,omitempty" structs:"status,omitempty"`
	Title     string   `json:"title,omitempty" structs:"title,omitempty"`
	SpaceID   string   `json:"spaceId,omitempty" structs:"spaceId,omitempty"`
	ParentID  string   `json:"parentId,omitempty" structs:"parentId,omitempty"`
	AuthorID  string   `json:"authorId,omitempty" structs:"authorId,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	Version   *Version `json:"version,omitempty" structs:"version,omitempty"`
}

type Pages struct {
	Results []*Page `json:"results,omitempty" structs:"results,omitempty"`
	Links   *Links  `json:"_links,omitempty" structs:"_links,omitempty"`
}

type Version struct {
	CreatedAt string `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	Message   string `json:"message,omitempty" structs:"message,omitempty"`
	Number    int    `json:"number,omitempty" structs:"number,omitempty"`
	MinorEdit bool   `json:"minorEdit,omitempty" structs:"minorEdit,omitempty"`
	AuthorId  string `json:"authorId,omitempty" structs:"authorId,omitempty"`
}

type Body struct {
	Storage        *BodyType `json:"storage,omitempty" structs:"storage,omitempty"`
	AtlasDocFormat *BodyType `json:"atlas_doc_format,omitempty" structs:"atlas_doc_format,omitempty"`
}

type GetPageByIdOptions struct {
	BodyFormat            string `url:"body-format,omitempty"`
	GetDraft              bool   `url:"get-draft,omitempty"`
	Version               int    `url:"version,omitempty"`
	SerializeIDsAsStrings bool   `url:"serialize-ids-as-strings,omitempty"`
}

type GetPagesInSpaceOptions struct {
	Status                string `url:"status,omitempty"`
	BodyFormat            string `url:"body-format,omitempty"`
	Cursor                string `url:"cursor,omitempty"`
	Limit                 int    `url:"limit,omitempty"`
	SerializeIDsAsStrings bool   `url:"serialize-ids-as-strings,omitempty"`
}

func (s *PageService) GetPageById(ctx context.Context, pageID string, options *GetPageByIdOptions) (*Page, *Response, error) {
	apiEndpoint := fmt.Sprintf("wiki/api/v2/pages/%s", pageID)
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

	page := new(Page)
	resp, err := s.client.Do(req, page)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	return page, resp, nil
}

func (s *PageService) GetPagesInSpace(ctx context.Context, spaceID string, options *GetPagesInSpaceOptions) (*Pages, *Response, error) {
	apiEndpoint := fmt.Sprintf("wiki/api/v2/spaces/%s/pages", spaceID)
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

	pages := new(Pages)
	resp, err := s.client.Do(req, pages)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if pages.Links.Next != "" {
		extraPages, _, err := s.GetNextPages(ctx, pages.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		pages.Results = append(pages.Results, extraPages.Results...)
	}
	return pages, resp, nil
}

func (s *PageService) GetNextPages(ctx context.Context, url string) (*Pages, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	pages := new(Pages)
	resp, err := s.client.Do(req, pages)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if pages.Links.Next != "" {
		extraPages, _, err := s.GetNextPages(ctx, pages.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		pages.Results = append(pages.Results, extraPages.Results...)
	}
	return pages, resp, nil
}
