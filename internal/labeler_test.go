package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/github"
)

// GithubRepoService is an interface defining the Wrapper Interface
// needed to test the github client.
type MockGithubIssuesService struct {
	Err      error
	Response *github.Response
	Labels   []*github.Label
	Owner    string
	Repo     string
	Number   int
}

func (mis *MockGithubIssuesService) AddLabelsToIssue(ctx context.Context, owner string, repo string, number int, labels []string) ([]*github.Label, *github.Response, error) {
	if owner != mis.Owner {
		return nil, nil, fmt.Errorf("owner did not equal expected owner: %s was: %s", mis.Owner, owner)
	}
	if repo != mis.Repo {
		return nil, nil, fmt.Errorf("repo did not equal expected repo: %s was: %s", mis.Repo, repo)
	}
	if number != mis.Number {
		return nil, nil, fmt.Errorf("number did not equal expected number: %d was: %d", mis.Number, number)
	}
	return mis.Labels, mis.Response, mis.Err
}

func TestAddLabel(t *testing.T) {
	mis := new(MockGithubIssuesService)
	mis.Owner = "owner"
	mis.Repo = "repo"
	mis.Number = 1
	labeler := GithubLabeler{issues: mis, ctx: context.Background()}
	err := labeler.AddLabel("owner", "repo", 1)
	if err != nil {
		t.Fatal(err)
	}
}
