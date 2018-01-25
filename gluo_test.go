package gluo

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Setenv("_LAMBDA_SERVER_PORT", "3000")
	os.Exit(m.Run())
}

func TestIsLambda(t *testing.T) {
	if !IsLambda() {
		t.Errorf("TestMain failed to set environment to simulate AWS Lambda")
	}
}

const testRequest = `{
    "body":"{\"name\":\"Gluo\"}",
    "resource":"/{proxy+}",
    "requestContext":{
        "resourceId":"123456",
        "apiId":"1234567890",
        "resourcePath":"/{proxy+}",
        "httpMethod":"POST",
        "requestId":"c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
        "accountId":"123456789012",
        "identity":{
            "apiKey":null,
            "userArn":null,
            "cognitoAuthenticationType":null,
            "caller":null,
            "userAgent":"Custom User Agent String",
            "user":null,
            "cognitoIdentityPoolId":null,
            "cognitoIdentityId":null,
            "cognitoAuthenticationProvider":null,
            "sourceIp":"127.0.0.1",
            "accountId":null
        },
        "stage":"prod",
        "X-Amzn-Trace-Id":"Root=1-5759e988-bd862e3fe1be46a994272793;Sampled=1"
    },
    "queryStringParameters":{
        "foo":"bar"
    },
    "headers":{
        "Via":"1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)",
        "Accept-Language":"en-US,en;q=0.8",
        "CloudFront-Is-Desktop-Viewer":"true",
        "CloudFront-Is-SmartTV-Viewer":"false",
        "CloudFront-Is-Mobile-Viewer":"false",
        "X-Forwarded-For":"127.0.0.1, 127.0.0.2",
        "CloudFront-Viewer-Country":"US",
        "Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Upgrade-Insecure-Requests":"1",
        "X-Forwarded-Port":"443",
        "Host":"1234567890.execute-api.us-east-1.amazonaws.com",
        "X-Forwarded-Proto":"https",
        "X-Amz-Cf-Id":"cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA==",
        "CloudFront-Is-Tablet-Viewer":"false",
        "Cache-Control":"max-age=0",
        "User-Agent":"Custom User Agent String",
        "CloudFront-Forwarded-Proto":"https",
        "Accept-Encoding":"gzip, deflate, sdch"
    },
    "pathParameters":{
        "proxy":"path/to/resource"
    },
    "httpMethod":"POST",
    "stageVariables":{
        "baz":"qux"
    },
    "path":"/path/to/resource"
}`

type helloRequest struct {
	Name string
}

func TestLambdaServe(t *testing.T) {
	la := LambdaAdapter{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r.Body)
			r.Body.Close()
			hr := helloRequest{}
			err := json.Unmarshal(buffer.Bytes(), &hr)
			if err != nil {
				result := fmt.Sprintf("Sorry, I didn't understand you.")
				w.Write([]byte(result))
				return
			}
			if hr.Name == "" {
				hr.Name = "stranger"
			}
			result := fmt.Sprintf("Hello, %s.", hr.Name)
			w.Write([]byte(result))
		}),
	}
	req := events.APIGatewayProxyRequest{}
	if err := json.Unmarshal([]byte(testRequest), &req); err != nil {
		t.Errorf("unexpected error on json.Unmarshal: %v", err)
	}
	res, err := la.Handle(context.TODO(), req)
	if err != nil {
		t.Error("LambdaHandler.Handle must return nil")
	}
	if res.StatusCode != 200 {
		t.Errorf("expected HTTP status code 200, not '%d'", res.StatusCode)
	}
	if res.Body != "Hello, Gluo." {
		t.Errorf("expected Body 'Hello, Gluo.', not '%s'", res.Body)
	}
	if res.IsBase64Encoded {
		t.Error("unexpected base 64 encoding")
	}
	contentType, _ := res.Headers["Content-Type"]
	if contentType != "application/json" {
		t.Errorf("expected Content-Type header 'application/json', not '%s'", contentType)
	}
}
