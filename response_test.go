package gluo

import (
	"testing"
)

func TestIsBinary(t *testing.T) {
	tests := []struct {
		contentType     string
		contentEncoding string
		expected        bool
	}{
		{
			"application/json",
			"",
			false,
		},
		{
			"application/json",
			"deflate",
			true,
		},
		{
			"application/json",
			"identity",
			false,
		},
		{
			"application/atom+xml",
			"",
			false,
		},
		{
			"image/png",
			"",
			true,
		},
		{
			"image/png",
			"identity",
			true,
		},
		{
			"image/png",
			"deflate",
			true,
		},
		{
			"text/plain",
			"identity",
			false,
		},
		{
			"text/plain",
			"gzip",
			true,
		},
		{
			"text/plain",
			"",
			false,
		},
		{
			"application/vnd.wap-wbxml",
			"",
			true,
		},
		{
			"application/geo+json-seq",
			"",
			false,
		},
		{
			"application/geo+json-seq",
			"gzip",
			true,
		},
	}
	for _, test := range tests {
		w := response{}
		w.Header().Set("Content-Type", test.contentType)
		w.Header().Set("Content-Encoding", test.contentEncoding)
		if w.isBinary() != test.expected {
			if test.expected {
				t.Errorf("'%s' (encoding: '%s') should be binary", test.contentType, test.contentEncoding)
			} else {
				t.Errorf("'%s' (encoding: '%s') should be text", test.contentType, test.contentEncoding)
			}
		}
	}
}
