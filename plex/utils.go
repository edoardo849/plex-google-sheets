package plex

import (
	"net/http"
	"runtime"
)

var (
	plexHeaders = map[string]string{
		"X-Plex-Product":           "golang-plex",
		"X-Plex-Version":           "0.0.1",
		"X-Plex-Client-Identifier": "golang-plex",
		"X-Plex-Platform":          "golang",
		"X-Plex-Platform-Version":  runtime.Version(),
		"X-Plex-Device":            "N/A",
		"X-Plex-Device-Name":       "golang-plex",
	}
)

func addPlexHeaders(r *http.Request) {
	for k, v := range plexHeaders {
		r.Header.Add(k, v)
	}
}
