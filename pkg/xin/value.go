package xin

import (
	"strconv"
)

type Value interface {
	String() string
}

type StringValue string

func (v StringValue) String() string {
	return "'" + string(v) + "'"
}

type IntValue int64

func (v IntValue) String() string {
	return strconv.FormatInt(int64(v), 10)
}
