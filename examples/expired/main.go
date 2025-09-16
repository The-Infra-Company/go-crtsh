package main

import (
	"context"
	"fmt"
	"log"

	gocrtsh "github.com/The-Infra-Company/go-crtsh"
)

// Example:
// ❯ ./expired_search
// Found 156 certificates for example.com (including expired)

func main() {
	client := gocrtsh.New()
	ctx := context.Background()

	opts := &gocrtsh.SearchOptions{
		IncludeExpired: true,
		Wildcard:       false,
	}

	records, err := client.Search(ctx, "example.com", opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d certificates for example.com (including expired)\n", len(records))
}
