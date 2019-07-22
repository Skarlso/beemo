package internal

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// NewGithubClient creates a wrapper around the github client.
func NewGithubClient(httpClient *http.Client) *github.Client {
	return github.NewClient(httpClient)
}

// AddLabel adds certain labels to a PR given an issueUrl point to the PR.
func AddLabel(owner, repo string, number int) error {
	LogDebug("Adding hook for issue number: ", number)
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if token == "" {
		return errors.New("GITHUB_OAUTH_TOKEN missing")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := NewGithubClient(tc)
	LogDebug("client created: ", client)
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, number, []string{"good first issue"})
	if err != nil {
		return err
	}
	return nil
}
