package common

import (
	"log"
	"net/url"

	"github.com/andygrunwald/go-jira"
	"github.com/mintel/atlassian-automator/pkg/confluenceclient"
	"github.com/mintel/atlassian-automator/pkg/jiraclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	goconfluence "github.com/virtomize/confluence-go-api"
)

type CollectedData struct {
	Summary     string
	Description string
}

var (
	ConfluenceClient  *goconfluence.API
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
	ConfluenceBaseURL = baseURL.ResolveReference(confluenceBaseRef)
	// Set up Jira Client
	JiraClient = jiraclient.ClientFor(baseURL, username, password)
	// Set up Confluence client
	ConfluenceClient = confluenceclient.ClientFor(baseURL, username, password)
}
