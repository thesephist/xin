package xin

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func charFromEscaper(escaper byte) rune {
	switch escaper {
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case '\\':
		return '\\'
	case '\'':
		return '\''
	default:
		return '?'
	}
}

func escapeString(s string) string {
	builder := strings.Builder{}
	max := len(s)

	for i := 0; i < max; i++ {
		c := s[i]
		if c == '\\' {
			i++
			next := s[i]
			if next == 'x' {
				hex := s[i+1 : i+3]
				i += 2

				codepoint, err := strconv.ParseInt(hex, 16, 32)
				fmt.Println("codepoint number:", hex, codepoint)
				if err != nil || !utf8.ValidRune(rune(codepoint)) {
					builder.WriteRune('?')
					continue
				}

				builder.WriteRune(rune(codepoint))
			} else {
				builder.WriteRune(charFromEscaper(next))
			}
		} else {
			builder.WriteByte(c)
		}
	}

	return builder.String()
}
