package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBar(t *testing.T) {
	bar := new(Bar)
	stdout := os.Stdout

	bar.Init(200, true)
	require.Equal(t, int64(200), bar.total)

	bar.Progress(50)
	require.Equal(t, float64(25), bar.state.percent)
	require.Equal(t, int64(50), bar.state.current)
	require.Equal(t, 25, bar.state.filledSteps)
	require.Equal(t, 75, bar.state.remainSteps)

	bar.Progress(140)
	require.Equal(t, float64(95), bar.state.percent)
	require.Equal(t, int64(190), bar.state.current)
	require.Equal(t, 95, bar.state.filledSteps)
	require.Equal(t, 5, bar.state.remainSteps)

	require.NotPanics(t, func() {
		bar.Progress(1000)
	})

	require.Equal(t, float64(100), bar.state.percent)
	require.Equal(t, int64(200), bar.state.current)
	require.Equal(t, 100, bar.state.filledSteps)
	require.Equal(t, 0, bar.state.remainSteps)

	err := bar.Progress(-1)
	require.Truef(t, errors.Is(err, ErrInvalidProgressStep), "actual error %q", err)

	bar.Finish()
	require.Equal(t, stdout, os.Stdout)
}
