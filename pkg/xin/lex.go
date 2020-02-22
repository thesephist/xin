package xin

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tkOpenParen = iota
	tkCloseParen

	tkBindForm
	tkIfForm
	tkDoForm
	tkImportForm

	tkName
	tkNumberLiteralInt
	tkNumberLiteralDecimal
	tkNumberLiteralHex
	tkStringLiteral
)

type tokenKind int

type token struct {
	kind  tokenKind
	value string
	position

	// we compute the integer and float values of
	// number literals at parse time for runtime efficiency
	intv  IntValue
	fracv FracValue
}

func (tk token) String() string {
	switch tk.kind {
	case tkOpenParen:
		return "("
	case tkCloseParen:
		return ")"
	case tkBindForm:
		return ":"
	case tkIfForm:
		return "if"
	case tkDoForm:
		return "do"
	case tkImportForm:
		return "import"
	case tkName, tkNumberLiteralInt, tkNumberLiteralDecimal, tkNumberLiteralHex:
		return tk.value
	case tkStringLiteral:
		return "'" + tk.value + "'"
	default:
		return "unknown token"
	}
}

type tokenStream []token

func (toks tokenStream) String() string {
	tokstrings := make([]string, len(toks))
	for i, tk := range toks {
		tokstrings[i] = tk.String()
	}
	return strings.Join(tokstrings, " ")
}

type position struct {
	path string
	line int
	col  int
}

func (p position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.path, p.line, p.col)
}

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
		return rune(escaper)
	}
}

func escapeString(s string) string {
	builder := strings.Builder{}
	max := len(s)

	for i := 0; i < max; i++ {
		c := s[i]
		if c == '\\' {
			i++

			if i >= len(s) {
				return builder.String()
			}
			next := s[i]

			if next == 'x' {
				hex := s[i+1 : i+3]
				i += 2

				codepoint, err := strconv.ParseInt(hex, 16, 32)
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

func bufToToken(s string, pos position) token {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		hexPart := s[2:]
		if _, err := strconv.ParseInt(hexPart, 16, 64); err == nil {
			v, _ := strconv.ParseInt(hexPart, 16, 64)
			return token{
				kind:     tkNumberLiteralHex,
				value:    hexPart,
				position: pos,
				intv:     IntValue(v),
			}
		} else {
			return token{
				kind: tkNumberLiteralHex,
				// TODO: make this throw a parse err
				value:    "0",
				position: pos,
				intv:     IntValue(0),
			}
		}
	} else if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		v, _ := strconv.ParseInt(s, 10, 64)
		return token{
			kind:     tkNumberLiteralInt,
			value:    s,
			position: pos,
			intv:     IntValue(v),
		}
	} else if _, err := strconv.ParseFloat(s, 64); err == nil {
		v, _ := strconv.ParseFloat(s, 64)
		return token{
			kind:     tkNumberLiteralDecimal,
			value:    s,
			position: pos,
			fracv:    FracValue(v),
		}
	} else {
		switch s {
		case ":":
			return token{
				kind:     tkBindForm,
				position: pos,
			}
		case "if":
			return token{
				kind:     tkIfForm,
				position: pos,
			}
		case "do":
			return token{
				kind:     tkDoForm,
				position: pos,
			}
		case "import":
			return token{
				kind:     tkImportForm,
				position: pos,
			}
		default:
			return token{
				kind:     tkName,
				value:    s,
				position: pos,
			}
		}
	}
}

func lex(path string, r io.Reader) (tokenStream, InterpreterError) {
	toks := make([]token, 0)
	rdr, err := newReader(path, r)
	if err != nil {
		return toks, nil
	}

	buf := ""
	clear := func() {
		if buf != "" {
			toks = append(toks, bufToToken(buf, rdr.currPos()))
			buf = ""
		}
	}
	for !rdr.done() {
		peeked := rdr.peek()
		switch {
		case peeked == ";":
			clear()
			rdr.upto("\n")
			rdr.skip()
		case peeked == "'":
			clear()
			rdr.skip()

			content := rdr.upto("'")
			for rdr.lookback() == "\\" {
				// count preceding backslashes
				// if there's an even number, like \\'
				// or \\\\', break loop
				backslashesBefore := 1
				for i := 2; i < rdr.index; i++ {
					if rdr.before(i) == "\\" {
						backslashesBefore++
					} else {
						break
					}
				}
				if backslashesBefore%2 == 0 {
					break
				}

				rdr.skip()
				content += "'"
				content += rdr.upto("'")
			}

			toks = append(toks, token{
				kind:     tkStringLiteral,
				value:    escapeString(content),
				position: rdr.currPos(),
			})
			rdr.skip()
		case peeked == "(":
			clear()
			toks = append(toks, token{
				kind:     tkOpenParen,
				position: rdr.currPos(),
			})
			rdr.skip()
		case peeked == ")":
			clear()
			toks = append(toks, token{
				kind:     tkCloseParen,
				position: rdr.currPos(),
			})
			rdr.skip()
		case unicode.IsSpace([]rune(peeked)[0]):
			clear()
			rdr.skip()
		default:
			buf += rdr.next()
		}
	}

	// clear out remaining items in the buffer
	clear()

	return toks, nil
}
