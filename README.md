# Gluo

Write your HTTP application in Go once, deploy it\* to your servers and AWS Lambda. Really, it's a drop-in replacement for net/http.

*\*: it must be compiled as* `GOOS=linux` to be deployed in AWS Lambda.

## Status

Alpha quality. It needs to be tested out extensively but it seems to work fine.

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![Coverage Status][5]][6]

[1]: https://godoc.org/github.com/imdario/gluo?status.svg
[2]: https://godoc.org/github.com/imdario/gluo
[3]: https://goreportcard.com/badge/github.com/imdario/gluo
[4]: https://goreportcard.com/report/github.com/imdario/gluo
[5]: https://coveralls.io/repos/github/imdario/gluo/badge.svg?branch=master
[6]: https://coveralls.io/github/imdario/gluo?branch=master

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

## Alternatives

Here are some similar projects, although the

## Contact me

If I can help you, you have an idea or you are using Gluo in your projects, don't hesitate to drop me a line (or a pull request): [@im_dario](https://twitter.com/im_dario)

## About

Written by [Dario Castañé](https://twitter.com/im_dario).

## License

[BSD 3-Clause](http://opensource.org/licenses/BSD-3-Clause) license.

## Support on Beerpay
Hey dude! Help me out for a couple of :beers:!

[![Beerpay](https://beerpay.io/imdario/gluo/badge.svg?style=beer-square)](https://beerpay.io/imdario/gluo)  [![Beerpay](https://beerpay.io/imdario/gluo/make-wish.svg?style=flat-square)](https://beerpay.io/imdario/gluo?focus=wish)
