# go-crtsh

A Go client library for the [crt.sh](https://crt.sh) certificate transparency database API.

## Installation

```bash
go get github.com/The-Infra-Company/go-crtsh
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    gocrtsh "github.com/The-Infra-Company/go-crtsh"
)

func main() {
    client := gocrtsh.New()
    ctx := context.Background()

    records, err := client.BasicSearch(ctx, "example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d certificates for example.com\n", len(records))
}
```

## Usage Examples

### Basic Search

Search for certificates issued specifically for a domain:

```go
records, err := client.BasicSearch(ctx, "example.com")
```

### Wildcard Search

Search for certificates issued for a domain and all its subdomains:

```go
records, err := client.SearchWithWildcard(ctx, "example.com", false)
```

### Include Expired Certificates

Search for both active and expired certificates:

```go
opts := &gocrtsh.SearchOptions{
    IncludeExpired: true,
    Wildcard:       false,
}

records, err := client.Search(ctx, "example.com", opts)
```

