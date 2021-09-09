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
loop:
	for i := 0; i <= len(inRune)-1; {
		switch true {
		case unicode.IsDigit(inRune[i]):
			return "", ErrInvalidString
		case string(inRune[i]) == `\`:
			switch true {
			case i == len(inRune)-1:
				return "", ErrInvalidString
			case !unicode.IsDigit(inRune[i+1]) && string(inRune[i+1]) != `\`:
				return "", ErrInvalidString
			case i+2 == len(inRune):
				out.WriteString(string(inRune[i+1]))
				break loop
			case unicode.IsDigit(inRune[i+2]):
				count, _ := strconv.Atoi(string(inRune[i+2]))
				out.WriteString(strings.Repeat(string(inRune[i+1]), count))
				i += 3
				continue loop
			default:
				out.WriteString(string(inRune[i+1]))
				i += 2
				continue loop
			}
		case i+1 <= len(inRune):
			switch true {
			case i+1 == len(inRune):
				out.WriteString(string(inRune[i]))
				break loop
			case unicode.IsDigit(inRune[i+1]):
				count, _ := strconv.Atoi(string(inRune[i+1]))
				out.WriteString(strings.Repeat(string(inRune[i]), count))
				i += 2
			default:
				out.WriteRune(inRune[i])
				i++
			}
		}
	}
	return out.String(), nil
}
