package pkg

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/github"

	"github.com/Skarlso/acquia-beemo/internal"

	"github.com/labstack/echo/v4"
)

// GithubRepoService is an interface defining the Wrapper Interface
// needed to test the github client.
type MockGithubIssuesService struct {
	Err      error
	Response *github.Response
	Labels   []*github.Label
	Owner    string
	Repo     string
}

func (mis *MockGithubIssuesService) AddLabelsToIssue(ctx context.Context, owner string, repo string, number int, labels []string) ([]*github.Label, *github.Response, error) {
	if owner != mis.Owner {
		return nil, nil, errors.New("owner did not equal expected owner: was: " + owner)
	}
	if repo != mis.Repo {
		return nil, nil, errors.New("repo did not equal expected repo: was: " + repo)
	}
	return mis.Labels, mis.Response, mis.Err
}

func TestGitWebHook(t *testing.T) {
	_ = os.Setenv("GITHUB_WEBHOOK_SECRET", "superawesomesecretgithubpassword")
	e := echo.New()
	mis := new(MockGithubIssuesService)
	mis.Owner = "Codertocat"
	mis.Repo = "Hello-World"
	internal.MockIssueService = mis
	t.Run("successfully extracting PR information from payload", func(t *testing.T) {
		payload, _ := ioutil.ReadFile(filepath.Join("test_data", "payload_sample.json"))
		req := httptest.NewRequest(echo.POST, "/githook", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		// Use https://www.freeformatter.com/hmac-generator.html#ad-output for example
		// to calculate a new sha if the fixture would change. The payload has to be
		// one line and no linebreak at the end.
		req.Header.Set("x-hub-signature", "sha1=21a102e67d5f897b97a96af53af6a13d9df038ff")
		req.Header.Set("x-github-event", "pull_request")
		req.Header.Set("X-github-delivery", "1234asdf")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = GitWebHook(c)

		if rec.Code != http.StatusOK {
			body, _ := ioutil.ReadAll(rec.Body)
			log.Println("Body: ", string(body))
			t.Fatalf("want response code %v got %v", http.StatusOK, rec.Code)
		}
	})
	t.Run("bad request on none pull_request type events", func(t *testing.T) {
		payload, _ := ioutil.ReadFile(filepath.Join("test_data", "payload_sample.json"))
		req := httptest.NewRequest(echo.POST, "/githook", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		// Use https://www.freeformatter.com/hmac-generator.html#ad-output for example
		// to calculate a new sha if the fixture would change. The payload has to be
		// one line and no linebreak at the end.
		req.Header.Set("x-hub-signature", "sha1=21a102e67d5f897b97a96af53af6a13d9df038ff")
		req.Header.Set("x-github-event", "push")
		req.Header.Set("X-github-delivery", "1234asdf")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = GitWebHook(c)

		if rec.Code != http.StatusBadRequest {
			body, _ := ioutil.ReadAll(rec.Body)
			log.Println("Body: ", string(body))
			t.Fatalf("want response code %v got %v", http.StatusBadRequest, rec.Code)
		}
	})
}
