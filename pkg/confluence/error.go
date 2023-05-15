package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Error struct {
	HTTPError error
	Errors    []*ConfluenceError `json:"errors"`
}

type ConfluenceError struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func NewConfluenceError(resp *Response, httpError error) error {
	if resp == nil {
		return fmt.Errorf("no response returned: %w", httpError)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", httpError.Error(), err)
	}
	cerr := Error{HTTPError: httpError}
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		err = json.Unmarshal(body, &cerr)
		if err != nil {
			return fmt.Errorf("%s: could not parse JSON: %w", httpError.Error(), err)
		}
	} else {
		if httpError == nil {
			return fmt.Errorf("got response status %s:%s", resp.Status, string(body))
		}
		return fmt.Errorf("%s: %s: %w", resp.Status, string(body), httpError)
	}

	return &cerr
}

// Error is a short string representing the error
func (e *Error) Error() string {
	var output string
	if len(e.Errors) > 0 {
		for _, err := range e.Errors {
			output += fmt.Sprintf("%v - %s - %s - %s ", err.Status, err.Code, err.Title, err.Detail)
		}
	}
	return output
}
