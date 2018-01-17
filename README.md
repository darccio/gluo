# Gluo

Write your HTTP application in Go once, deploy it everywhere (on-premise & AWS Lambda supported)

## Status

Alpha quality. It needs to be tested out extensively but it seems to work fine.

[![GoDoc][1]][2]

[1]: https://godoc.org/github.com/imdario/gluo?status.svg
[2]: https://godoc.org/github.com/imdario/gluo

## Installation

    go get github.com/imdario/gluo

    import (
        "github.com/imdario/gluo"
    )

## Usage

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

## Contact me

If I can help you, you have an idea or you are using Gluo in your projects, don't hesitate to drop me a line (or a pull request): [@im_dario](https://twitter.com/im_dario)

## About

Written by [Dario Castañé](https://twitter.com/im_dario).

## License

[BSD 3-Clause](http://opensource.org/licenses/BSD-3-Clause) license.
