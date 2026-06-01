package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type headerFlags []string

func (h headerFlags) String() string {
	return strings.Join(h, ",")
}

func (h *headerFlags) Set(value string) error {
	if !strings.Contains(value, ":") {
		return fmt.Errorf("header %q must be in \"Name: value\" form", value)
	}
	*h = append(*h, value)
	return nil
}

type options struct {
	url     string
	headers headerFlags
}

func main() {
	// struct to hold the options returned
	var opts options

	flag.Var(&opts.headers, "header", "Multiple optional headers")

	flag.Parse()

	opts.url = flag.Arg(0)
	if opts.url == "" {
		log.Fatalf("URL not provided")
	}

	fmt.Fprintf(os.Stderr, "The provided URL is: %s\n", opts.url)

	req, err := http.NewRequest("get", opts.url, nil)
	if err != nil {
		log.Fatalf("Fatal error creating request for %v", opts.url)
	}

	for _, strHeader := range opts.headers {
		splitHeader := strings.Split(strHeader, ":")
		key := splitHeader[0]
		value := splitHeader[1]

		if key == "" || value == "" {
			log.Fatalf("Header %v is invalid", strHeader)
		}

		fmt.Fprintf(os.Stderr, "Adding key: %s; value: %s; to headers...\n", key, value)

		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)

	io.Copy(os.Stdout, resp.Body)
}
