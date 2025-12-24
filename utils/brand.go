package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type BrandPosition string

const (
	TopLeft     BrandPosition = "top-left"
	TopRight    BrandPosition = "top-right"
	BottomLeft  BrandPosition = "bottom-left"
	BottomRight BrandPosition = "bottom-right"
	Center      BrandPosition = "center"
)

type BrandingOptions struct {
	Position  BrandPosition
	BrandLogo string
	Width     int
	Height    int
}

func downloadImage(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func decodeImage(data []byte, ext string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch strings.ToLower(ext) {
	case "jpg", "jpeg":
		return jpeg.Decode(reader)
	case "png":
		return png.Decode(reader)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
}

func encodeImage(img image.Image, ext string) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch strings.ToLower(ext) {
	case "jpg", "jpeg":
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
		if err != nil {
			return nil, err
		}
	case "png":
		err := png.Encode(buf, img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	return buf.Bytes(), nil
}

func BrandImage(imageBuffer []byte, options BrandingOptions, fileExt string) ([]byte, error) {
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}

	tempFileName := fmt.Sprintf("%s.%s", uuid.New().String(), fileExt)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	logoExt := filepath.Ext(options.BrandLogo)
	if logoExt != "" {
		logoExt = strings.TrimPrefix(logoExt, ".")
	}
	brandLogoFileName := fmt.Sprintf("%s_logo.%s", uuid.New().String(), logoExt)
	brandLogoPath := filepath.Join(tempDir, brandLogoFileName)

	defer func() {
		os.Remove(tempFilePath)
		os.Remove(brandLogoPath)
	}()

	mainImg, err := decodeImage(imageBuffer, fileExt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode main image: %w", err)
	}

	logoBuffer, err := downloadImage(options.BrandLogo)
	if err != nil {
		return nil, fmt.Errorf("failed to download logo: %w", err)
	}

	logoImg, err := decodeImage(logoBuffer, logoExt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode logo: %w", err)
	}

	bounds := mainImg.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	maxLogoSize := int(math.Min(float64(imgWidth), float64(imgHeight)) * 0.25)
	requestedLogoSize := int(math.Max(float64(options.Width), float64(options.Height)))

	var finalLogoWidth, finalLogoHeight uint

	if requestedLogoSize > maxLogoSize {
		scaleFactor := float64(maxLogoSize) / float64(requestedLogoSize)
		finalLogoWidth = uint(math.Round(float64(options.Width) * scaleFactor))
		finalLogoHeight = uint(math.Round(float64(options.Height) * scaleFactor))
	} else if requestedLogoSize < int(float64(maxLogoSize)*0.3) {
		scaleFactor := (float64(maxLogoSize) * 0.3) / float64(requestedLogoSize)
		finalLogoWidth = uint(math.Round(float64(options.Width) * scaleFactor))
		finalLogoHeight = uint(math.Round(float64(options.Height) * scaleFactor))
	} else {
		finalLogoWidth = uint(options.Width)
		finalLogoHeight = uint(options.Height)
	}

	resizedLogo := resize.Resize(finalLogoWidth, finalLogoHeight, logoImg, resize.Lanczos3)

	offsetX := 20
	offsetY := 20

	var x, y int
	logoBounds := resizedLogo.Bounds()
	logoWidth := logoBounds.Dx()
	logoHeight := logoBounds.Dy()

	switch options.Position {
	case TopLeft:
		x = offsetX
		y = offsetY
	case TopRight:
		x = imgWidth - logoWidth - offsetX
		y = offsetY
	case BottomLeft:
		x = offsetX
		y = imgHeight - logoHeight - offsetY
	case BottomRight:
		x = imgWidth - logoWidth - offsetX
		y = imgHeight - logoHeight - offsetY
	case Center:
		x = (imgWidth - logoWidth) / 2
		y = (imgHeight - logoHeight) / 2
	default:
		x = imgWidth - logoWidth - offsetX
		y = imgHeight - logoHeight - offsetY
	}

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, mainImg, bounds.Min, draw.Src)

	logoPoint := image.Point{X: x, Y: y}
	logoRect := image.Rectangle{Min: logoPoint, Max: logoPoint.Add(logoBounds.Size())}
	draw.Draw(dst, logoRect, resizedLogo, logoBounds.Min, draw.Over)

	brandedBuffer, err := encodeImage(dst, fileExt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode branded image: %w", err)
	}

	return brandedBuffer, nil
}
