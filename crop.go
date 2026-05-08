package main

import (
	"image"
)

// CropTo32 crops the image to a 3:2 aspect ratio by taking the largest
// centered 3:2 rectangle from the source.
type CropTo32 struct{}

func NewCropTo32() *CropTo32 {
	return &CropTo32{}
}

func (c *CropTo32) Name() string { return "crop_3:2" }

func (c *CropTo32) Process(img image.Image) (image.Image, error) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	targetW, targetH := w, h

	if w*2 > h*3 {
		// Wider than 3:2 — constrain by height, crop sides.
		targetW = h * 3 / 2
	} else if w*2 < h*3 {
		// Taller than 3:2 — constrain by width, crop top/bottom.
		targetH = w * 2 / 3
	} else {
		// Already 3:2.
		return img, nil
	}

	x0 := b.Min.X + (w-targetW)/2
	y0 := b.Min.Y + (h-targetH)/2
	rect := image.Rect(x0, y0, x0+targetW, y0+targetH)

	type subImager interface {
		SubImage(image.Rectangle) image.Image
	}
	if si, ok := img.(subImager); ok {
		return si.SubImage(rect), nil
	}

	dst := image.NewRGBA(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst, nil
}
