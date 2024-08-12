package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder
	if len(str) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(str[0])) {
		return "", ErrInvalidString
	}

	slashIsOpen := false
	for ind, chr := range str {
		lastChar := ""
		if builder.Len() > 0 {
			lastChar = string(builder.String()[builder.Len()-1])
		}

		switch {
		case unicode.IsLetter(chr) || chr == '\n':
			builder.WriteString(string(chr))
		case chr == '\\':
			if !slashIsOpen {
				slashIsOpen = true
			} else {
				builder.WriteString(string(chr))
				slashIsOpen = false
			}
		case unicode.IsDigit(chr):
			switch {
			case chr == '0':
				// if digit is 0 then remove last letter from resulting string
				removeLastChar(&builder)
			case slashIsOpen:
				builder.WriteString(string(chr))
				slashIsOpen = false
			case len(str) > ind+1 && unicode.IsDigit(rune(str[ind+1])):
				return "", ErrInvalidString
			default:
				nTimes, err := strconv.Atoi(string(chr))
				if err == nil {
					builder.WriteString(strings.Repeat(lastChar, nTimes-1))
				}
			}
		}
	}

	return builder.String(), nil
}

func removeLastChar(b *strings.Builder) {
	if b.Len() == 0 {
		return
	}
	str := b.String()
	str = str[:len(str)-1]
	b.Reset()
	b.WriteString(str)
}
