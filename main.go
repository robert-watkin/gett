package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	url := flag.Arg(0)
	if url == "" {
		log.Fatalf("URL not provided")
	}

	fmt.Fprintf(os.Stderr, "The provided URL is: %s\n", url)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to GET: %v", err)
	}
	defer response.Body.Close()

	fmt.Println("Response status: ", response.Status)

	io.Copy(os.Stdout, response.Body)
}
