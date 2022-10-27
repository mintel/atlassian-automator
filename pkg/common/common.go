package common

import (
	"context"
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

func AtlassianSetup(ctx context.Context, baseURL *url.URL, username string, password string) {

	var ok bool

	// Set up Jira Client
	jiraClient, err := jiraclient.ClientFor(ctx, baseURL, username, password)
	if err != nil {
		log.Fatal(err)
	}
	JiraClient, ok = jiraClient.(*jira.Client)
	if !ok {
		log.Fatal(ok)
	}

	// Set up Confluence client
	confluenceBaseRef, err := url.Parse("/wiki")
	if err != nil {
		log.Fatal(err)
	}
	ConfluenceBaseURL = baseURL.ResolveReference(confluenceBaseRef)
	confluenceClient, err := confluenceclient.ClientFor(ctx, baseURL, username, password)
	if err != nil {
		log.Fatal(err)
	}
	ConfluenceClient, ok = confluenceClient.(*goconfluence.API)
	if !ok {
		log.Fatal(ok)
	}
}
