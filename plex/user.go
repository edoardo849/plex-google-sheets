package plex

import (
	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	accountInfoURL = "https://plex.tv/users/account.json"
	authURL        = "https://plex.tv/users/sign_in.json"
)

// Login is the information required to perform an authentication against the API
type Login struct {
	Username string
	Password string
	Token    string
}

type account struct {
	User         User         `json:"user"`
	Subscription Subscription `json:"subscription"`
	Entitlements []string     `json:"entitlements"`
}

// Subscription is...
type Subscription struct {
	Active   bool     `json:"active"`
	Status   string   `json:"status"`
	Plan     string   `json:"plan"`
	Features []string `json:"features"`
}

// User is the object returned by the Plex login
type User struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	JoinedAt string `json:"joined_at"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Thumb    string `json:"thumb"`
	Title    string `json:"title"`
	Token    string `json:"authToken"`
}

// Do performs the authentication
func (l Login) Do() (User, error) {

	if len(l.Token) == 0 {
		log.Debug("Token not set: authenticating")
		return authenticate(l.Username, l.Password)
	}
	log.Debug("Token set: loading account details")
	return loadAccountDetails(l.Token)
}

func authenticate(uname, pwd string) (User, error) {
	u := User{}

	req, err := http.NewRequest("POST", authURL, nil)

	if err != nil {
		return u, errors.Wrap(err, "POST request building failed")
	}

	req.SetBasicAuth(uname, pwd)

	addPlexHeaders(req)

	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {

		return u, errors.Wrap(err, "API call failed")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return u, errors.Wrapf(errors.New(http.StatusText(resp.StatusCode)), "Unexpected status code %d from the API", resp.StatusCode)
	}

	decoder := json.NewDecoder(io.ReadCloser(resp.Body))
	var a account
	err = decoder.Decode(&a)

	if err != nil {
		return u, errors.Wrap(err, "Failed to decode the response body")
	}
	u = a.User
	log.WithFields(log.Fields{
		"email": u.Email,
		"thumb": u.Thumb,
		"ID":    u.ID,
	}).Info("User Details")

	return u, nil
}

// TODO find a usage for this.
func loadAccountDetails(token string) (User, error) {
	u := User{}

	req, err := http.NewRequest("GET", accountInfoURL, nil)

	if err != nil {
		return u, err
	}

	req.URL.RawQuery = fmt.Sprintf("X-Plex-Token=%s", token)

	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {
		return u, err
	}

	if resp.StatusCode != 200 {
		return u, errors.New(http.StatusText(resp.StatusCode))
	}

	decoder := json.NewDecoder(io.ReadCloser(resp.Body))
	var a account
	err = decoder.Decode(&a)

	if err != nil {
		return u, err

	}
	return a.User, nil
}
