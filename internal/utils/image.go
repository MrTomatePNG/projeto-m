package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

const (
	TargetWidth  = 800
	TargetHeight = 1200
	JPEGQuality  = 70
)

func ResizeImage(b []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	resizedImg := resize.Resize(uint(TargetWidth), 0, img, resize.Lanczos3)
	return resizedImg, nil
}

func PadToCanvas(img image.Image, targetWidth, targetHeight int, bg color.Color) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	srcB := img.Bounds()
	srcW := srcB.Dx()
	srcH := srcB.Dy()
	if srcW <= 0 || srcH <= 0 {
		return canvas
	}

	offsetX := (targetWidth - srcW) / 2
	offsetY := (targetHeight - srcH) / 2
	dst := image.Rect(offsetX, offsetY, offsetX+srcW, offsetY+srcH)
	draw.Draw(canvas, dst, img, srcB.Min, draw.Over)

	return canvas
}

func CenterCropTo(img image.Image, targetWidth, targetHeight int) *image.RGBA {
	srcB := img.Bounds()
	srcW := srcB.Dx()
	srcH := srcB.Dy()
	if srcW <= 0 || srcH <= 0 {
		return image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	}

	scaleW := float64(targetWidth) / float64(srcW)
	scaleH := float64(targetHeight) / float64(srcH)
	scale := scaleW
	if scaleH > scale {
		scale = scaleH
	}

	newW := int(float64(srcW) * scale)
	newH := int(float64(srcH) * scale)
	if newW < targetWidth {
		newW = targetWidth
	}
	if newH < targetHeight {
		newH = targetHeight
	}

	resized := resize.Resize(uint(newW), uint(newH), img, resize.Lanczos3)
	rb := resized.Bounds()
	offsetX := (rb.Dx() - targetWidth) / 2
	offsetY := (rb.Dy() - targetHeight) / 2
	srcCrop := image.Rect(offsetX, offsetY, offsetX+targetWidth, offsetY+targetHeight)
	srcCrop = srcCrop.Intersect(rb)

	out := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.Draw(out, out.Bounds(), resized, srcCrop.Min, draw.Src)
	return out
}

func CreateCanvasWithImage(resizedImg image.Image) *image.RGBA {
	newHeight := int(float64(TargetWidth) * float64(resizedImg.Bounds().Dy()) / float64(resizedImg.Bounds().Dx()))
	newImg := image.NewRGBA(image.Rect(0, 0, TargetWidth, TargetHeight))
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	// Centralizar imagem
	offSetY := (TargetHeight - newHeight) / 2
	draw.Draw(newImg, image.Rect(0, offSetY, TargetWidth, offSetY+newHeight), resizedImg, image.Pt(0, 0), draw.Over)

	return newImg
}

func ProcessImageToVerticalCanvasRGBA(b []byte) (*image.RGBA, error) {
	resized, err := ResizeImage(b)
	if err != nil {
		return nil, err
	}
	return CreateCanvasWithImage(resized), nil
}

func EncodeJPEG(filePath string, img *image.RGBA) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer outFile.Close()

	opts := &jpeg.Options{Quality: JPEGQuality}
	return jpeg.Encode(outFile, img, opts)
}
