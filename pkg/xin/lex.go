package xin

import (
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	tkOpenParen = iota
	tkCloseParen

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
}

func (tk token) String() string {
	switch tk.kind {
	case tkOpenParen:
		return "("
	case tkCloseParen:
		return ")"
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
	filePath string
	line     int
	col      int
}

func bufToToken(s string, pos position) token {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		hexPart := s[2:]
		if _, err := strconv.ParseInt(hexPart, 16, 64); err == nil {
			return token{
				kind:     tkNumberLiteralHex,
				value:    hexPart,
				position: pos,
			}
		} else {
			return token{
				kind: tkNumberLiteralHex,
				// TODO: make this throw a parse err
				value:    "0",
				position: pos,
			}
		}
	} else if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return token{
			kind:     tkNumberLiteralInt,
			value:    s,
			position: pos,
		}
	} else if _, err := strconv.ParseFloat(s, 64); err == nil {
		return token{
			kind:     tkNumberLiteralDecimal,
			value:    s,
			position: pos,
		}
	} else {
		return token{
			kind:     tkName,
			value:    s,
			position: pos,
		}
	}
}

// TODO: lexer should be able to throw parse errors
func lex(r io.Reader) tokenStream {
	toks := make([]token, 0)
	rdr, err := newReader(r)
	if err != nil {
		return toks
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

	return toks
}
