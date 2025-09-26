package sheet

import (
	"github.com/murasame29/hackathon-util/internal/datasource"
	"google.golang.org/api/sheets/v4"
)

type DataSource struct {
	sheetID string
	range_  string
	ss      *sheets.Service
}

func NewDataSource(sheetID, range_ string) *DataSource {
	return &DataSource{sheetID: sheetID, range_: range_}
}

func (d *DataSource) Read() (*datasource.ReadDataSourceResult, error) {
	sheetValues, err := d.ss.Spreadsheets.Values.Get(d.sheetID, d.range_).Do()
	if err != nil {
		return nil, err
	}

	result := new(datasource.ReadDataSourceResult)
	result.Teams = make(map[string][]string)
	result.Members = make([]string, 0)
	result.TeamNames = make([]string, 0)

	indexes := make(map[int]string)
	var csvIndex int

	for _, records := range sheetValues.Values {

		if csvIndex == 0 {
			for i, item := range records {
				indexes[i] = item.(string)
			}
		}

		var (
			teamName string
			teams    []string
		)
		for i, item := range records {
			switch indexes[i] {
			case "team_name":
				if item == "" {
					continue
				}
				teamName = item.(string)
			case "メンバー":
				if item == "" {
					continue
				}
				teams = append(teams, item.(string))
			}
		}

		result.Members = append(result.Members, teams...)
		result.TeamNames = append(result.TeamNames, teamName)
		result.Teams[teamName] = teams

		csvIndex++
	}
	return result, nil
}
