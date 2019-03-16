# go-paycor
Go library for interacting with the Paycor's Reporting API.  This implements the minium needed to generate & download an exisitng report.

You'll have to talk to them about getting access for your company.

## Use
```go
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/MikeAlbertFleetSolutions/go-paycor"
)

func main() {
	var err error

	// connection to paycor
	paycorClient := paycor.NewClient("xxx", "xxx", "secure.paycor.com")

	// get naughtylist report from paycor
	var naughtylist []byte
	naughtylist, err = paycorClient.GetReportByName("hr naughty list")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// csv
	r := csv.NewReader(bytes.NewReader(naughtylist))

	// process rows
	for {
		// read row
		var record []string
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%+v", err)
		}

		fmt.Printf("%#v\n", record)
	}
}
```
