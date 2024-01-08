# Codefresh SDK for Golang

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/codefresh-inc/codefresh-io%2Fgo-sdk%2Fgo-sdk?type=cf-1)]( https://g.codefresh.io/public/accounts/codefresh-inc/pipelines/codefresh-io/go-sdk/go-sdk) 
[![GoDoc](https://godoc.org/github.com/codefresh-io/go-sdk?status.svg)](https://godoc.org/github.com/codefresh-io/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/codefresh-io/go-sdk)](https://goreportcard.com/report/github.com/codefresh-io/go-sdk)

# Start

`go get -u github.com/codefresh-io/go-sdk`

```go
import (
    "fmt"
    "os"

    "github.com/codefresh-io/go-sdk/pkg/utils"
    "github.com/codefresh-io/go-sdk/pkg"
)

func main() {
    path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
    options, err := utils.ReadAuthContext(path, "")
    if err != nil {
        fmt.Println("Failed to read codefresh config file")
        panic(err)
    }
    clientOptions := codefresh.ClientOptions{Host: options.URL,
        Auth: codefresh.AuthOptions{Token: options.Token}}
    cf := codefresh.New(&clientOptions)
    pipelines, err := cf.Pipelines().List()
    if err != nil {
        fmt.Println("Failed to get Pipelines from Codefresh API")
        panic(err)
    }
    for _, p := range pipelines {
        fmt.Printf("Pipeline: %+v\n\n", p)
    }
}

```

This is not an official Codefresh project.
