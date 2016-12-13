package gsheets

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

// ClientFactory is the function that produces an HTTP client to be used with the Google API. Useful for testing and mocking
type ClientFactory func(configFile string) *http.Client

// DefaultClientFactory is the function that includes the required authorization in each request to the Google API
var DefaultClientFactory = func(configFile string) *http.Client {
	ctx := context.Background()

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)

	}

	return getClient(ctx, config)

}

// NewService performs an authentication against the Google Sheets API and returns a Sheet Service
func NewService(configFile string, cf ClientFactory) (*sheets.Service, error) {

	client := cf(configFile)

	return sheets.New(client)
}

// Test is a test
// https://godoc.org/google.golang.org/api/sheets/v4
func Test(srv *sheets.Service) {

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetID := "1_gLn3VkhOc129zwbRdQMpMjZpAGzFpgmnW0_5y8hVgY"
	readRange := "Class Data!A2:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) > 0 {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%s, %s\n", row[0], row[4])
		}
	} else {
		fmt.Print("No data found.")
	}

}

// TestWrite attempts to write values to a spreadsheet
// Documentation here: https://developers.google.com/sheets/api/guides/value
// and here https://developers.google.com/sheets/api/samples/writing
// and here https://godoc.org/google.golang.org/api/sheets/v4
func TestWrite(srv *sheets.Service) error {
	spreadsheetID := "1_gLn3VkhOc129zwbRdQMpMjZpAGzFpgmnW0_5y8hVgY"
	writeRange := "Class Data!A1:B"

	data := []*sheets.ValueRange{
		&sheets.ValueRange{
			Range:          writeRange,
			MajorDimension: "ROWS",
			Values: [][]interface{}{
				[]interface{}{"Name", "Position"},
				[]interface{}{"Edo", "Dev"},
			},
		},
	}

	batchReq := sheets.BatchUpdateValuesRequest{
		Data: data,
		IncludeValuesInResponse: false,
		ValueInputOption:        "USER_ENTERED",
	}

	resp, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetID, &batchReq).Do()

	if err != nil {
		return errors.Wrap(err, "Failed to batch updated the spreadsheet")
	}

	log.Infof("Batch update succeded with %d cells", resp.TotalUpdatedCells)

	return nil
}
