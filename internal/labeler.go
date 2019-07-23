package internal

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Labeler defines a service which can add a label to a pull request.
type Labeler interface {
	AddLabel(string, string, int) error
}

// GithubLabeler is a concrete implementation of a Labeler that provides
// Github labelling capabilities.
type GithubLabeler struct {
	issues IssuesService
	ctx    context.Context
}

// NewGithubLabeler returns a concrete github label implementation.
func NewGithubLabeler() Labeler {
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return GithubLabeler{issues: client.Issues, ctx: ctx}
}

// IssuesService is an interface defining the Wrapper Interface
// needed to test the github client.
type IssuesService interface {
	AddLabelsToIssue(ctx context.Context, owner string, repo string, number int, labels []string) ([]*github.Label, *github.Response, error)
}

// AddLabel adds certain labels to a PR given an issueUrl point to the PR.
func (gt GithubLabeler) AddLabel(owner, repo string, number int) error {
	LogDebug("Adding hook for issue number: ", number)
	_, _, err := gt.issues.AddLabelsToIssue(gt.ctx, owner, repo, number, []string{"good first issue"})
	if err != nil {
		return err
	}
	log.Println("Label successfully added to PR.")
	return nil
}
