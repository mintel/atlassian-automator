package common

import (
	"log"
	"net/url"

	"github.com/andygrunwald/go-jira"
	"github.com/mintel/atlassian-automator/pkg/confluence"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CollectedData struct {
	Summary     string
	Description string
}

var (
	// ConfluenceAPI     *goconfluence.API
	// ConfluenceBaseURL *url.URL
	ConfluenceClient *confluence.Client
	JiraClient       *jira.Client
	PromErrors       = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "atlassian_automator_errors_total",
			Help: "The number of errors encountered by the main package",
		},
		[]string{
			"package",
		})
)

func AtlassianSetup(baseURL *url.URL, username string, password string) {
	// confluenceBaseRef, err := url.Parse("/wiki")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// confluenceAPIRef, err := url.Parse("/wiki/rest/api")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ConfluenceBaseURL = baseURL.ResolveReference(confluenceBaseRef)
	// confluenceAPIURL := *baseURL.ResolveReference(confluenceAPIRef)

	var err error

	// Set up Jira client
	ctp := confluence.BasicAuthTransport{
		Username: username,
		APIToken: password,
	}
	ConfluenceClient, err = confluence.NewClient(baseURL.String(), ctp.Client())
	if err != nil {
		log.Fatal(err)
	}
	jtp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	JiraClient, err = jira.NewClient(jtp.Client(), baseURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// Set up Confluence client
	// ConfluenceAPI, err = goconfluence.NewAPI(confluenceAPIURL.String(), username, password)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
