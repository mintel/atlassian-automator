//go:generate interfacer -for github.com/andygrunwald/go-jira.Client -as jiraclient.Client -o client_interface.go
//go:generate mockgen -package mock_jiraclient -destination mock_jiraclient/client.go . Client
package jiraclient

import (
	"log"
	"net/url"

	"github.com/andygrunwald/go-jira"
)

func ClientFor(baseURL *url.URL, username string, password string) *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	JiraClient, err := jira.NewClient(tp.Client(), baseURL.String())
	if err != nil {
		log.Fatal(err)
	}
	return JiraClient
}
