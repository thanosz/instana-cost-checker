package main

import (
	"bytes"
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
	var verbose bool
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
	flag.BoolVar(&verbose, "verbose", false, "Verbose output for each day")
	flag.Parse()

	maxallowedBytes, _ := datasize.ParseString(maxallowed)

	if month < 1 || month > 12 || year > int(current_year) || len(token) == 0 || len(endpoint) == 0 || maxallowedBytes.Bytes() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	url := "https://" + endpoint + "/api/instana/usage/api/" + strconv.Itoa(month) + "/" + strconv.Itoa(year)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the authorization header
	req.Header.Add("Authorization", "apiToken "+token)

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
		fmt.Println("Response received:", string(body))
		return
	}

	// Calculate the sum
	var totalBytesIngestedInfrastructure uint64
	var totalBytesIngestedOtlpTraces uint64
	var totalBytesIngestedTraces uint64
	var totalBytesMobile uint64
	var totalBytesSpans uint64
	var totalBytesWebsite uint64
	var dayLog bytes.Buffer

	for _, entry := range data {
		sec := entry.Time / 1000
		msec := entry.Time % 1000
		timestamp := time.Unix(sec, msec*int64(time.Millisecond))
		fmt.Fprintf(&dayLog, "%s\n", timestamp)
		for _, item := range entry.Items {
			//fmt.Println(time.Unix(sec, msec*int64(time.Millisecond)))
			switch item.Name {
			case "bytes_ingested_infrastructure_acceptor":
				fmt.Fprintf(&dayLog, "	 (infra), %s\n", humanize.Bytes(item.Sims))
				totalBytesIngestedInfrastructure += item.Sims
			case "bytes_ingested_traces_otlp_acceptor":
				fmt.Fprintf(&dayLog, "	 (trace), %s\n", humanize.Bytes(item.Sims))
				totalBytesIngestedOtlpTraces += item.Sims
			case "bytes_ingested_traces_acceptor":
				fmt.Fprintf(&dayLog, "	 (trace), %s\n", humanize.Bytes(item.Sims))
				totalBytesIngestedTraces += item.Sims

			case "bytes_ingested_eum_mobile_eum_acceptor":
				fmt.Fprintf(&dayLog, "  (mobile), %s\n", humanize.Bytes(item.Sims))
				totalBytesMobile += item.Sims
			case "bytes_ingested_eum_spans_eum_acceptor":
				fmt.Fprintf(&dayLog, "   (spans), %s\n", humanize.Bytes(item.Sims))
				totalBytesSpans += item.Sims
			case "bytes_ingested_eum_website_eum_acceptor":
				fmt.Fprintf(&dayLog, " (website), %s\n", humanize.Bytes(item.Sims))
				totalBytesWebsite += item.Sims
			}
		}
	}
	// Print the results
	if verbose {
		fmt.Println(&dayLog)
	}
	fmt.Printf("Totals:\n")
	fmt.Printf("         infra: %s\n", humanize.Bytes(totalBytesIngestedInfrastructure))
	fmt.Printf("   otlp traces: %s\n", humanize.Bytes(totalBytesIngestedOtlpTraces))
	fmt.Printf("  agent traces: %s\n", humanize.Bytes(totalBytesIngestedTraces))
	fmt.Printf("        mobile: %s\n", humanize.Bytes(totalBytesMobile))
	fmt.Printf("         spans: %s\n", humanize.Bytes(totalBytesSpans))
	fmt.Printf("       website: %s\n", humanize.Bytes(totalBytesWebsite))

	usage := totalBytesIngestedInfrastructure + totalBytesIngestedOtlpTraces + totalBytesIngestedTraces + totalBytesMobile + totalBytesSpans + totalBytesWebsite
	fmt.Printf("\nTotal Usage for month %s: %s (%s) (%d%%)\n", time.Month(month), humanize.Bytes(usage), humanize.Comma(int64(usage)), int(float64(usage)/float64(maxallowedBytes.Bytes())*100))

	if usage >= uint64(float64(maxallowedBytes.Bytes())*threshold) {
		fmt.Printf("\nThreshold warning!\n")
		os.Exit(1)
	}
	os.Exit(0)
}
