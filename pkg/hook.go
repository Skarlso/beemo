package pkg

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Skarlso/acquia-beemo/internal"
	"github.com/labstack/echo/v4"
)

// Hook represent a github based webhook context.
type Hook struct {
	Signature string
	Event     string
	ID        string
	Payload   []byte
}

// Owner is the owner of the repository.
type Owner struct {
	Login string `json:"login"`
}

// Repository is information about the repository attached to the PR.
type Repository struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}

// Payload contains information about the event like, user, commit id and so on.
type Payload struct {
	Action string     `json:"action"`
	Number int        `json:"number"`
	Repo   Repository `json:"repository"`
}

const (
	// Ping event name.
	Ping = "ping"
	// PR event name.
	PR = "pull_request"
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

func parseRequest(secret []byte, req *http.Request) (Hook, error) {
	h := Hook{}

	if h.Signature = req.Header.Get("x-hub-signature"); len(h.Signature) == 0 {
		return Hook{}, errors.New("no signature")
	}

	if h.Event = req.Header.Get("x-github-event"); len(h.Event) == 0 {
		return Hook{}, errors.New("no event")
	}

	if h.Event != PR {
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

// GitWebHook creates a hook handler with an injected labeler.
func GitWebHook(labeler internal.Labeler) echo.HandlerFunc {
	return func(c echo.Context) error {
		internal.LogDebug("[DEBUG] Received request from: ", c.Request().UserAgent())
		secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
		if secret == "" {
			return c.String(http.StatusBadRequest, "GITHUB_WEBHOOK_SECRET is empty")
		}

		h, err := parseRequest([]byte(secret), c.Request())
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if h.Event == Ping {
			return c.NoContent(http.StatusOK)
		}

		c.Request().Header.Set("Content-type", "application/json")

		p := new(Payload)
		if err := json.Unmarshal(h.Payload, p); err != nil {
			return c.String(http.StatusBadRequest, "error in unmarshalling json payload")
		}

		if p.Action != "opened" {
			return c.String(http.StatusOK, "skipped; status was not opened but: "+p.Action)
		}

		err = labeler.AddLabel(p.Repo.Owner.Login, p.Repo.Name, p.Number)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, "successfully processed event")
	}
}
