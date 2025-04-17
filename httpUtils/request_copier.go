package httpUtils

import (
	"bytes"
	"io"
	"net/http"
)

// RequestCopier keeps info for the last request that passed through it, and is an http.RoundTripper,
// meant to be the Transport field on a http.Client.Transport
type RequestCopier struct {
	Headers       http.Header
	BodySent      []byte
	wrappedClient *http.Client
}

// RoundTrip meets the interface of http.RoundTripper. Copies the request before sending it to the wrapped client
func (rc *RequestCopier) RoundTrip(req *http.Request) (res *http.Response, err error) {
	rc.Headers = req.Header.Clone()
	rc.BodySent, err = io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(rc.BodySent))
	return rc.wrappedClient.Do(req)
}

// NewCopyRoundTripperOnDefaultClient returns a RequestCopier, as well as the provided http.Client with the copier on it
func NewCopyRoundTripperWithClient(clientIn *http.Client) (rtHoldingRequestInformation *RequestCopier, clientThatCopiesOutput *http.Client) {
	clientOut := *clientIn
	crt := &RequestCopier{
		Headers:       http.Header{},
		BodySent:      []byte{},
		wrappedClient: clientIn,
	}
	clientOut.Transport = crt
	return crt, &clientOut
}

// NewCopyRoundTripperOnDefaultClient returns a RequestCopier that is on the default http.Client
func NewCopyRoundTripperOnDefaultClient() (*RequestCopier, *http.Client) {
	return NewCopyRoundTripperWithClient(http.DefaultClient)
}
