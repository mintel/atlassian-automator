package lastupdate

import (
	"context"
	"fmt"
	"log"
	"path"
	"sort"
	"time"

	"github.com/mintel/atlassian-automator/pkg/common"
	"github.com/mintel/atlassian-automator/pkg/confluence"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	pkg string = "lastupdate"
)

var (
	PromLastUpdatePagesTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "atlassian_automator_lastupdate_pages_total",
			Help: "The number of pages monitored by lastupdate jobs",
		},
		[]string{"job_name"},
	)
)

type Config struct {
	Duration     string `yaml:"duration"`
	ParentPageID string `yaml:"parentPageID"`
	SpaceKey     string `yaml:"spaceKey"`
	Type         string `yaml:"type"`
	ResultsLimit int    `yaml:"resultsLimit"`
}

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

// findChildren takes a confluence.Pages object and iterates through it to find child pages of the provided
// parentPageID, then calls itself to find children of those children and so on.
func findChildren(pages confluence.Pages, parentPageID string) confluence.Pages {
	var children confluence.Pages

	for _, page := range pages.Results {
		if page.ParentID == parentPageID {
			children.Results = append(children.Results, page)
		}
	}

	if len(children.Results) > 0 {
		for _, page := range children.Results {
			next := findChildren(pages, page.ID)
			children.Results = append(children.Results, next.Results...)
		}
	}
	return children
}

func findNotUpdated(pages confluence.Pages, duration time.Duration) (*confluence.Pages, error) {
	var output confluence.Pages
	timeLayout := "2006-01-02T15:04:05.000Z"
	oldestAllowedTime := time.Now().Add(-duration)

	for _, page := range pages.Results {
		lastUpdatedTime, err := time.Parse(timeLayout, page.Version.CreatedAt)
		if err != nil {
			return &confluence.Pages{}, nil
		}
		if lastUpdatedTime.Before(oldestAllowedTime) {
			output.Results = append(output.Results, page)
		}
	}
	return &output, nil
}

// Run gets metadata for all the pages in a Confluence space and then filters out the pages that 1) are not descendents
// of the provided parent page and 2) have been updated within the specified time period. Doing it this way (rather than
// using the Children API endpoint) dramatically reduces the number of API calls and the total time taken.
func Run(ctx context.Context, jobName string, cfg Config) ([]common.CollectedData, error) {
	var collectedData []common.CollectedData

	// Get the space's ID using the key provided in the config
	spaceID, err := getSpaceIDFromKey(ctx, cfg.SpaceKey)
	if err != nil {
		common.PromErrors.WithLabelValues(pkg).Inc()
		log.Printf("%s: couldn't get space ID from key %s: %s", jobName, cfg.SpaceKey, err)
		return nil, err
	}

	// Get all pages in the space
	pagesInSpace, _, err := common.ConfluenceClient.Page.GetPagesInSpace(ctx, spaceID, &confluence.GetPagesInSpaceOptions{
		Limit:                 cfg.ResultsLimit,
		SerializeIDsAsStrings: true,
	})
	if err != nil {
		common.PromErrors.WithLabelValues(pkg).Inc()
		log.Printf("%s: couldn't get pages in space %s: %s", jobName, cfg.SpaceKey, err)
		return nil, err
	}

	// Filter the list so it only includes children of the parent page ID provided in the config
	children := findChildren(*pagesInSpace, cfg.ParentPageID)
	PromLastUpdatePagesTotal.WithLabelValues(jobName).Set(float64(len(children.Results)))

	// Find pages not updated since the provided duration
	duration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		common.PromErrors.WithLabelValues(pkg).Inc()
		log.Printf("%s: couldn't parse duration %s: %s", jobName, cfg.Duration, err)
		return nil, err
	}
	notUpdated, err := findNotUpdated(children, duration)
	if err != nil {
		common.PromErrors.WithLabelValues(pkg).Inc()
		log.Printf("%s: couldn't parse duration %s: %s", jobName, cfg.Duration, err)
		return nil, err
	}
	log.Printf("%s: %+v pages found older than %s", jobName, len(notUpdated.Results), cfg.Duration)
	sort.SliceStable(notUpdated.Results, func(i, j int) bool {
		return notUpdated.Results[i].Version.CreatedAt > notUpdated.Results[j].Version.CreatedAt
	})

	baseURL := common.ConfluenceClient.BaseURL
	for _, page := range notUpdated.Results {
		pageURL := path.Join(baseURL.String(), fmt.Sprintf("/wiki/pages/viewpage.action?pageId=%s", page.ID))
		collectedData = append(collectedData,
			common.CollectedData{
				Summary:     fmt.Sprintf("\"%s\" has not been updated since %s", page.Title, page.Version.CreatedAt),
				Description: fmt.Sprintf("Page ID: %s\nPage Title: %s\nLast Updated: %s\nURL: %s\n\nPlease review this page and update any out-of-date information if required.\n\n", page.ID, page.Title, page.Version.CreatedAt, pageURL),
			},
		)
	}

	return collectedData, nil
}
