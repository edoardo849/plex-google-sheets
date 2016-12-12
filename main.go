package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/edoardo849/plex-google-sheets/plex"
	"github.com/pkg/errors"
)

func init() {

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

}

func main() {

	pwd := os.Getenv("PLEX_PASSWORD")
	uname := os.Getenv("PLEX_USERNAME")
	token := os.Getenv("PLEX_TOKEN")

	login := plex.Login{
		Username: uname,
		Password: pwd,
		Token:    token,
	}

	u, err := login.Do()

	if err != nil {
		e := errors.Wrap(err, "New user creation failed")
		log.Fatal(e)
	}

	log.Infof("User %s authenticated", u.Username)
}
