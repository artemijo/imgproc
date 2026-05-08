package main

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// ResizeFixed resizes the image to a fixed width x height.
// Uses CatmullRom (Lanczos-quality) interpolation.
type ResizeFixed struct {
	Width  int
	Height int
}

func NewResizeFixed(w, h int) *ResizeFixed {
	return &ResizeFixed{Width: w, Height: h}
}

func (r *ResizeFixed) Name() string { return "resize_fixed" }

func (r *ResizeFixed) Process(img image.Image) (image.Image, error) {
	b := img.Bounds()
	if b.Dx() == r.Width && b.Dy() == r.Height {
		return img, nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)
	return dst, nil
}
