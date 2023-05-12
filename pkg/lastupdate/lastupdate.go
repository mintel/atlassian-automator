package lastupdate

import (
	"context"
	"fmt"

	"github.com/mintel/atlassian-automator/pkg/common"
	"github.com/mintel/atlassian-automator/pkg/confluence"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// const (
// 	pkg string = "lastupdate"
// )

var (
	PromLastUpdatePagesTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "atlassian_automator_lastupdate_pages_total",
			Help: "The number of pages monitored by lastupdate jobs",
		},
		[]string{"name"},
	)
)

type Config struct {
	Duration     string `yaml:"duration"`
	ParentPageID string `yaml:"parentPageID"`
	SpaceKey     string `yaml:"spaceKey"`
	Type         string `yaml:"type"`
	ResultsLimit int    `yaml:"resultsLimit"`
}

// findChildren takes a confluence.Pages object and iterates through it to find child pages of the provided
// parentPageID, then calls itself to find children of those children and so on.
func findChildren(allPages confluence.Pages, parentPageID string) confluence.Pages {
	var children confluence.Pages

	for _, page := range allPages.Results {
		if page.ParentID == parentPageID {
			children.Results = append(children.Results, page)
		}
	}

	if len(children.Results) > 0 {
		for _, page := range children.Results {
			next := findChildren(allPages, page.ID)
			children.Results = append(children.Results, next.Results...)
		}
	}
	return children
}

// func filterResults(ancestorID string, olderThan time.Duration, allPages confluence.Pages) error {

// 	var filtered confluence.Pages

// 	for _, result := range allPages.Results {

// 		// Filter out anything not under the specified page
// 		wanted := false
// 		for _, ancestor := range result. {
// 			if ancestor.ID == ancestorID {
// 				wanted = true
// 				break
// 			}
// 		}

// 		// Filter out anything not older than the provided date
// 		if wanted {
// 			layout := "2006-01-02T15:04:05.000Z"
// 			lastUpdatedTime, err := time.Parse(layout, result.History.LastUpdated.When)
// 			if err != nil {
// 				common.PromErrors.WithLabelValues(pkg).Inc()
// 				return nil, err
// 			}
// 			oldestAllowedTime := time.Now().Add(-olderThan)
// 			if !lastUpdatedTime.Before(oldestAllowedTime) {
// 				wanted = false
// 			}
// 		}

// 		if wanted {
// 			filtered = append(filtered, result)
// 		}
// 	}
// 	return filtered, nil
// }

// func Run(jobName string, api goconfluence.API, cfg Config, baseURL *url.URL) ([]common.CollectedData, error) {

// 	var collectedData []common.CollectedData

// 	// Get content by query
// 	res, err := api.GetContent(goconfluence.ContentQuery{
// 		SpaceKey: cfg.SpaceKey,
// 		Expand:   []string{"history.lastUpdated", "ancestors"},
// 		Type:     cfg.Type,
// 		Limit:    cfg.ResultsLimit,
// 		OrderBy:  "history.lastUpdated.when desc",
// 	})
// 	if err != nil {
// 		common.PromErrors.WithLabelValues(pkg).Inc()
// 		log.Printf("%s: error getting content from confluence: %s", jobName, err)
// 		return nil, err
// 	}

// 	duration, err := time.ParseDuration(cfg.Duration)
// 	if err != nil {
// 		common.PromErrors.WithLabelValues(pkg).Inc()
// 		log.Printf("%s: couldn't parse duration %s: %s", jobName, cfg.Duration, err)
// 		return nil, err
// 	}

// 	allPages, err := filterResults(cfg.ParentPageID, duration, res)
// 	if err != nil {
// 		common.PromErrors.WithLabelValues(pkg).Inc()
// 		log.Printf("%s: couldn't filter results: %s", jobName, err)
// 		return nil, err
// 	}
// 	log.Printf("%s: %+v pages found older than %s", jobName, len(allPages), cfg.Duration)
// 	sort.SliceStable(allPages, func(i, j int) bool {
// 		return allPages[i].History.LastUpdated.When > allPages[j].History.LastUpdated.When
// 	})

// 	for _, page := range allPages {
// 		pageURL := *baseURL
// 		pageURL.Path = path.Join(pageURL.Path, page.Links.TinyUI)
// 		collectedData = append(collectedData,
// 			common.CollectedData{
// 				Summary:     fmt.Sprintf("\"%s\" has not been updated since %s", page.Title, page.History.LastUpdated.When),
// 				Description: fmt.Sprintf("Page ID: %s\nPage Title: %s\nLast Updated: %s\nURL: %s\n\nPlease review this page and update any out-of-date information if required.\n\n", page.ID, page.Title, page.History.LastUpdated.When, pageURL.String()),
// 			},
// 		)
// 	}
// 	return collectedData, nil
// }

// getSpaceIDFromKey gets the ID of a space using the provided key
func getSpaceIDFromKey(ctx context.Context, key string) (string, error) {

	spaces, _, err := common.ConfluenceClient.Space.GetSpaces(ctx, &confluence.GetSpacesOptions{
		Keys:                  []string{key},
		SerializeIDsAsStrings: true,
	})
	if err != nil {
		return "", err
	}
	if len(spaces.Results) == 0 {
		return "", fmt.Errorf("no spaces found with key %s", key)
	}
	if len(spaces.Results) > 1 {
		return "", fmt.Errorf("more than one space found with key %s", key)
	}
	return spaces.Results[0].ID, nil
}

func Run(ctx context.Context, jobName string, cfg Config) ([]common.CollectedData, error) {
	var collectedData []common.CollectedData

	// Get the space's ID using the key provided in the config
	spaceID, err := getSpaceIDFromKey(ctx, cfg.SpaceKey)
	if err != nil {
		return nil, err
	}

	// Get all pages in the space
	pages, _, err := common.ConfluenceClient.Page.GetPagesInSpace(ctx, spaceID, &confluence.GetPagesInSpaceOptions{
		Limit:                 cfg.ResultsLimit,
		SerializeIDsAsStrings: true,
	})
	if err != nil {
		return nil, err
	}

	// Filter the list so it only includes children of the parent page ID provided in the config
	filtered := findChildren(*pages, cfg.ParentPageID)
	PromLastUpdatePagesTotal.WithLabelValues(jobName).Set(float64(len(filtered.Results)))
	fmt.Print(filtered)

	// Find pages not updated since the provided duration
	// duration, err := time.ParseDuration(cfg.Duration)
	// if err != nil {
	// 	common.PromErrors.WithLabelValues(pkg).Inc()
	// 	log.Printf("%s: couldn't parse duration %s: %s", jobName, cfg.Duration, err)
	// 	return nil, err
	// }
	return collectedData, nil
}
