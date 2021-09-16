package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var R = regexp.MustCompile(`(^\d)|([^\\]\d\d)|(\\[a-z]|\\[A-Z])|` + "(`)")

func Unpack(input string) (string, error) {
	if !validate(input) {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	runes := []rune(input)

	for i := 0; i <= len(runes)-1; i++ {
		char := string(runes[i])

		if char == `\` {
			builder.WriteRune(runes[i+1])
			i++

			continue
		}

		if unicode.IsDigit(runes[i]) {
			repeatAmount, err := strconv.Atoi(char)
			if err != nil {
				return "", err
			}
			repeatLastChar(&builder, repeatAmount)
			continue
		}
		builder.WriteString(char)
	}

	return builder.String(), nil
}

func repeatLastChar(builder *strings.Builder, repeatAmount int) {
	currentString := []rune(builder.String())
	lastChar := string(currentString[len(currentString)-1])

	if repeatAmount > 0 {
		builder.WriteString(strings.Repeat(lastChar, repeatAmount-1))
	} else {
		// if repeatAmount == 0 we need to cut previous char
		builder.Reset()
		builder.WriteString(string(currentString[:len(currentString)-1]))
	}
}

func validate(input string) bool {
	return len(R.FindStringSubmatch(input)) == 0
}
