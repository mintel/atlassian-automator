package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/andygrunwald/go-jira"
	"github.com/mintel/atlassian-automator/pkg/lastupdate"
	goconfluence "github.com/virtomize/confluence-go-api"
	"gopkg.in/yaml.v3"
)

var (
	baseURL         *url.URL
	confluenceAPI   *goconfluence.API
	debugMode       bool
	goconfluenceURL url.URL
	gojiraURL       url.URL
	jiraClient      *jira.Client
	osSignals       chan os.Signal
	wg              sync.WaitGroup
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

func hasExistingJiraIssue(itemTitle string, projectKey string, jiraClient *jira.Client) bool {
	// Escape quotes in the title so its parsed correctly by Jira's JQL parser
	itemTitle = strings.ReplaceAll(itemTitle, `"`, `\"`)
	// Wrap the itemTitle in "\ \" so Jira does a direct match.
	//https://confluence.atlassian.com/jirasoftwareserver/search-syntax-for-text-fields-939938747.html
	jql := fmt.Sprintf("project = \"%s\" AND summary ~ \"\\\"%s\\\"\"", projectKey, itemTitle)
	log.Printf("Searching for existing issue \"%s\" in project %s\n", itemTitle, projectKey)
	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		log.Printf("Issue search failed for JQL: %s", jql)
		panic(err)
	}

	if len(issues) == 0 {
		return false
	} else if len(issues) > 1 {
		log.Printf("Found multiple issues that match \"%s\":", itemTitle)
		for _, x := range issues {
			log.Printf("%s ", x.Key)
		}
	}
	return true
}

func issueRaiser(cfg *IssueConfig) {
	defer wg.Done()
	intervalDuration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		log.Fatal(err)
	}
	timer := time.NewTimer(intervalDuration)
	log.Printf("%s: waiting for %s", cfg.Name, cfg.Interval)
	for {
		select {
		case signal := <-osSignals:
			log.Printf("%s signal received. Shutting down.", signal.String())
			return
		case <-timer.C:
			log.Printf("%s: running job", cfg.Name)
			if cfg.LastUpdate != (lastupdate.Config{}) {
				allPages, err := lastupdate.Run(*confluenceAPI, cfg.LastUpdate, &goconfluenceURL)
				if err != nil {
					log.Fatal(err)
				}
				for _, page := range allPages {
					if debugMode {
						fmt.Print(page.Summary + "\n\n")
						fmt.Print(page.Description)
					} else {
						if !hasExistingJiraIssue(page.Summary, cfg.JiraProjectKey, jiraClient) {
							issue := jira.Issue{
								Fields: &jira.IssueFields{
									Type:        jira.IssueType{Name: "Task"},
									Project:     jira.Project{Key: cfg.JiraProjectKey},
									Description: page.Description,
									Summary:     page.Summary,
									Labels:      cfg.JiraLabels,
								},
							}
							createdIssue, resp, err := jiraClient.Issue.Create(&issue)
							if err != nil {
								log.Printf("Unable to create Jira issue for %s \n %s \n", cfg.Name, err)
								log.Print(resp)
								continue
							}
							fmt.Printf("%s: %+v\n", createdIssue.Key, createdIssue.Self)
							log.Printf("Created Jira Issue '%s' in project: %s' \n", createdIssue.Key, cfg.JiraProjectKey)
						}
					}
				}
			}
			log.Printf("%s: job complete.", cfg.Name)
			log.Printf("%s: waiting for %s", cfg.Name, cfg.Interval)
			timer = time.NewTimer(intervalDuration)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {

	_ = kong.Parse(&args)

	if err := ValidateConfigPath(args.Config); err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig(args.Config)
	if err != nil {
		log.Fatal(err)
	}

	debugMode = cfg.Debug
	// goconfluence.SetDebug(debugMode)

	// Create URLs for API libraries
	baseURL, err = url.Parse(cfg.Atlassian.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	goconfluenceURL = *baseURL
	goconfluenceURL.Path = path.Join(goconfluenceURL.Path, "/wiki/rest/api")
	gojiraURL = *baseURL

	// Set up Jira client
	tp := jira.BasicAuthTransport{
		Username: args.AtlassianUsername,
		Password: args.AtlassianToken,
	}
	jiraClient, err = jira.NewClient(tp.Client(), gojiraURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// Set up Confluence client
	confluenceAPI, err = goconfluence.NewAPI(goconfluenceURL.String(), args.AtlassianUsername, args.AtlassianToken)
	if err != nil {
		log.Fatal(err)
	}

	osSignals = make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	if cfg.IssueConfigs != nil {
		wg.Add(len(cfg.IssueConfigs))
		for _, ic := range cfg.IssueConfigs {
			go issueRaiser(&ic)
		}
	}

	wg.Wait()

}
