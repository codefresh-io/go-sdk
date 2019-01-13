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
    "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

func main() {
    path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
    options := utils.ReadAuthContext(path, "")
    cf := codefresh.New(options)
    cf.GetPipelines()
}
```

This is not an official Codefresh project.