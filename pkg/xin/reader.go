package xin

import (
	"io"
	"io/ioutil"
)

type reader struct {
	source string
	index  int
	max    int

	position
}

func newReader(path string, r io.Reader) (*reader, error) {
	allBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	asString := string(allBytes)
	rdr := reader{
		source: asString,
		max:    len(asString),
	}
	rdr.position.path = path
	return &rdr, nil
}
func (rdr *reader) done() bool {
	return rdr.index >= rdr.max
}

func (rdr *reader) next() string {
	rdr.skip()
	return rdr.source[rdr.index-1 : rdr.index]
}

func (rdr *reader) skip() {
	rdr.index++
	if rdr.done() {
		return
	}

	if rdr.source[rdr.index] == '\n' {
		rdr.line++
		rdr.col = 0
	} else {
		rdr.col++
	}
}

func (rdr *reader) lookback() string {
	return rdr.source[rdr.index-1 : rdr.index]
}

func (rdr *reader) before(n int) string {
	if n > rdr.index {
		panic("Tried to look before() start of string")
	}

	return rdr.source[rdr.index-n : rdr.index-n+1]
}

func (rdr *reader) peek() string {
	return rdr.source[rdr.index : rdr.index+1]
}

func (rdr *reader) upto(end string) string {
	s := ""
	for !rdr.done() && rdr.peek() != end {
		s += rdr.next()
	}
	return s
}

func (rdr *reader) currPos() position {
	return position{
		path: rdr.path,
		line: rdr.position.line + 1,
		col:  rdr.position.col + 1,
	}
}
