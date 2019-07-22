package pkg

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// Hook represent a github based webhook context.
type Hook struct {
	Signature string
	Event     string
	ID        string
	Payload   []byte
}

// Repository contains information about the repository. All we care about
// here are the possible urls for identification.
type Repository struct {
	GitURL  string `json:"git_url"`
	SSHURL  string `json:"ssh_url"`
	HTMLURL string `json:"html_url"`
}

// Payload contains information about the event like, user, commit id and so on.
// All we care about for the sake of identification is the repository.
type Payload struct {
	Repo Repository `json:"repository"`
}

const (
	// Ping event name.
	Ping = "ping"
	// Push event name.
	Push = "push"
)

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	_, _ = computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {
	signaturePrefix := "sha1="
	signatureLength := 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	_, _ = hex.Decode(actual, []byte(signature[5:]))
	expected := signBody(secret, body)
	return hmac.Equal(expected, actual)
}

func parse(secret []byte, req *http.Request) (Hook, error) {
	h := Hook{}

	if h.Signature = req.Header.Get("x-hub-signature"); len(h.Signature) == 0 {
		return Hook{}, errors.New("no signature")
	}

	if h.Event = req.Header.Get("x-github-event"); len(h.Event) == 0 {
		return Hook{}, errors.New("no event")
	}

	if h.Event != Push {
		if h.Event == Ping {
			return Hook{Event: "ping"}, nil
		}
		return Hook{}, errors.New("invalid event")
	}

	if h.ID = req.Header.Get("x-github-delivery"); len(h.ID) == 0 {
		return Hook{}, errors.New("no event id")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return Hook{}, err
	}

	if !verifySignature(secret, h.Signature, body) {
		return Hook{}, errors.New("invalid signature")
	}

	h.Payload = body
	return h, err
}

// GitWebHook handles callbacks from GitHub's webhook system.
func GitWebHook(c echo.Context) error {
	LogDebug("[DEBUG] Received request from: ", c.Request().UserAgent())
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		return c.String(http.StatusBadRequest, "GITHUB_WEBHOOK_SECRET is empty")
	}

	h, err := parse([]byte(secret), c.Request())

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if h.Event == "ping" {
		return c.NoContent(http.StatusOK)
	}

	c.Request().Header.Set("Content-type", "application/json")

	p := new(Payload)
	if err := json.Unmarshal(h.Payload, p); err != nil {
		return c.String(http.StatusBadRequest, "error in unmarshalling json payload")
	}

	log.Printf("Received: GitURL: %s, SSHURL: %s\n", p.Repo.GitURL, p.Repo.SSHURL)

	return c.String(http.StatusOK, "successfully processed event")
}
