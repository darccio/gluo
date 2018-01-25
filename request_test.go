package gluo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestBase64Body(t *testing.T) {
	var (
		expected = "Hello world!"
		body     = base64.StdEncoding.EncodeToString([]byte(expected))
		buffer   bytes.Buffer
	)
	e := events.APIGatewayProxyRequest{
		IsBase64Encoded: true,
		Body:            body,
	}
	r, err := getBodyReader(e)
	if err != nil {
		t.Errorf("unexpected error on getBodyReader: %v", err)
	}
	_, err = buffer.ReadFrom(r)
	if err != nil {
		t.Errorf("unexpected error reading body in base 64 from getBodyReader: %v", err)
	}
	result := buffer.String()
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestWrongBase64Body(t *testing.T) {
	var (
		body   = "Hello world!"
		buffer bytes.Buffer
	)
	e := events.APIGatewayProxyRequest{
		IsBase64Encoded: true,
		Body:            body,
	}
	r, err := getBodyReader(e)
	if err != nil {
		t.Errorf("unexpected error on getBodyReader: %v", err)
	}
	_, err = buffer.ReadFrom(r)
	if err == nil {
		t.Errorf("expected error reading body in base 64 from getBodyReader, got nil")
	}
}

func TestXRayHeader(t *testing.T) {
	event := events.APIGatewayProxyRequest{}
	json.Unmarshal([]byte(testRequest), &event)
	ctx := context.WithValue(context.Background(), "X-Amzn-Trace-Id", "Root=1-5759e988-bd862e3fe1be46a994272793;Sampled=1")
	req, err := request(ctx, event)
	if err != nil {
		t.Errorf("unexpected error on request: %v", err)
	}
	if req.Header.Get("X-Amzn-Trace-Id") == "" {
		t.Errorf("expected X-Amzn-Trace-Id injected but it wasn't found")
	}
}
