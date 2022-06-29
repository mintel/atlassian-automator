package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/alecthomas/kong"
	goconfluence "github.com/virtomize/confluence-go-api"
)

var args struct {
	Username string `env:"ATLASSIAN_USERNAME"`
	Token    string `env:"ATLASSIAN_TOKEN"`
}

func filterResults(ancestorID string, olderThan time.Duration, cs *goconfluence.ContentSearch) []goconfluence.Content {

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
				log.Fatal(err)
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
	return filtered
}

func main() {

	// goconfluence.SetDebug(true)
	_ = kong.Parse(&args)

	// Initialize a new api instance
	api, err := goconfluence.NewAPI("https://mintel.atlassian.net/wiki/rest/api", args.Username, args.Token)
	if err != nil {
		log.Fatal(err)
	}

	// Get content by query
	res, err := api.GetContent(goconfluence.ContentQuery{
		SpaceKey: "CI",
		Expand:   []string{"history.lastUpdated", "ancestors"},
		Type:     "page",
		Limit:    10000,
		// Ordering by lastUpdated.when is not supported by the API so we have to get EVERYTHING and then sort within Go :(
		// OrderBy: "history.lastUpdated.when desc",
	})
	if err != nil {
		log.Fatal(err)
	}

	var oneYear time.Duration = 48 * time.Hour
	allPages := filterResults("2950594561", oneYear, res)
	sort.SliceStable(allPages, func(i, j int) bool {
		return allPages[i].History.LastUpdated.When > allPages[j].History.LastUpdated.When
	})
	for _, page := range allPages {
		fmt.Printf("%+v\n", page.History.LastUpdated.When)
	}

}
