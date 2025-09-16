package main

import (
	"context"
	"fmt"
	"log"

	gocrtsh "github.com/The-Infra-Company/go-crtsh"
)

// Example:
// ❯ ./wildcard_search
// Found 127 certificates for *.example.com

func main() {
	client := gocrtsh.New()
	ctx := context.Background()

	records, err := client.SearchWithWildcard(ctx, "example.com", false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d certificates for *.example.com\n", len(records))
}
