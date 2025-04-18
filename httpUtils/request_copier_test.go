package httpUtils

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestRequestCopier(t *testing.T) {
	expBody := "expBody"
	expBodyReadCloser := io.NopCloser(strings.NewReader(expBody))
	defer expBodyReadCloser.Close()
	exp := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "aScheme",
		},
		Header:  http.Header{"a": []string{"b"}},
		Body:    expBodyReadCloser,
		GetBody: func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(expBody)), nil },
	}
	copier, client := NewCopyRoundTripperOnDefaultClient()
	resp, _ := client.Do(exp)
	if resp != nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, exp.URL, copier.URL)
	assert.True(t, reflect.DeepEqual(exp.Header, copier.Headers))
	assert.Equal(t, expBody, string(copier.BodySent))
}
