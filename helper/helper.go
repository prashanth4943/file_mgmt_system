package helper

import (
	"bytes"
	"compress/gzip"
	"context"
	"file_mgmt_system/middleware"
	"io"
)

func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(middleware.UserKey).(string)
	return email, ok
}

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
