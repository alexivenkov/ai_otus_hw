package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testCases := []struct {
		expectedName       string
		expectedValue      string
		expectedNeedRemove bool
	}{
		{
			expectedName:       "BAR",
			expectedValue:      "bar",
			expectedNeedRemove: false,
		},
		{
			expectedName:       "EMPTY",
			expectedValue:      "",
			expectedNeedRemove: false,
		},
		{
			expectedName:       "FOO",
			expectedValue:      "   foo\nwith new line",
			expectedNeedRemove: false,
		},
		{
			expectedName:       "HELLO",
			expectedValue:      "\"hello\"",
			expectedNeedRemove: false,
		},
		{
			expectedName:       "UNSET",
			expectedNeedRemove: true,
		},
	}

	env, err := ReadDir("testdata/env")
	require.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.expectedName, func(t *testing.T) {
			envValue, ok := env[tc.expectedName]

			require.True(t, ok)
			require.Equal(t, tc.expectedValue, envValue.Value)
			require.Equal(t, tc.expectedNeedRemove, envValue.NeedRemove)
		})
	}
}

func TestEmptyDir(t *testing.T) {
	err := os.Mkdir("testdata/empty", 0777)
	require.NoError(t, err)
	defer os.Remove("testdata/empty")

	env, err := ReadDir("testdata/empty")
	require.NoError(t, err)
	require.Equal(t, Environment{}, env)
}

func TestUnreadableFile(t *testing.T) {
	err := os.Mkdir("testdata/test", 0777)
	require.NoError(t, err)

	f, err := os.Create("testdata/test/RESTRICTED")
	require.NoError(t, err)

	err = f.Chmod(0000)
	require.NoError(t, err)

	defer os.RemoveAll("testdata/test")

	_, err = ReadDir("testdata/test")
	require.ErrorIs(t, err, os.ErrPermission)
}
