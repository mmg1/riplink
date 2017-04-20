package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mschwager/riplink/src/requests"
)

func main() {
	var queryUrl string
	flag.StringVar(&queryUrl, "url", "https://google.com", "URL to query")

	var timeout int
	flag.IntVar(&timeout, "timeout", 5, "Timeout in seconds")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")

	var depth uint
	flag.UintVar(&depth, "depth", 1, "Follow discovered links this deep")

	var sameDomain bool
	flag.BoolVar(&sameDomain, "same-domain", false, "Only query links of the same domain as the initial URL")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, queryUrl, depth, sameDomain, results)

	for result := range results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}

		if verbose || result.Code < 200 || result.Code > 299 {
			fmt.Println(result.Url, result.Code)
		}
	}
}
