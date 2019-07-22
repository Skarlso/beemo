package internal

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubRepoService is an interface defining the Wrapper Interface
// needed to test the github client.
type IssuesService interface {
	AddLabelsToIssue(ctx context.Context, owner string, repo string, number int, labels []string) ([]*github.Label, *github.Response, error)
}

// GithubClient is a client that has the ability to replace the actual
// git client.
type GithubClient struct {
	Issues IssuesService
	*github.Client
}

// MockClient is an exported mock client that can be used to mock the behavior of AddLabel.
var MockIssueService IssuesService

// NewGithubClient creates a wrapper around the github client. This is
// needed in order to decouple gaia from github client to be
// able to unit test createGithubWebhook and ultimately have
// the ability to replace github with anything else.
func NewGithubClient(httpClient *http.Client) GithubClient {
	if MockIssueService != nil {
		return GithubClient{
			Issues: MockIssueService,
		}
	}
	githubClient := github.NewClient(httpClient)

	return GithubClient{
		Issues: githubClient.Issues,
	}
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
	log.Println("Label successfully added to PR.")
	return nil
}
