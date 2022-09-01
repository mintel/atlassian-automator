package common

import (
	"log"
	"net/url"

	"github.com/andygrunwald/go-jira"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	goconfluence "github.com/virtomize/confluence-go-api"
)

type CollectedData struct {
	Summary     string
	Description string
}

var (
	ConfluenceAPI     *goconfluence.API
	ConfluenceBaseURL *url.URL
	JiraClient        *jira.Client
	PromErrors        = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "atlassian_automator_errors_total",
			Help: "The number of errors encountered by the main package",
		},
		[]string{
			"package",
		})
)

func AtlassianSetup(baseURL *url.URL, username string, password string) {
	confluenceBaseRef, err := url.Parse("/wiki")
	if err != nil {
		log.Fatal(err)
	}
	confluenceAPIRef, err := url.Parse("/wiki/rest/api")
	if err != nil {
		log.Fatal(err)
	}
	ConfluenceBaseURL = baseURL.ResolveReference(confluenceBaseRef)
	confluenceAPIURL := *baseURL.ResolveReference(confluenceAPIRef)

	// Set up Jira client
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	JiraClient, err = jira.NewClient(tp.Client(), baseURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// Set up Confluence client
	ConfluenceAPI, err = goconfluence.NewAPI(confluenceAPIURL.String(), username, password)
	if err != nil {
		log.Fatal(err)
	}
}
