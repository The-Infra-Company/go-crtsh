package main

import (
	"context"
	"fmt"
	"log"

	gocrtsh "github.com/The-Infra-Company/go-crtsh"
)

// Example:
// ❯ ./basic_search
// Found 33 certificates for example.com

func main() {
	client := gocrtsh.New()
	ctx := context.Background()

	records, err := client.BasicSearch(ctx, "example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d certificates for example.com\n", len(records))
}
