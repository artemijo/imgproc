package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	inputDir := flag.String("input", filepath.Join(exeDir, "input"), "Input folder with JPG files")
	outputDir := flag.String("output", filepath.Join(exeDir, "output"), "Output folder for processed files")
	quality := flag.Int("quality", 95, "JPEG output quality (1-100)")
	flag.Parse()

	if *quality < 1 || *quality > 100 {
		fmt.Fprintf(os.Stderr, "Error: quality must be 1-100, got %d\n", *quality)
		os.Exit(1)
	}

	info, err := os.Stat(*inputDir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: input %q is not a valid directory\n", *inputDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot create output dir %q: %v\n", *outputDir, err)
		os.Exit(1)
	}

	pipeline := NewPipeline(NewCropTo43(), NewAutoCorrect(), NewResizeTo900())

	files, err := collectJPGs(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning input dir: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		log.Println("No JPG files found in input directory")
		os.Exit(0)
	}

	log.Printf("Processing %d files: %s → %s (quality %d)", len(files), *inputDir, *outputDir, *quality)

	ok, failed := 0, 0
	for i, path := range files {
		name := filepath.Base(path)
		log.Printf("[%d/%d] %s", i+1, len(files), name)

		if err := processFile(path, *outputDir, *quality, pipeline); err != nil {
			log.Printf("[%d/%d] %s: FAILED — %v", i+1, len(files), name, err)
			failed++
		} else {
			ok++
		}
	}

	log.Printf("Done: %d ok, %d failed", ok, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func collectJPGs(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func processFile(inputPath, outputDir string, quality int, pipe *Pipeline) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	result, err := pipe.Run(img)
	if err != nil {
		return err
	}

	outName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + ".jpg"
	outPath := filepath.Join(outputDir, outName)

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, result, &jpeg.Options{Quality: quality}); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
