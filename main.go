package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type options struct {
	url     string
	headers headerFlags
	timeout time.Duration
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

func main() {
	// struct to hold the options returned
	opts, err := parseArgs()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	if err := run(opts); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
