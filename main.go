package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
	timeout time.Duration
}

func main() {
	// struct to hold the options returned
	opts, err := parseArgs()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	req, err := http.NewRequest("GET", opts.url, nil)
	if err != nil {
		log.Fatalf("Fatal error creating request for %v", opts.url)
	}

	for _, strHeader := range opts.headers {
		splitHeader := strings.SplitN(strHeader, ":", 2)
		key := splitHeader[0]
		value := splitHeader[1]

		if key == "" || value == "" {
			log.Fatalf("Header %v is invalid", strHeader)
		}

		fmt.Fprintf(os.Stderr, "Adding key: %s; value: %s; to headers...\n", key, value)

		req.Header.Add(key, value)
	}

	client := &http.Client{Timeout: opts.timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)

	io.Copy(os.Stdout, resp.Body)
}

func parseArgs() (options, error) {
	var opts options

	fs := flag.NewFlagSet("gett", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Var(&opts.headers, "header", "Multiple optional headers")
	fs.DurationVar(&opts.timeout, "timeout", 30*time.Second, "request timeout")
	flag.Duration("timeout", opts.timeout, "Timeout for the request")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return opts, fmt.Errorf("Failed to parse options: %v", err)
	}

	if fs.NArg() != 1 {
		return opts, errors.New("usage: gett [flags] URL")
	}
	opts.url = fs.Arg(0)

	return opts, nil
}
