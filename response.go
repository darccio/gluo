package gluo

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"mime"
	"net/http"
	"strings"
)

type response struct {
	event       events.APIGatewayProxyResponse
	wroteHeader bool
	header      http.Header
	buffer      bytes.Buffer
}

func (w *response) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *response) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if len(data) == 0 {
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

func (w *response) finish() (events.APIGatewayProxyResponse, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.event.IsBase64Encoded = w.isBinary()
	if w.event.IsBase64Encoded {
		w.event.Body = base64.StdEncoding.EncodeToString(w.buffer.Bytes())
	} else {
		w.event.Body = w.buffer.String()
	}
	w.event.Headers = make(map[string]string)
	for name, value := range w.Header() {
		if len(value) > 0 {
			w.event.Headers[name] = value[0]
		}
	}
	return w.event, nil
}

func (w *response) getContentType() string {
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(w.buffer.Bytes())
	}
	return contentType
}

func (w *response) isBinary() bool {
	encoding := w.Header().Get("Content-Encoding")
	if len(encoding) > 0 && encoding != "identity" {
		return true
	}
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		// Content-Type shouldn't be empty but this is a last resort
		w.Header().Set("Content-Type", "application/octet-stream")
		return true
	}
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		// we respect the original Content-Type header but assume it's binary
		return true
	}
	// https://www.iana.org/assignments/media-types/media-types.xhtml
	if strings.HasPrefix(mimeType, "text/") {
		// we shortcut explicit text MIME types
		return false
	}
	switch {
	case strings.HasSuffix(mimeType, "xml") && !strings.HasSuffix(mimeType, "wbxml"):
		fallthrough
	case strings.HasSuffix(mimeType, "json-seq"):
		fallthrough
	case strings.HasSuffix(mimeType, "json"):
		return false
	}
	return true
}
