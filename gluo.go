package gluo

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

func isLambda() bool {
	functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	return functionName != ""
}

type lambdaAdapter struct {
	http.Handler
}

func (la lambdaAdapter) Serve(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rq, err := request(event)
	rs := events.APIGatewayProxyResponse{}
	if err != nil {
		return rs, err
	}
	w := response{event: &rs}
	la.Handler.ServeHTTP(&w, rq)
	w.finishRequest()
	return rs, nil
}

// ListenAndServe calls http.ListenAndServe with handler to handle
// requests on incoming connections when invoked outside AWS Lambda
// environment.
//
// In AWS Lambda environment it calls lambda.Start with handler to
// handle requests transparently as if it was working in a non
// serverless environment. Addr argument is ignored in AWS Lambda.
//
// A trivial example is:
//
//	package main
//
//	import (
//		"github.com/imdario/gluo"
//		"net/http"
//	)
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Hello from Gluo"))
//	}
//
//	func main() {
//		gluo.ListenAndServe(":3000", http.HandlerFunc(handler))
//	}
//
// ListenAndServe always returns a non-nil error. Under AWS Lambda,
// it always returns nil.
func ListenAndServe(addr string, handler http.Handler) error {
	// TODO Handle nil handler as DefaultServeMux
	if isLambda() {
		adapter := lambdaAdapter{handler}
		lambda.Start(adapter.Serve)
		return nil
	} else {
		return http.ListenAndServe(addr, handler)
	}
}
