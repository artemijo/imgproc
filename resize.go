package main

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// ResizeTo900 resizes the image so the wide side is 900 pixels.
// Uses CatmullRom (Lanczos-quality) interpolation.
type ResizeTo900 struct {
	TargetSize int
}

func NewResizeTo900() *ResizeTo900 {
	return &ResizeTo900{TargetSize: 900}
}

func (r *ResizeTo900) Name() string { return "resize_900" }

func (r *ResizeTo900) Process(img image.Image) (image.Image, error) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	wideSide := w
	if h > w {
		wideSide = h
	}

	if wideSide == r.TargetSize {
		return img, nil
	}

	scale := float64(r.TargetSize) / float64(wideSide)
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)

	return dst, nil
}
