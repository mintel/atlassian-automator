package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/andygrunwald/go-jira"
	"github.com/mintel/atlassian-automator/pkg/common"
	"github.com/mintel/atlassian-automator/pkg/lastupdate"
	goconfluence "github.com/virtomize/confluence-go-api"
	"gopkg.in/yaml.v3"
)

var (
	baseURL           *url.URL
	confluenceAPI     *goconfluence.API
	confluenceAPIURL  url.URL
	confluenceBaseURL url.URL
	debugMode         bool
	jiraClient        *jira.Client
	wg                sync.WaitGroup
)

var args struct {
	AtlassianToken       string `env:"ATLASSIAN_TOKEN" required:"" help:"Your Atlassian API token. Either the environment variable or the flag MUST be set."`
	AtlassianUsername    string `env:"ATLASSIAN_USERNAME" required:"" help:"Your Atlassian user name. Either the environment variable or the flag MUST be set."`
	Config               string `env:"CONFIG_FILE" default:"config.yaml" type:"path" help:"Path to atlasstian-automator config file."`
	ListenAddress        string `default:":8080"`
	RedisAuthToken       string `env:"REDIS_AUTH_TOKEN"`
	RedisPort            string `env:"REDIS_PORT"`
	RedisPrimaryEndpoint string `env:"REDIS_PRIMARY_ENDPOINT"`
	RedisSSL             bool   `env:"REDIS_SSL" default:"true"`
}

type Config struct {
	Atlassian struct {
		BaseURL string `yaml:"baseURL"`
	} `yaml:"atlassian"`
	Debug        bool          `yaml:"debug"`
	IssueConfigs []IssueConfig `yaml:"issues"`
}

type IssueConfig struct {
	Interval       string            `yaml:"interval"`
	JiraLabels     []string          `yaml:"jiraLabels"`
	JiraProjectKey string            `yaml:"jiraProjectKey"`
	LastUpdate     lastupdate.Config `yaml:"lastUpdate"`
	Name           string            `yaml:"name"`
	RetryInterval  string            `yaml:"retryInterval"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func hasExistingJiraIssue(itemTitle string, projectKey string, jiraClient *jira.Client) (bool, error) {
	// Escape quotes in the title so its parsed correctly by Jira's JQL parser
	itemTitle = strings.ReplaceAll(itemTitle, `"`, `\"`)
	// Wrap the itemTitle in "\ \" so Jira does a direct match.
	//https://confluence.atlassian.com/jirasoftwareserver/search-syntax-for-text-fields-939938747.html
	jql := fmt.Sprintf("project = \"%s\" AND summary ~ \"\\\"%s\\\"\"", projectKey, itemTitle)
	log.Printf("searching for existing issue \"%s\" in project %s\n", itemTitle, projectKey)
	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return false, err
	}

	if len(issues) == 0 {
		return false, nil
	} else if len(issues) > 1 {
		log.Printf("found multiple issues that match \"%s\":", itemTitle)
		for _, x := range issues {
			log.Printf("%s ", x.Key)
		}
	}
	return true, nil
}

func raiseIssue(page *common.CollectedData, jiraProjectKey string, jiraLabels []string) (*jira.Issue, *jira.Response, error) {
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Type:        jira.IssueType{Name: "Task"},
			Project:     jira.Project{Key: jiraProjectKey},
			Description: page.Description,
			Summary:     page.Summary,
			Labels:      jiraLabels,
		},
	}
	jiraIssue, jiraResponse, err := jiraClient.Issue.Create(&issue)
	if err != nil {
		return nil, nil, err
	}
	return jiraIssue, jiraResponse, nil
}

func issueRaiser(ctx context.Context, cfg *IssueConfig) {
	defer wg.Done()
	intervalDuration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		log.Fatal(err)
	}
	retryIntervalDuration, err := time.ParseDuration(cfg.RetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			log.Printf("%s: running job", cfg.Name)
			if cfg.LastUpdate != (lastupdate.Config{}) {
				allPages, err := lastupdate.Run(*confluenceAPI, cfg.LastUpdate, &confluenceBaseURL)
				if err != nil {
					log.Print(err)
					log.Printf("Retrying in %s", cfg.RetryInterval)
					timer = time.NewTimer(retryIntervalDuration)
					break
				}
				for _, page := range allPages {
					if debugMode {
						fmt.Print(page.Summary + "\n\n")
						fmt.Print(page.Description)
					} else {
						exists, err := hasExistingJiraIssue(page.Summary, cfg.JiraProjectKey, jiraClient)
						if err != nil {
							log.Print(err)
							break
						}
						if !exists {
							log.Printf("%s: creating issue for %s", cfg.Name, page.Summary)
							jiraIssue, _, err := raiseIssue(&page, cfg.JiraProjectKey, cfg.JiraLabels)
							if err != nil {
								log.Print(err)
								break
							} else {
								log.Printf("%s: issue created for %s: %s", cfg.Name, page.Summary, jiraIssue.Key)
							}
						} else {
							log.Printf("%s: issue already exists for %s", cfg.Name, page.Summary)
						}
					}
				}
			}
			log.Printf("%s: job complete.", cfg.Name)
			log.Printf("%s: waiting for %s", cfg.Name, cfg.Interval)
			timer = time.NewTimer(intervalDuration)
		}
	}
}

func main() {

	// Parse command line arguments
	_ = kong.Parse(&args)

	// Parse config file
	if err := ValidateConfigPath(args.Config); err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(args.Config)
	if err != nil {
		log.Fatal(err)
	}

	// Set global debugMode variable
	debugMode = cfg.Debug

	// Create URLs for API libraries
	baseURL, err = url.Parse(cfg.Atlassian.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	confluenceBaseRef, err := url.Parse("/wiki")
	if err != nil {
		log.Fatal(err)
	}
	confluenceAPIRef, err := url.Parse("/wiki/rest/api")
	if err != nil {
		log.Fatal(err)
	}
	confluenceBaseURL = *baseURL.ResolveReference(confluenceBaseRef)
	confluenceAPIURL = *baseURL.ResolveReference(confluenceAPIRef)

	// Set up Jira client
	tp := jira.BasicAuthTransport{
		Username: args.AtlassianUsername,
		Password: args.AtlassianToken,
	}
	jiraClient, err = jira.NewClient(tp.Client(), baseURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// Set up Confluence client
	confluenceAPI, err = goconfluence.NewAPI(confluenceAPIURL.String(), args.AtlassianUsername, args.AtlassianToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set up OS signal notifications
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start job goroutines
	if cfg.IssueConfigs != nil {
		wg.Add(len(cfg.IssueConfigs))
		for _, ic := range cfg.IssueConfigs {
			go issueRaiser(ctx, &ic)
		}
	}

	// Start HTTP server goroutine for healthchecks
	httpServer := http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	go httpServer.ListenAndServe()

	// Sit and wait for an OS SIGTERM / SIGINT then shut everything down when received
	<-ctx.Done()
	log.Print("shutting down gracefully")
	stop()
	wg.Wait()

	// Give 5s more to process existing requests
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(timeoutCtx); err != nil {
		log.Fatal(err)
	}

}
