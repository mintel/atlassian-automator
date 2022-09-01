package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mintel/atlassian-automator/pkg/common"
	"github.com/mintel/atlassian-automator/pkg/issueraiser"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

var (
	baseURL *url.URL
)

var args struct {
	AtlassianToken    string `env:"ATLASSIAN_TOKEN" required:"" help:"Your Atlassian API token. Either the environment variable or the flag MUST be set."`
	AtlassianUsername string `env:"ATLASSIAN_USERNAME" required:"" help:"Your Atlassian user name. Either the environment variable or the flag MUST be set."`
	Config            string `env:"CONFIG_FILE" default:"config.yaml" type:"path" help:"Path to atlassian-automator config file."`
	ListenAddress     string `env:"LISTEN_ADDRESS" default:":8000" help:"Address on which HTTP server will listen (for healthchecks and metrics)."`
}

type Config struct {
	Atlassian struct {
		BaseURL string `yaml:"baseURL"`
	} `yaml:"atlassian"`
	IssueConfigs []issueraiser.Config `yaml:"issues"`
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

	// Setup Atlassian clients
	baseURL, err = url.Parse(cfg.Atlassian.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	common.AtlassianSetup(baseURL, args.AtlassianUsername, args.AtlassianToken)

	// Set up OS signal notifications
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start job goroutines
	var wg sync.WaitGroup
	if cfg.IssueConfigs != nil {
		wg.Add(len(cfg.IssueConfigs))
		for _, ic := range cfg.IssueConfigs {
			go issueraiser.Run(ctx, &wg, &ic)
		}
	}

	// Start HTTP server goroutine for healthchecks
	httpServer := http.Server{
		Addr: args.ListenAddress,
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

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
