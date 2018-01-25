package gluo

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func mustBuildURL(event events.APIGatewayProxyRequest) string {
	query := url.Values{}
	for name, value := range event.QueryStringParameters {
		query.Add(name, value)
	}
	u := url.URL{
		Scheme:   event.Headers["X-Forwarded-Proto"],
		Host:     event.Headers["Host"],
		Path:     event.Path,
		RawQuery: query.Encode(),
	}
	return u.String()
}

func getBodyReader(event events.APIGatewayProxyRequest) (io.Reader, error) {
	var body io.Reader = strings.NewReader(event.Body)
	if event.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}
	return body, nil
}

// An unexported type to be used as the key for types in this package.
type contextKey struct{}

// The key for a LambdaContext in Contexts.
var reqCtxKey = &contextKey{}

// APIGatewayContext returns the APIGatewayProxyRequestContext value stored in ctx.
func APIGatewayContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	c, ok := ctx.Value(reqCtxKey).(events.APIGatewayProxyRequestContext)
	return c, ok
}

// LambdaContext return the LambdaContext value stored in ctx.
func LambdaContext(ctx context.Context) (*lambdacontext.LambdaContext, bool) {
	return lambdacontext.FromContext(ctx)
}

func request(ctx context.Context, event events.APIGatewayProxyRequest) (*http.Request, error) {
	body, err := getBodyReader(event)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(event.HTTPMethod, mustBuildURL(event), body)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, reqCtxKey, event.RequestContext)
	req = req.WithContext(ctx)
	req.Header = http.Header{}
	for name, value := range event.Headers {
		req.Header.Add(name, value)
	}
	if req.ContentLength > -1 {
		req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	}
	req.Header.Set("X-Request-ID", event.RequestContext.RequestID)
	req.Header.Set("X-Stage", event.RequestContext.Stage)
	// AWS X-Ray
	if traceID := ctx.Value("x-amzn-trace-id"); traceID != nil {
		req.Header.Set("X-Amzn-Trace-Id", fmt.Sprintf("%v", traceID))
	}
	req.RemoteAddr = event.RequestContext.Identity.SourceIP
	return req, nil
}
