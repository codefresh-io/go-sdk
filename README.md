# Codefresh SDK for Golang

# Start

`go get -u github.com/codefresh-io`

```go

import (
    "fmt"
    "os"

    "github.com/codefresh-io/go-sdk/pkg/util"
    "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

func main() {
    path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
    authOptions := util.ReadAuthContext(path, "")
    cf := codefresh.New(authOptions)
    cf.GetPipelines()
}
```