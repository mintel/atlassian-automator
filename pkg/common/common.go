package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CollectedData struct {
	Summary     string
	Description string
}

var (
	PromErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "atlassian_automator_errors_total",
			Help: "The number of errors encountered by the main package",
		},
		[]string{
			"package",
		})
)
