package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/murasame29/hackathon-util/internal/datasource"
)

type DataSource struct {
	path string
}

func NewDataSource(path string) *DataSource {
	return &DataSource{path: path}
}

func (d *DataSource) Read() (*datasource.ReadDataSourceResult, error) {
	file, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := new(datasource.ReadDataSourceResult)
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
			case "team_name":
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
