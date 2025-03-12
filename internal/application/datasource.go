package application

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"google.golang.org/api/sheets/v4"
)

type ReadDataSourceResult struct {
	TeamNames []string
	Members   []string
	Teams     map[string][]string
}

type DataSource interface {
	Read(string) (*ReadDataSourceResult, error)
}

type DataSourceCSV struct {
	path string
}

func NewDataSourceCSV(path string) *DataSourceCSV {
	return &DataSourceCSV{path: path}
}

func (d *DataSourceCSV) Read() (*ReadDataSourceResult, error) {
	file, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := new(ReadDataSourceResult)
	result.Teams = make(map[string][]string)
	result.Members = make([]string, 0)
	result.TeamNames = make([]string, 0)

	indexes := make(map[int]string)
	var csvIndex int

	teamTables := csv.NewReader(file)
	for {
		record, err := teamTables.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if csvIndex == 0 {
			for i, item := range record {
				indexes[i] = item
			}
		}

		var (
			teamName string
			teams    []string
		)
		for i, item := range record {
			switch indexes[i] {
			case "チーム名":
				if item == "" {
					continue
				}
				teamName = item
			case "メンバー":
				if item == "" {
					continue
				}
				teams = append(teams, item)
			}
		}

		result.Members = append(result.Members, teams...)
		result.TeamNames = append(result.TeamNames, teamName)
		result.Teams[teamName] = teams

		csvIndex++
	}
	return result, nil
}

type DataSourceGoogleSheets struct {
	sheetID string
	range_  string
	ss      *sheets.Service
}

func NewDataSourceGoogleSheets(sheetID, range_ string) *DataSourceGoogleSheets {
	return &DataSourceGoogleSheets{sheetID: sheetID, range_: range_}
}

func (d *DataSourceGoogleSheets) Read() (*ReadDataSourceResult, error) {
	sheetValues, err := d.ss.Spreadsheets.Values.Get(d.sheetID, d.range_).Do()
	if err != nil {
		return nil, err
	}

	result := new(ReadDataSourceResult)
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
			case "チーム名":
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
