//go:generate interfacer -for github.com/virtomize/confluence-go-api.API -as confluenceclient.Client -o client_interface.go
//go:generate mockgen -package mock_confluenceclient -destination mock_confluenceclient/client.go . Client
package confluenceclient

import (
	"context"
	"net/url"

	goconfluence "github.com/virtomize/confluence-go-api"
)

type contextKey string

var (
	clientContextKey contextKey = "confluenceclient.client"
)

func ClientFor(ctx context.Context, baseURL *url.URL, username string, password string) (Client, error) {

	client, ok := ctx.Value(clientContextKey).(Client)
	if ok && client != nil {
		return client, nil
	}

	confluenceAPIRef, err := url.Parse("/wiki/rest/api")
	if err != nil {
		return nil, err
	}
	confluenceAPIURL := *baseURL.ResolveReference(confluenceAPIRef)
	client, err = goconfluence.NewAPI(confluenceAPIURL.String(), username, password)
	if err != nil {
		return nil, err
	}
	return client, nil
}
