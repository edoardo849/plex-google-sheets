package plex

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

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

// Do performs the authentication: if a cache login is found, it will load the user from that file
func (l Login) Do() (*User, error) {
	cacheFile, err := cache.FilePath(".credentials", url.QueryEscape("plex.tv-plex-google-sheets.json"))
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	u, err := userFromFile(cacheFile)
	if err != nil {
		log.Warn("Token not set: authenticating")
		promptLogin(&l)
		u, err = authenticate(&l)

		if err != nil {
			log.Fatal("Failed to authenticate to the Plex service")
		}

		saveUser(cacheFile, u)
		return u, nil
	}

	log.Debug("Loading account details")
	return u, nil
}

func promptLogin(l *Login) {
	fmt.Println("Enter Plex email: ")
	var email string
	if _, err := fmt.Scan(&email); err != nil {
		log.Fatalf("Unable to read the email %v", err)
	}
	l.Username = email

	fmt.Println("Enter Plex password: ")
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Couldn't parse the password %v", err)
	}
	l.Password = string(pwd)
}

func userFromFile(file string) (*User, error) {

	f, err := cache.Open(file)

	if err != nil {
		return nil, errors.Wrapf(err, "The cache file %s was not found", file)
	}
	u := &User{}
	err = json.NewDecoder(f).Decode(u)
	defer f.Close()
	log.Debug("User loaded from cache")
	return u, err

}

// saveToken uses a file path to create a file and store the
// token in it.
func saveUser(file string, u *User) {
	log.Debugf("Saving credential file to: %s\n", file)
	f, err := cache.OpenOrCreate(file)

	if err != nil {
		log.Fatalf("Unable to cache user: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(u)
}

func authenticate(l *Login) (*User, error) {

	req, err := http.NewRequest("POST", authURL, nil)

	if err != nil {
		return nil, errors.Wrap(err, "POST request building failed")
	}

	req.SetBasicAuth(l.Username, l.Password)

	addPlexHeaders(req)

	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {

		return nil, errors.Wrap(err, "API call failed")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.Wrapf(errors.New(http.StatusText(resp.StatusCode)), "Unexpected status code %d from the API", resp.StatusCode)
	}

	decoder := json.NewDecoder(io.ReadCloser(resp.Body))
	var a account
	err = decoder.Decode(&a)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode the response body")
	}

	u := a.User

	l.Token = a.User.Token
	log.WithFields(log.Fields{
		"email": u.Email,
		"thumb": u.Thumb,
		"ID":    u.ID,
	}).Info("User Details")

	return &u, nil
}

func loadAccountDetails(token string) (*User, error) {

	req, err := http.NewRequest("GET", accountInfoURL, nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = fmt.Sprintf("X-Plex-Token=%s", token)
	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	decoder := json.NewDecoder(io.ReadCloser(resp.Body))
	var a account
	err = decoder.Decode(&a)

	if err != nil {
		return nil, err
	}
	return &a.User, nil
}
