//go:generate interfacer -for github.com/virtomize/confluence-go-api.API -as confluenceclient.Client -o client_interface.go
//go:generate mockgen -package mock_confluenceclient -destination mock_confluenceclient/client.go . Client
package confluenceclient

import (
	"log"
	"net/url"

	goconfluence "github.com/virtomize/confluence-go-api"
)

func ClientFor(baseURL *url.URL, username string, password string) *goconfluence.API {
	confluenceAPIRef, err := url.Parse("/wiki/rest/api")
	if err != nil {
		log.Fatal(err)
	}
	confluenceAPIURL := *baseURL.ResolveReference(confluenceAPIRef)
	ConfluenceClient, err := goconfluence.NewAPI(confluenceAPIURL.String(), username, password)
	if err != nil {
		log.Fatal(err)
	}
	return ConfluenceClient
}
