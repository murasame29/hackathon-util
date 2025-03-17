package gs

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSpreadSeet struct {
	ss *sheets.Service
}

func New(credential option.ClientOption) *GoogleSpreadSeet {
	ss, err := sheets.NewService(context.Background(), credential)
	if err != nil {
		panic(err)
	}

	return &GoogleSpreadSeet{ss}
}

func (gs *GoogleSpreadSeet) Read(spreadsheetID string, range_ string) ([][]string, error) {
	sheetValues, err := gs.ss.Spreadsheets.Values.Get(spreadsheetID, range_).Do()
	if err != nil {
		return nil, err
	}

	if len(sheetValues.Values) == 0 {
		return nil, fmt.Errorf("'%s' is empty", range_)
	}

	var values [][]string
	for _, row := range sheetValues.Values {
		var rowValues []string
		for _, col := range row {
			rowValues = append(rowValues, col.(string))
		}
		values = append(values, rowValues)
	}

	return values, nil
}
