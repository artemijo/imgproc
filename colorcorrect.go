package main

import (
	"image"
	"image/color"
	"math"
)

// AutoCorrect applies gray-world white balance and auto-levels (percentile stretch).
type AutoCorrect struct {
	LowPct  float64 // lower percentile for stretch (default 0.005 = 0.5%)
	HighPct float64 // upper percentile for stretch (default 0.995 = 99.5%)
}

func NewAutoCorrect() *AutoCorrect {
	return &AutoCorrect{LowPct: 0.01, HighPct: 0.99}
}

func (a *AutoCorrect) Name() string { return "auto_correct" }

func (a *AutoCorrect) Process(img image.Image) (image.Image, error) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	totalPixels := w * h

	// Pass 1: compute average R/G/B and build histograms.
	var sumR, sumG, sumB uint64
	histR := [256]uint64{}
	histG := [256]uint64{}
	histB := [256]uint64{}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// RGBA() returns premultiplied 0-65535, convert to 0-255.
			ri := r >> 8
			gi := g >> 8
			bi := b >> 8
			sumR += uint64(ri)
			sumG += uint64(gi)
			sumB += uint64(bi)
			histR[ri]++
			histG[gi]++
			histB[bi]++
		}
	}

	// Gray world: scale channels so each matches the average brightness.
	avgR := float64(sumR) / float64(totalPixels)
	avgG := float64(sumG) / float64(totalPixels)
	avgB := float64(sumB) / float64(totalPixels)
	avgGray := (avgR + avgG + avgB) / 3.0

	// Gray world with damping: apply only 40% of the correction.
	const wbStrength = 0.4
	wbR, wbG, wbB := 1.0, 1.0, 1.0
	if avgR > 1.0 {
		wbR = 1.0 + (avgGray/avgR-1.0)*wbStrength
	}
	if avgG > 1.0 {
		wbG = 1.0 + (avgGray/avgG-1.0)*wbStrength
	}
	if avgB > 1.0 {
		wbB = 1.0 + (avgGray/avgB-1.0)*wbStrength
	}

	// Apply white balance to histograms to find percentiles on corrected data.
	// Build lookup: for each original value, compute WB-corrected value.
	wbLookup := func(val uint8, wb float64) uint8 {
		v := float64(val) * wb
		if v > 255 {
			v = 255
		}
		return uint8(v)
	}

	// Build corrected histograms.
	corrHistR := [256]uint64{}
	corrHistG := [256]uint64{}
	corrHistB := [256]uint64{}
	for i := 0; i < 256; i++ {
		corrHistR[wbLookup(uint8(i), wbR)] += histR[i]
		corrHistG[wbLookup(uint8(i), wbG)] += histG[i]
		corrHistB[wbLookup(uint8(i), wbB)] += histB[i]
	}

	// Find percentiles.
	lowThresh := float64(totalPixels) * a.LowPct
	highThresh := float64(totalPixels) * a.HighPct

	pctlR := percentile(corrHistR[:], lowThresh, highThresh)
	pctlG := percentile(corrHistG[:], lowThresh, highThresh)
	pctlB := percentile(corrHistB[:], lowThresh, highThresh)

	// Pass 2: apply white balance + auto-levels.
	dst := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			ri := r >> 8
			gi := g >> 8
			bi := b >> 8
			ai := a >> 8

			ro := stretch(float64(ri)*wbR, float64(pctlR[0]), float64(pctlR[1]))
			go_ := stretch(float64(gi)*wbG, float64(pctlG[0]), float64(pctlG[1]))
			bo := stretch(float64(bi)*wbB, float64(pctlB[0]), float64(pctlB[1]))

			dst.SetRGBA(x, y, color.RGBA{R: ro, G: go_, B: bo, A: uint8(ai)})
		}
	}

	return dst, nil
}

// percentile returns (low, high) values at the given cumulative thresholds.
func percentile(hist []uint64, lowThresh, highThresh float64) [2]uint8 {
	var cum uint64
	lowVal, highVal := uint8(0), uint8(255)
	lowFound := false

	for i, count := range hist {
		cum += count
		if !lowFound && float64(cum) >= lowThresh {
			lowVal = uint8(i)
			lowFound = true
		}
		if float64(cum) >= highThresh {
			highVal = uint8(i)
			break
		}
	}

	if highVal <= lowVal {
		highVal = lowVal + 1
	}
	return [2]uint8{lowVal, highVal}
}

// stretch maps val from [lo, hi] to [0, 255] and clamps.
func stretch(val, lo, hi float64) uint8 {
	if hi <= lo {
		return uint8(math.Min(math.Max(val, 0), 255))
	}
	out := (val - lo) / (hi - lo) * 255.0
	return uint8(math.Min(math.Max(out, 0), 255))
}
