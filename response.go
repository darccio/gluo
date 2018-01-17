package gluo

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type response struct {
	event       *events.APIGatewayProxyResponse
	wroteHeader bool
	buffer      bytes.Buffer
	header      http.Header
}

func (w *response) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *response) Write(data []byte) (int, error) {
	lenData := len(data)
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if lenData == 0 {
		return 0, nil
	}
	return w.buffer.Write(data)
}

func (w *response) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.event.StatusCode = code
}

func (w *response) finishRequest() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.event.Body = w.buffer.String()
	w.buffer.Reset()
	w.event.Headers = make(map[string]string)
	for name, value := range w.header {
		w.event.Headers[name] = value[0]
	}
	w.header = make(http.Header)
}
