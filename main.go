package main

import (
	"context"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/edoardo849/plex-google-sheets/gsheets"
	"github.com/edoardo849/plex-google-sheets/plex"
	"github.com/pkg/errors"
)

func init() {

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

}

// TODO 2 binaries: one for setup and the other for executing everything
func main() {
	ctx := context.Background()

	gConf := os.Getenv("GSHEETS_CONF")

	login := plex.Login{}

	u, err := login.Do()

	if err != nil {
		e := errors.Wrap(err, "New user creation failed")
		log.Fatal(e)
	}

	log.Infof("User %s authenticated", u.Username)

	srv, err := gsheets.NewService(ctx, gConf, gsheets.DefaultClientFactory)
	if err != nil {
		log.Fatal(err)
	}
	gsheets.TestWrite(srv)
	// Now, spawn 2 goroutines, one that sends data from plex and one that receives data and sets into gsheets

	// Prepare gsheets: open sheet
	// Create sheets for movies and series (or do nothing if sheets already exists)
	// Prepare data and send batchUpdate in batches of ~100 or more movies.

}
