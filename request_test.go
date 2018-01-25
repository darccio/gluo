package gluo

import (
	"bytes"
	"encoding/base64"
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
