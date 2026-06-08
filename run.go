package main

import (
	"fmt"
	"io"
	"log"
	"mime"
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

	contentType := resp.Header.Get("Content-Type")

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 || mediaType != "application/json" {
		return fmt.Errorf("Response body was not JSON")
	}

	io.Copy(os.Stdout, resp.Body)
	return nil
}
