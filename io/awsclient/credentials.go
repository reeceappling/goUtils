package awsclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	v2signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

func GenerateSignedReqV2(ctx context.Context, url, requestText, acceptType, apiKey string) (reqToFire *http.Request, err error) {

	reqToFire, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(requestText))
	if err != nil {
		return nil, err
	}

	reqToFire.Header.Set("Accept", acceptType)
	reqToFire.Header["x-api-key"] = []string{apiKey} // Set will mangle the key to title case

	signer := v2signer.NewSigner()
	hash := sha256.New()
	hash.Write([]byte(requestText))
	if credentials.Expired() {
		credentials, err = refreshCredentials(ctx)
		if err != nil {
			return nil, err
		}
	}
	if err := signer.SignHTTP(ctx,
		credentials,
		reqToFire,
		hex.EncodeToString(hash.Sum(nil)),
		"execute-api",
		AwsRegion,
		time.Now(),
	); err != nil {
		return nil, err
	}

	return reqToFire, nil
}
