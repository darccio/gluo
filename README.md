# Gluo

Write your HTTP application in Go once, deploy it everywhere (on-premise & AWS Lambda supported)

# Installation

```
go get github.com/imdario/gluo

import (
	"github.com/imdario/gluo"
)
```

# Usage

Just call `gluo.ListenAndServe` instead of `http.ListenAndServe`.

```go
package main

import (
        "github.com/imdario/gluo"
        "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello from Gluo"))
}

func main() {
        gluo.ListenAndServe(":3000", http.HandlerFunc(handler))
}
```
