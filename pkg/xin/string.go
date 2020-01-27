package xin

import (
	"bytes"
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

type StringValue []byte

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Equal(o Value) bool {
	if ov, ok := o.(StringValue); ok {
		return bytes.Compare(v, ov) == 0
	}

	return false
}

func strGetForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstStr) {
			return firstStr[secondInt : secondInt+1], nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func strSizeForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstString, ok := first.(StringValue); ok {
		return IntValue(len(firstString)), nil
	}

	return IntValue(1), nil
}

func strSliceForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
			required: 3,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}
	third, err := unlazy(args[2])
	if err != nil {
		return nil, err
	}

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	thirdInt, tok := third.(IntValue)
	if fok && sok && tok {
		max := len(firstStr)
		inRange := func(iv IntValue) bool {
			return int(iv) >= 0 && int(iv) < max
		}

		if int(secondInt) >= max {
			secondInt = IntValue(max) - 1
		}
		if int(thirdInt) >= max {
			thirdInt = IntValue(max) - 1
		}

		if inRange(secondInt) && inRange(thirdInt) {
			return firstStr[secondInt:thirdInt], nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
