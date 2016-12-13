package plex

import (
	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/edoardo849/plex-google-sheets/cache"

	"github.com/pkg/errors"
)

const (
	accountInfoURL = "https://plex.tv/users/account.json"
	authURL        = "https://plex.tv/users/sign_in.json"
)

// Login is the information required to perform an authentication against the API
type Login struct {
	Username string `json:"-"`
	Password string `json:"-"`
	Token    string `json:"authToken"`
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
	cacheFile, err := cache.FilePath(".credentials", url.QueryEscape("plex.tv-plex-google-sheets.json"))
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Warn("Token not set: authenticating")
		u, err := authenticate(&l)

		if err != nil {
			log.Fatal("Failed to authenticate to the Plex service")
		}

		saveToken(cacheFile, l)
		return u, nil
	}
	log.Debug("Token set: loading account details")
	return loadAccountDetails(tok)
}

func tokenFromFile(file string) (string, error) {
	f, err := cache.Open(file)

	if err != nil {
		return "", err
	}
	l := &Login{}
	err = json.NewDecoder(f).Decode(l)
	defer f.Close()

	return l.Token, err

}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, l Login) {
	log.Debugf("Saving credential file to: %s\n", file)
	f, err := cache.OpenOrCreate(file)

	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(l)
}

func authenticate(l *Login) (User, error) {
	u := User{}

	req, err := http.NewRequest("POST", authURL, nil)

	if err != nil {
		return u, errors.Wrap(err, "POST request building failed")
	}

	req.SetBasicAuth(l.Username, l.Password)

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

	l.Token = a.User.Token
	log.WithFields(log.Fields{
		"email": u.Email,
		"thumb": u.Thumb,
		"ID":    u.ID,
	}).Info("User Details")

	return u, nil
}

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
