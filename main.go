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
	user, err := plex.NewUser(plex.SetUsername("edoardo849@gmail.com"), plex.SetPassword("petergower849"))

	if err != nil {
		e := errors.Wrap(err, "New user creation failed")
		log.Fatal(e)
	}

	log.WithFields(log.Fields{

		"ID":       user.ID,
		"Email":    user.Email,
		"Thumb":    user.Thumb,
		"Title":    user.Title,
		"Token":    user.Token,
		"Username": user.Username,
	}).Info("User")

}
