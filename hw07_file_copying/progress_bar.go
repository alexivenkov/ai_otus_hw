package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidProgressStep = errors.New("progress step should be positive integer")

type State struct {
	percent     float64
	current     int64
	filledSteps int
	remainSteps int
}

type Bar struct {
	total          int64
	state          State
	originalStdout *os.File
}

func (b *Bar) Init(total int64, disableOutput bool) {
	b.total = total
	b.originalStdout = os.Stdout

	if disableOutput {
		os.Stdout = nil
	}
}

func (b *Bar) Reset() {
	b.total = 0
	b.state = State{}
}

func (b *Bar) Progress(step int64) error {
	if step < 0 {
		return ErrInvalidProgressStep
	}

	if b.state.current+step > b.total {
		b.state.current = b.total
		b.state.percent = 100
		b.state.filledSteps = 100
		b.state.remainSteps = 0

		return nil
	}

	b.state.current += step
	b.state.percent = float64(b.state.current) / float64(b.total) * 100
	b.state.filledSteps = int(b.state.percent)
	b.state.remainSteps = 100 - b.state.filledSteps

	return nil
}

func (b *Bar) Print() {
	fmt.Printf("\r:[")

	for i := 0; i < b.state.filledSteps; i++ {
		fmt.Print("=")
	}

	for i := 0; i < b.state.remainSteps; i++ {
		fmt.Print(" ")
	}

	fmt.Printf("]: %.2f%%", b.state.percent)
}

func (b *Bar) Finish() {
	fmt.Printf("\n")
	os.Stdout = b.originalStdout
}
