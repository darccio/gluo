package gluo

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"net/url"
	"strings"
)

func mustBuildURL(event events.APIGatewayProxyRequest) string {
	u := url.URL{
		Scheme: event.Headers["X-Forwarded-Proto"],
		Host:   event.Headers["Host"],
		Path:   event.Path,
	}
	query := url.Values{}
	for name, value := range event.QueryStringParameters {
		query.Add(name, value)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func request(event events.APIGatewayProxyRequest) (*http.Request, error) {
	body := strings.NewReader(event.Body)
	rq, err := http.NewRequest(event.HTTPMethod, mustBuildURL(event), body)
	if err != nil {
		return nil, err
	}
	rq.Header = http.Header{}
	for name, value := range event.Headers {
		rq.Header.Add(name, value)
	}
	return rq, nil
}
