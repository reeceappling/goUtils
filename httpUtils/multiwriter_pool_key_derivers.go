package httpUtils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RequestPoolKeyDeriver func(*http.Request) (string, error)

func GeneratePoolKeyDeriver(u URLPoolKeyDeriver, h HeaderPoolKeyDeriver, b BodyPoolKeyDeriver) RequestPoolKeyDeriver {
	return func(r *http.Request) (out string, err error) {
		out, uPart, hPart, bPart := "", "", "", ""
		if u != nil {
			uPart, err = u(r.URL)
			if err != nil {
				return
			}
		}
		if h != nil {
			hPart, err = h(r.Header)
			if err != nil {
				return
			}
		}
		if b != nil {
			bodyBytes, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			bPart, err = b(bytes.NewReader(bodyBytes))
			if err != nil {
				return
			}
		}

		return uPart + hPart + bPart, nil
	}
}

var (
	_ URLPoolKeyDeriver = defaultUrlKeyDeriver
	_ URLPoolKeyDeriver = defaultShortenedUrlKeyDeriver
)

type URLPoolKeyDeriver func(url *url.URL) (string, error)

func defaultShortenedUrlKeyDeriver(urlPtr *url.URL) (string, error) {
	if urlPtr == nil {
		return "", nil
	}
	URL := *urlPtr
	return fmt.Sprintf("%s%s%s%s", URL.Scheme, URL.Host, URL.Port(), URL.EscapedPath()), nil // TODO: params?
}
func defaultUrlKeyDeriver(urlPtr *url.URL) (string, error) {
	if urlPtr == nil {
		return "", nil
	}
	URL := *urlPtr
	return fmt.Sprintf("%s://%s:%s%s", URL.Scheme, URL.Host, URL.Port(), URL.EscapedPath()), nil
}

var (
	_ HeaderPoolKeyDeriver = nilHeaderKeyDeriver
	_ HeaderPoolKeyDeriver = defaultHeaderKeyDeriver
)

type HeaderPoolKeyDeriver func(http.Header) (string, error)

func nilHeaderKeyDeriver(http.Header) (string, error) {
	return "", nil
}
func defaultHeaderKeyDeriver(http.Header) (string, error) {
	return "", nil // TODO: change
}

var (
	_ BodyPoolKeyDeriver = defaultBodyPoolKeyDeriver
	_ BodyPoolKeyDeriver = defaultShortenedBodyPoolKeyDeriver
)

type BodyPoolKeyDeriver func(io.Reader) (string, error) // TODO: turn into reader as input?

func defaultBodyPoolKeyDeriver(body io.Reader) (string, error) {
	bs, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
func defaultShortenedBodyPoolKeyDeriver(body io.Reader) (string, error) { // TODO: use
	h := sha256.New()
	_, err := io.Copy(h, body)
	if err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}
