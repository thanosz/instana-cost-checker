package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/dustin/go-humanize"
)

type Item struct {
	Name string `json:"name"`
	Sims uint64 `json:"sims"`
}

type Data struct {
	Time  int64  `json:"time"`
	Items []Item `json:"items"`
}

func main() {

	var month int
	var year int
	var token string
	var endpoint string
	var maxallowed string
	var threshold float64

	current_year, current_month, _ := time.Now().Date()

	flag.Usage = func() {
		fmt.Printf("\nCheck the amount of data ingested by Instana server and produce a warning if exceeds the specified threshold; returns 1 if threshold exceeded, 0 otherwise.\n")

		fmt.Printf("Example: check current month usage and produce a warning if the total ingested data is at 70%% of the allowed ingested data.\n\n")
		fmt.Printf("	instana-cost-checker -token TOKEN -endpoint unit-tenant.instana.io -maxallowed 7TB -threshold 0.7\n\n")

		fmt.Printf("Options:\n")
		flag.PrintDefaults()

	}

	flag.IntVar(&month, "month", int(current_month), "The month of the year to request data for (optional, skip for current month)")
	flag.IntVar(&year, "year", int(current_year), "The year (optional, skip for current year)")
	flag.StringVar(&token, "token", "", "The authentication token to use (required)")
	flag.StringVar(&endpoint, "endpoint", "", "The endpoint to connect to (e.g. unit-tenant.instana.io, required)")
	flag.StringVar(&maxallowed, "maxallowed", "", "The maximum entitled data usage in MB, GB or TB (e.g. 7TB, required)")
	flag.Float64Var(&threshold, "threshold", 0.8, "The percentage to multiply with to generate a warning (optional)")

	flag.Parse()
	maxallowed_d, _ := datasize.ParseString(maxallowed)
	//fmt.Printf("Threshold %d \n", thres.Bytes())

	if month < int(current_month) || year < int(current_year) || len(token) == 0 || len(endpoint) == 0 || maxallowed_d == 0 {
		flag.Usage()
		os.Exit(0)
	}

	url := "https://" + endpoint + "/api/instana/usage/api/" + strconv.Itoa(month) + "/" + strconv.Itoa(year)
	apiToken := "apiToken " + token // Replace with your actual API token

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the authorization header
	req.Header.Add("Authorization", apiToken)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse the JSON response
	var data []Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Calculate the sum
	var totalBytesIngestedInfrastructure uint64
	var totalBytesIngestedTraces uint64

	for _, entry := range data {
		for _, item := range entry.Items {
			switch item.Name {
			case "bytes_ingested_infrastructure_acceptor":
				totalBytesIngestedInfrastructure += item.Sims
			case "bytes_ingested_traces_otlp_acceptor":
				totalBytesIngestedTraces += item.Sims
			}
		}
	}

	// Print the results
	fmt.Printf("Total bytes_ingested_infrastructure_acceptor: %s\n", humanize.Bytes(totalBytesIngestedInfrastructure))
	fmt.Printf("Total bytes_ingested_traces_otlp_acceptor: %s\n", humanize.Bytes(totalBytesIngestedTraces))
	total := totalBytesIngestedInfrastructure + totalBytesIngestedTraces
	fmt.Printf("Total Usage for month (%d): %s\n", month, humanize.Bytes(total))

	if total >= uint64(float64(maxallowed_d.Bytes())*threshold) {
		fmt.Printf("Threshold warning \n")
		os.Exit(1)

	}
	os.Exit(0)

}
