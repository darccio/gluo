package gluo

import (
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"net/http"
	"os"
)

// EntryPoint is the the exported handler called by AWS Lambda.
type EntryPoint apigatewayproxy.Handler

func isLambda() bool {
	if len(os.Args) == 0 {
		panic("unable to tell if we are running in a lambda or not")
	}
	if os.Args[0] == "/usr/bin/python2.7" {
		return true
	}
	return false
}

func Init(entryPoint *EntryPoint, handle http.HandlerFunc) {
	if isLambda() {
		listener := net.Listen()
		*entryPoint = apigatewayproxy.New(listener, nil).Handle
		go http.Serve(listener, handle)
	}
}

func ListenAndServe(addr string, handle http.HandlerFunc) {
	http.ListenAndServe(addr, handle)
}
