package gluo

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

// IsLambda checks if it is running on AWS Lambda.
func IsLambda() bool {
	return os.Getenv("_LAMBDA_SERVER_PORT") != ""
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
	w.finish()
	return rs, nil
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
		adapter := lambdaAdapter{handler}
		lambda.Start(adapter.Serve)
		return nil
	}
	return http.ListenAndServe(addr, handler)
}
