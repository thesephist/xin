package xin

import (
	"strconv"
)

type Value interface {
	String() string
}

type String string

func (s String) String() string {
	return string(s)
}

type Int int

func (i Int) String() string {
	return strconv.Itoa(int(i))
}
