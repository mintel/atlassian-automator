# Atlassian Automator

The Atlassian Automator checks last updated dates/times for pages in Confluence and raises JIRA tickets when they are older than a date set by the user. This may be extended to other use-cases in future (hence the generic name).

```
$ ./atlassian-automator -h
Usage: atlassian-automator --atlassian-token=STRING --atlassian-username=STRING

Flags:
  -h, --help                         Show context-sensitive help.
      --atlassian-token=STRING       Your Atlassian API token. Either the environment variable or the flag MUST be set
                                     ($ATLASSIAN_TOKEN).
      --atlassian-username=STRING    Your Atlassian user name. Either the environment variable or the flag MUST be set
                                     ($ATLASSIAN_USERNAME).
      --config="config.yaml"         Path to atlassian-automator config file ($CONFIG_FILE).
      --listen-address=":8000"       Address on which HTTP server will listen (for healthchecks and metrics)
                                     ($LISTEN_ADDRESS).
```

## Configuration

1. Ensure your `$ATLASSIAN_USERNAME` and `$ATLASSIAN_TOKEN` environment variables / command line arguments are set
2. Create a file called `config.yaml` and point to it using the `--config` argument:

```yaml
# Set to true if you want to see what issues would be raised in Jira but don't want to actually raise any
debug: false

# atlassian config goes here
atlassian:
  # baseURL is the base Atlassian URL for your organisation
  baseURL: https://${DOMAIN}.atlassian.net

# issues is a list of jobs that will create issues in Jira
issues:
  # name is a unique name for your job
- name: my-example-job
  # jiraProjectKey determines the project in which your issue(s) will be raised
  jiraProjectKey: EXAMPLE
  # jiraLabels is a list of labels you want to apply to any new issue(s)
  jiraLabels:
  - area/Documentation
  - area/Ops
  # interval determines how often this job will run. See https://pkg.go.dev/time#ParseDuration for valid input
  interval: 1m
  
  # lastUpdate is a dict which corresponds to the package of the same name. We may add more packages later that do other
  # things
  lastUpdate:
    # duration determines how old a page has to be before it is passed to the issue creator
    duration: 24h
    # parentPageID will restrict checks for updated pages to only those that are children of this one
    parentPageID: 1234567890
    # spaceKey is the key of the space in which the pages you want to check live
    spaceKey: MYSPACE
    # type is a comma-separated list of type(s) you want to check for updates. See 
    # https://developer.atlassian.com/cloud/confluence/rest/api-group-content---children-and-descendants/#api-wiki-rest-api-content-id-child-type-get # for the valid types.
    type: page
    # resultsLimit is the max number of results you want the API call to return
    resultsLimit: 1000
  # retryInterval is the length of time you wish to wait before retrying a job if it fails for some reason
  retryInterval: 1m
```
