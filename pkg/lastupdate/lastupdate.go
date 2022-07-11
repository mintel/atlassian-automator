package lastupdate

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"sort"
	"time"

	"github.com/mintel/atlassian-automator/pkg/common"
	goconfluence "github.com/virtomize/confluence-go-api"
)

type Config struct {
	Duration     string `yaml:"duration"`
	ParentPageID string `yaml:"parentPageID"`
	SpaceKey     string `yaml:"spaceKey"`
	Type         string `yaml:"type"`
	ResultsLimit int    `yaml:"resultsLimit"`
}

func filterResults(ancestorID string, olderThan time.Duration, cs *goconfluence.ContentSearch) ([]goconfluence.Content, error) {

	var filtered []goconfluence.Content

	for _, result := range cs.Results {

		// Filter out anything not under the specified page
		wanted := false
		for _, ancestor := range result.Ancestors {
			if ancestor.ID == ancestorID {
				wanted = true
				break
			}
		}

		// Filter out anything not older than the provided date
		if wanted {
			layout := "2006-01-02T15:04:05.000Z"
			lastUpdatedTime, err := time.Parse(layout, result.History.LastUpdated.When)
			if err != nil {
				return nil, err
			}
			oldestAllowedTime := time.Now().Add(-olderThan)
			if !lastUpdatedTime.Before(oldestAllowedTime) {
				wanted = false
			}
		}

		if wanted {
			filtered = append(filtered, result)
		}
	}
	return filtered, nil
}

func Run(api goconfluence.API, cfg Config, baseURL *url.URL) ([]common.CollectedData, error) {

	var collectedData []common.CollectedData

	// Get content by query
	res, err := api.GetContent(goconfluence.ContentQuery{
		SpaceKey: cfg.SpaceKey,
		Expand:   []string{"history.lastUpdated", "ancestors"},
		Type:     cfg.Type,
		Limit:    cfg.ResultsLimit,
		// Ordering by lastUpdated.when is not supported by the API so we have to get EVERYTHING and then sort within Go :(
		// OrderBy: "history.lastUpdated.when desc",
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	allPages, err := filterResults(cfg.ParentPageID, duration, res)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Printf("%+v pages found older than %s", len(allPages), cfg.Duration)
	sort.SliceStable(allPages, func(i, j int) bool {
		return allPages[i].History.LastUpdated.When > allPages[j].History.LastUpdated.When
	})

	for _, page := range allPages {
		pageURL := *baseURL
		pageURL.Path = path.Join(pageURL.Path, page.Links.TinyUI)
		collectedData = append(collectedData,
			common.CollectedData{
				Summary:     fmt.Sprintf("\"%s\" has not been updated since %s", page.Title, page.History.LastUpdated.When),
				Description: fmt.Sprintf("Page ID: %s\nPage Title: %s\nLast Updated: %s\nURL: %s\n\nPlease review this page and update any out-of-date information if required.\n\n", page.ID, page.Title, page.History.LastUpdated.When, pageURL.String()),
			},
		)
	}
	return collectedData, nil
}
