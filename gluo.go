package gluo

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

// IsLambda checks if it is running on AWS Lambda.
func IsLambda() bool {
	_, ok := os.LookupEnv("_LAMBDA_SERVER_PORT")
	return ok
}

// LambdaAdapter allows to wrap a http.Handler for working on AWS Lambda.
type LambdaAdapter struct {
	http.Handler
}

// Handle handles an AWS Lambda request, converting an APIGatewayProxyRequest into a http.Request.
func (la LambdaAdapter) Handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := request(ctx, event)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	w := &response{}
	la.Handler.ServeHTTP(w, req)
	return w.finish()
}

// ListenAndServe calls http.ListenAndServe with handler to handle
// requests on incoming connections when invoked outside AWS Lambda
// environment.
//
// If handler is nil, it uses http.DefaultServeMux as default handler.
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
	if handler == nil {
		handler = http.DefaultServeMux
	}
	if IsLambda() {
		adapter := LambdaAdapter{handler}
		lambda.Start(adapter.Handle)
		return nil
	}
	return http.ListenAndServe(addr, handler)
}
