package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8000"

func ReadBody(rc io.ReadCloser) []byte {
	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	rc.Close()
	return b
}

func NewBody(b interface{}) io.Reader {
	bodyBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bodyBytes)
}

func NewJsonRequest(method, url, key string, body any) *http.Request {
	req, _ := http.NewRequest(method, BuildURL(url), NewBody(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", key)
	return req
}

func BuildURL(path string) string {
	return baseURL + path
}

func StringPtr(s string) *string { return &s }

func IntPtr[N int | int8 | int16 | int32 | int64](n N) *N { return &n }
