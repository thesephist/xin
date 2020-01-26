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

func (rdr *reader) done() bool {
	return rdr.index >= rdr.max
}

func (rdr *reader) next() string {
	rdr.index++
	return rdr.source[rdr.index-1 : rdr.index]
}

func (rdr *reader) skip() {
	rdr.index++
}

func (rdr *reader) lookback() string {
	return rdr.source[rdr.index-1 : rdr.index]
}

func (rdr *reader) peek() string {
	return rdr.source[rdr.index : rdr.index+1]
}

func (rdr *reader) upto(end string) string {
	s := ""
	for rdr.peek() != end {
		s += rdr.next()
	}
	return s
}

func (rdr *reader) currPos() position {
	return rdr.position
}

func newReader(r io.Reader) (reader, error) {
	allBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return reader{}, err
	}

	asString := string(allBytes)
	rdr := reader{
		source: asString,
		max:    len(asString),
	}
	return rdr, nil
}
