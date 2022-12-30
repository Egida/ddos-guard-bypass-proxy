package common

import (
	"bytes"
	"encoding/base64"
	"io"
)

func copySlice(values []string) []string {
	out := make([]string, len(values))

	copy(out, values)

	return out
}

func CopyHeader(header map[string][]string) map[string][]string {
	out := map[string][]string{}

	for k, v := range header {
		out[k] = copySlice(v)
	}

	return out
}

func MarshalBody(reader io.Reader) (string, error) {
	if reader == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	body := base64.StdEncoding.EncodeToString(bodyBytes)

	return body, err
}

func UnmarshalBody(data string) (io.ReadCloser, error) {
	if data == "" {
		return nil, nil
	}

	bodyBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	body := io.NopCloser(bytes.NewReader(bodyBytes))

	return body, nil
}
