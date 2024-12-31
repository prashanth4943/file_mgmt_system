package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os/exec"

	"github.com/h2non/filetype"
	"github.com/nfnt/resize"
)

func (d *DownloadService) GenerateThumbnail(fileBytes []byte) ([]byte, error) {

	kind, err := filetype.Match(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("could not determine file type: %w", err)
	}

	var thumbnailBytes []byte

	switch kind.MIME.Value {
	case "image/jpeg", "image/png":
		// Generate thumbnail for images
		thumbnailBytes, err = generateImageThumbnail(fileBytes)
	case "application/pdf":
		// Generate thumbnail for PDFs
		thumbnailBytes, err = generatePDFThumbnail(fileBytes)
	case "video/mp4", "video/quicktime": // video/quicktime covers MOV files
		// Generate thumbnail for videos
		thumbnailBytes, err = generateVideoThumbnail(fileBytes)
	default:
		// Unsupported file types
		return nil, fmt.Errorf("unsupported file type: %s", kind.MIME.Value)
	}

	if err != nil {
		return nil, err
	}

	return thumbnailBytes, nil
}

func generateImageThumbnail(fileBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	thumbnail := resize.Resize(100, 100, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumbnail, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func generatePDFThumbnail(fileBytes []byte) ([]byte, error) {
	cmd := exec.Command("convert", "pdf:-[0]", "-thumbnail", "100x100", "jpeg:-")
	cmd.Stdin = bytes.NewReader(fileBytes)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func generateVideoThumbnail(fileBytes []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "thumbnail,scale=100:100", "-frames:v", "1", "-f", "image2pipe", "-vcodec", "mjpeg", "-")
	cmd.Stdin = bytes.NewReader(fileBytes)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
