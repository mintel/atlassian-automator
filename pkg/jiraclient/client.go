//go:generate interfacer -for github.com/andygrunwald/go-jira.Client -as jiraclient.Client -o client_interface.go
//go:generate mockgen -package mock_jiraclient -destination mock_jiraclient/client.go . Client
package jiraclient

import (
	"context"
	"net/url"

	"github.com/andygrunwald/go-jira"
)

type contextKey string

var (
	clientContextKey contextKey = "jiraclient.client"
)

func ClientFor(ctx context.Context, baseURL *url.URL, username string, password string) (Client, error) {

	client, ok := ctx.Value(clientContextKey).(Client)
	if ok && client != nil {
		return client, nil
	}

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(tp.Client(), baseURL.String())
	if err != nil {
		return nil, err
	}

	return client, nil
}
