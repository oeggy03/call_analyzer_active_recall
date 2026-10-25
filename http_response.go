package main

import (
	"fmt"
	"io"
)

const maxResponseBodyBytes = 10 << 20

func readResponseBody(body io.Reader, secrets ...string) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxResponseBodyBytes+1))
	if err != nil {
		return nil, redactedError(fmt.Sprintf("reading response body: %v", err), secrets...)
	}
	if len(data) > maxResponseBodyBytes {
		return nil, redactedError("response body is too large", secrets...)
	}
	return data, nil
}
