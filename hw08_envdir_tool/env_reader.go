package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var ErrInvalidValue = errors.New("environment variables should not contain = symbols")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment, 10)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		file, err := os.Open(dir + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		if fileInfo.Size() == 0 {
			env[fileInfo.Name()] = EnvValue{
				NeedRemove: true,
			}
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			value, err := prepareValue(scanner.Text())
			if err != nil {
				return nil, err
			}

			envVal := EnvValue{
				Value: value,
			}
			env[fileInfo.Name()] = envVal
			break
		}

		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	return env, nil
}

func prepareValue(value string) (string, error) {
	value = strings.ReplaceAll(value, "\x00", "\n")
	if strings.Contains(value, "=") {
		return "", ErrInvalidValue
	}

	return strings.TrimRight(value, "\t "), nil
}
