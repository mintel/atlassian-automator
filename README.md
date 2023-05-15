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
1. Copy `config.example.yaml` to `config.yaml` and edit as appropriate
1. `go build`
1. `./atlassian-automator`

## Development

To develop and test this you'll need the following environment variables set:

  * `ATLASSIAN_USERNAME` (your email address)
  * `ATLASSIAN_TOKEN` (the Atlassian token you created above):

If you don't have an existing org you can use, you can get one [here](https://www.atlassian.com/try/cloud/signup?product=confluence.ondemand,jira-software.ondemand,jira-servicedesk.ondemand,jira-core.ondemand&developer=true)) with:


## Thanks

[go-jira](https://github.com/andygrunwald/go-jira/tree/main) was an important reference for the `confluence` v2 API package and we have copied a lot of its code either as-is or used it as inspiration.
