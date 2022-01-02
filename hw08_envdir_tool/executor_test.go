package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		testName  string
		command   string
		envValues Environment
		osEnvs    map[string]string
	}{
		{
			testName:  "without env",
			command:   "echo",
			envValues: nil,
		},
		{
			testName: "single env",
			command:  "echo",
			envValues: Environment{
				"BAR": EnvValue{
					Value: "bar",
				},
			},
		},
		{
			testName: "multiple env",
			command:  "echo",
			envValues: Environment{
				"BAR": EnvValue{
					Value: "bar",
				},
				"FOO": EnvValue{
					Value: "foo",
				},
			},
		},
		{
			testName: "unset env",
			command:  "echo",
			envValues: Environment{
				"UNSET": EnvValue{
					NeedRemove: true,
				},
			},
			osEnvs: map[string]string{"UNSET": "UNSET"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			for name, value := range tc.osEnvs {
				err := os.Setenv(name, value)
				require.NoError(t, err)
			}

			command := []string{tc.command}
			exitCode := RunCmd(command, tc.envValues)
			require.Equal(t, 0, exitCode)

			for name, value := range tc.envValues {
				require.Equal(t, os.Getenv(name), value.Value)
			}
		})
	}
}

func TestRunCmdUnsuccessfulStatus(t *testing.T) {
	f, err := os.Create("testdata/permitted")
	defer os.Remove("testdata/permitted")
	require.NoError(t, err)

	err = f.Chmod(0000)
	require.NoError(t, err)

	cmd := []string{"cat", "testdata/permitted"}
	code := RunCmd(cmd, Environment{})
	require.Equal(t, 1, code)
}
