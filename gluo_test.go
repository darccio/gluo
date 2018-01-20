package gluo

import (
	"flag"
	"testing"
	"os"
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
