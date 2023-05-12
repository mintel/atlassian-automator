package confluence

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// ChildrenService handles Child Pages for the Confluence instance / API.
//
// Confluence API docs: https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-children/#api-group-children
type ChildrenService service

// ChildPage represents a Confluence child page.
type ChildPage struct {
	ID            string `json:"id,omitempty" structs:"id,omitempty"`
	Status        string `json:"status,omitempty" structs:"status,omitempty"`
	Title         string `json:"title,omitempty" structs:"title,omitempty"`
	SpaceID       string `json:"spaceId,omitempty" structs:"spaceId,omitempty"`
	ChildPosition int    `json:"childPosition,omitempty" structs:"childPosition,omitempty"`
}

// ChildPages contains a list of child pages plus a URL for the next list of pages if the limit for this requesthas been
// reached
type ChildPages struct {
	Results []*ChildPage `json:"results,omitempty" structs:"results,omitempty"`
	Links   *Links       `json:"_links,omitempty" structs:"_links,omitempty"`
}

// ChildCustomContent represents a Confluence ChildCustomContent response
type ChildCustomContent struct {
	ID      string `json:"id,omitempty" structs:"id,omitempty"`
	Status  string `json:"status,omitempty" structs:"status,omitempty"`
	Title   string `json:"title,omitempty" structs:"title,omitempty"`
	Type    string `json:"type,omitempty" structs:"type,omitempty"`
	SpaceID string `json:"spaceId,omitempty" structs:"spaceId,omitempty"`
}

// ChildCustomContents contains a list of ChildCustomConten objects plus a URL for the next list of pages if the limit
// for this requesthas been reached
type ChildCustomContents struct {
	Results []*ChildCustomContent `json:"results,omitempty" structs:"results,omitempty"`
	Links   *Links                `json:"_links,omitempty" structs:"_links,omitempty"`
}

// GetChildPagesOptions is the query parameters that can be passed to the GetChildPages API call
type GetChildPagesOptions struct {
	Cursor                string `url:"cursor,omitempty"`
	Limit                 int    `url:"limit,omitempty"`
	Sort                  string `url:"sort,omitempty"`
	SerializeIDsAsStrings bool   `url:"serialize-ids-as-strings,omitempty"`
}

// GetChildCustomContentOptions is the query parameters that can be passed to the GetChildCustomContentOptions API call
type GetChildCustomContentOptions struct {
	Cursor                string `url:"cursor,omitempty"`
	Limit                 int    `url:"limit,omitempty"`
	Sort                  string `url:"sort,omitempty"`
	SerializeIDsAsStrings bool   `url:"serialize-ids-as-strings,omitempty"`
}

// GetChildPages returns all child pages for given page id. The number of results is limited by the limit parameter and
// additional results (if available) will be retrieved using the GetNextChildPages function if necessary.
func (s *ChildrenService) GetChildPages(ctx context.Context, pageID string, options *GetChildPagesOptions) (*ChildPages, *Response, error) {
	apiEndpoint := fmt.Sprintf("wiki/api/v2/pages/%s/children", pageID)
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

	childPages := new(ChildPages)
	resp, err := s.client.Do(req, childPages)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if childPages.Links.Next != "" {
		extraChildPages, _, err := s.GetNextChildPages(ctx, childPages.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		childPages.Results = append(childPages.Results, extraChildPages.Results...)
	}
	return childPages, resp, nil
}

// GetNextChildPages takes the URL from the response of a GetChildPages request and continues to call itself until there
// are no pages left (i.e. a links.next value is not provided in the final API response). Returns all resulting
// ChildPage objects.
func (s *ChildrenService) GetNextChildPages(ctx context.Context, url string) (*ChildPages, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	childPages := new(ChildPages)
	resp, err := s.client.Do(req, childPages)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if childPages.Links.Next != "" {
		extraChildPages, _, err := s.GetNextChildPages(ctx, childPages.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		childPages.Results = append(childPages.Results, extraChildPages.Results...)
	}
	return childPages, resp, nil
}

// GetChildCustomContent returns all child custom content for given custom content id. The number of results is limited
// by the limit parameter and additional results (if available) will be retrieved using the GetNextChildCustomContents
// function if necessary.
func (s *ChildrenService) GetChildCustomContent(ctx context.Context, pageID string, options *GetChildCustomContentOptions) (*ChildCustomContents, *Response, error) {
	apiEndpoint := fmt.Sprintf("wiki/api/v2/custom-content/%s/children", pageID)
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

	childCustomContents := new(ChildCustomContents)
	resp, err := s.client.Do(req, childCustomContents)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if childCustomContents.Links.Next != "" {
		extraChildCustomContents, _, err := s.GetNextChildCustomContents(ctx, childCustomContents.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		childCustomContents.Results = append(childCustomContents.Results, extraChildCustomContents.Results...)
	}
	return childCustomContents, resp, nil
}

// GetNextChildCustomContents takes the URL from the response of a GetChildCustomContents request and continues to call
// itself until there are no pages left (i.e. a links.next value is not provided in the final API response). Returns all
// resulting GetNextChildCustomContent objects.
func (s *ChildrenService) GetNextChildCustomContents(ctx context.Context, url string) (*ChildCustomContents, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	childCustomContents := new(ChildCustomContents)
	resp, err := s.client.Do(req, childCustomContents)
	if err != nil {
		cerr := NewConfluenceError(resp, err)
		return nil, resp, cerr
	}
	if childCustomContents.Links.Next != "" {
		extraChildCustomContents, _, err := s.GetNextChildCustomContents(ctx, childCustomContents.Links.Next)
		if err != nil {
			return nil, nil, err
		}
		childCustomContents.Results = append(childCustomContents.Results, extraChildCustomContents.Results...)
	}
	return childCustomContents, resp, nil
}
