package main

import (
	"image"
)

// CropTo43 crops the image to a 4:3 aspect ratio by taking the largest
// centered 4:3 rectangle from the source.
type CropTo43 struct{}

func NewCropTo43() *CropTo43 {
	return &CropTo43{}
}

func (c *CropTo43) Name() string { return "crop_4:3" }

func (c *CropTo43) Process(img image.Image) (image.Image, error) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	targetW, targetH := w, h

	if w*3 > h*4 {
		// Wider than 4:3 — constrain by height, crop sides.
		targetW = h * 4 / 3
	} else if w*3 < h*4 {
		// Taller than 4:3 — constrain by width, crop top/bottom.
		targetH = w * 3 / 4
	} else {
		// Already 4:3.
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

	// Fallback: manual copy.
	dst := image.NewRGBA(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst, nil
}
