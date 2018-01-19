package gluo

import (
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
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

func request(event events.APIGatewayProxyRequest) (*http.Request, error) {
	body, err := getBodyReader(event)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(event.HTTPMethod, mustBuildURL(event), body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{}
	for name, value := range event.Headers {
		req.Header.Add(name, value)
	}
	if req.ContentLength > -1 {
		req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	}
	req.Header.Set("X-Request-ID", event.RequestContext.RequestID)
	req.Header.Set("X-Stage", event.RequestContext.Stage)
	req.RemoteAddr = event.RequestContext.Identity.SourceIP
	return req, nil
}
