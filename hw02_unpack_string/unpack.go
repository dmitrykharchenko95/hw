package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	inRune := []rune(in)
	if in == "" {
		return "", nil
	}
	var out strings.Builder
	for i := 0; i <= len(inRune)-1; {
		if unicode.IsDigit(inRune[i]) {
			return "", ErrInvalidString
		}
		if string(inRune[i]) == `\` {
			if !unicode.IsDigit(inRune[i+1]) && string(inRune[i+1]) != `\` {
				return "", ErrInvalidString
			} else if i+2 == len(inRune) {
				out.WriteString(string(inRune[i+1]))
				break
			} else if unicode.IsDigit(inRune[i+2]) {
				count, _ := strconv.Atoi(string(in[i+2]))
				out.WriteString(strings.Repeat(string(inRune[i+1]), count))
				i += 3
			} else {
				out.WriteString(string(inRune[i+1]))
				i += 2
			}
		}
		if i+1 <= len(inRune) {
			if i+1 == len(inRune) {
				out.WriteString(string(inRune[i]))
				break
			} else if unicode.IsDigit(inRune[i+1]) {
				count, _ := strconv.Atoi(string(inRune[i+1]))
				out.WriteString(strings.Repeat(string(inRune[i]), count))
				i += 2
			} else {
				out.WriteRune(inRune[i])
				i += 1
			}

		}

	}
	return out.String(), nil
}
