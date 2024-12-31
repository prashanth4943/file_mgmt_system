package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/nfnt/resize"
)

func (d *DownloadService) CompressImage(fileData io.Reader, format string, width uint) ([]byte, error) {
	img, _, err := image.Decode(fileData)
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Resize(width, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if strings.Contains(format, "jpeg") || strings.Contains(format, "jpg") {
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 70}) // Adjust quality (70%)
	} else if strings.Contains(format, "png") {
		err = png.Encode(&buf, resizedImg)
	} else {
		err = fmt.Errorf("unsupported format")
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
