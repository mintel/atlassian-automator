package issueraiser

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/mintel/atlassian-automator/pkg/common"
	"github.com/mintel/atlassian-automator/pkg/lastupdate"
)

const (
	pkg string = "issueraiser"
)

type Config struct {
	Interval       string            `yaml:"interval"`
	JiraLabels     []string          `yaml:"jiraLabels"`
	JiraProjectKey string            `yaml:"jiraProjectKey"`
	LastUpdate     lastupdate.Config `yaml:"lastUpdate"`
	Name           string            `yaml:"name"`
	RetryInterval  string            `yaml:"retryInterval"`
}

func hasExistingJiraIssue(itemTitle string, projectKey string, jiraClient *jira.Client) (bool, error) {
	// Escape quotes in the title so its parsed correctly by Jira's JQL parser
	itemTitle = strings.ReplaceAll(itemTitle, `"`, `\\\"`)
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
	jiraIssue, jiraResponse, err := common.JiraClient.Issue.Create(&issue)
	if err != nil {
		return nil, nil, err
	}
	return jiraIssue, jiraResponse, nil
}

func lastUpdateRaiser(ctx context.Context, cfg *Config) {

	var allPages []common.CollectedData

	retryIntervalDuration, err := time.ParseDuration(cfg.RetryInterval)
	if err != nil {
		log.Fatal(err)
	}

	for {
		allPages, err = lastupdate.Run(ctx, cfg.Name, cfg.LastUpdate)
		if err != nil {
			log.Printf("retrying in %s", cfg.RetryInterval)
			timer := time.NewTimer(retryIntervalDuration)
			<-timer.C
		} else {
			break
		}
	}

	for _, page := range allPages {
		exists, err := hasExistingJiraIssue(page.Summary, cfg.JiraProjectKey, common.JiraClient)
		if err != nil {
			common.PromErrors.WithLabelValues(pkg).Inc()
			log.Print(err)
			break
		}
		if !exists {
			log.Printf("%s: creating issue for \"%s\"", cfg.Name, page.Summary)
			jiraIssue, _, err := raiseIssue(&page, cfg.JiraProjectKey, cfg.JiraLabels)
			if err != nil {
				common.PromErrors.WithLabelValues(pkg).Inc()
				log.Print(err)
				break
			} else {
				log.Printf("%s: issue created for \"%s\": %s", cfg.Name, page.Summary, jiraIssue.Key)
			}
		} else {
			log.Printf("%s: issue already exists for \"%s\"", cfg.Name, page.Summary)
		}
	}
}

func Run(ctx context.Context, wg *sync.WaitGroup, cfg *Config) {
	defer wg.Done()
	intervalDuration, err := time.ParseDuration(cfg.Interval)
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
			// Check for lastupdate config (not empty)
			if cfg.LastUpdate != (lastupdate.Config{}) {
				lastUpdateRaiser(ctx, cfg)
			}
			log.Printf("%s: job complete.", cfg.Name)
			log.Printf("%s: waiting for %s", cfg.Name, cfg.Interval)
			timer = time.NewTimer(intervalDuration)
		}
	}
}
