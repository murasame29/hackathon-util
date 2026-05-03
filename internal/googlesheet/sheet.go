package googlesheet

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GetTeamData fetches team data from a Google Spreadsheet.
func GetTeamData(spreadsheetID, rangeStr, credentialsFile string) ([][]any, error) {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("googlesheet: create service: %w", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeStr).Do()
	if err != nil {
		return nil, fmt.Errorf("googlesheet: get values %s/%s: %w", spreadsheetID, rangeStr, err)
	}

	return resp.Values, nil
}
