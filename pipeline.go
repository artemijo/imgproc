package main

import (
	"fmt"
	"image"
)

// Step represents a single image processing operation.
type Step interface {
	Name() string
	Process(img image.Image) (image.Image, error)
}

// Pipeline chains multiple Steps and executes them sequentially.
type Pipeline struct {
	steps []Step
}

func NewPipeline(steps ...Step) *Pipeline {
	return &Pipeline{steps: steps}
}

func (p *Pipeline) Run(img image.Image) (image.Image, error) {
	var err error
	for _, step := range p.steps {
		img, err = step.Process(img)
		if err != nil {
			return nil, fmt.Errorf("step %q: %w", step.Name(), err)
		}
	}
	return img, nil
}
