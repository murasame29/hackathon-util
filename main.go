package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetID = "1jtEqNR2eFKQTpgnuGmPJR3_bxScJXC-stePGSaiPEHE"

func main() {
	credential := option.WithCredentialsFile("./credential.json")
	srv, err := sheets.NewService(context.TODO(), credential)
	if err != nil {
		log.Fatal(err)
	}
	readRange := "チームごと!A2:F19"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()

	if err != nil {
		log.Fatalln(err)
	}
	if len(resp.Values) == 0 {
		log.Fatalln("data not found")
	}

	for _, row := range resp.Values {
		for _, col := range row {
			fmt.Printf("%s, ", col)
		}
		fmt.Println()
	}
}
