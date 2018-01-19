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
	rq, err := http.NewRequest(event.HTTPMethod, mustBuildURL(event), body)
	if err != nil {
		return nil, err
	}
	rq.Header = http.Header{}
	for name, value := range event.Headers {
		rq.Header.Add(name, value)
	}
	if rq.ContentLength > -1 {
		rq.Header.Set("Content-Length", strconv.FormatInt(rq.ContentLength, 10))
	}
	rq.Header.Set("X-Request-ID", event.RequestContext.RequestID)
	rq.Header.Set("X-Stage", event.RequestContext.Stage)
	rq.RemoteAddr = event.RequestContext.Identity.SourceIP
	return rq, nil
}
