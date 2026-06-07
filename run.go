package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func run(opts options) error {
	req, err := http.NewRequest("GET", opts.url, nil)
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)

	io.Copy(os.Stdout, resp.Body)
	return nil
}
