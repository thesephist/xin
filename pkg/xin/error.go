package xin

import (
	"fmt"
	"strings"
)

func FormatError(e InterpreterError) string {
	return e.Error() + "\n\t at " + e.pos().String()
}

type InterpreterError interface {
	error
	pos() position
}

type UndefinedNameError struct {
	name     string
	position position
}

func (e UndefinedNameError) Error() string {
	return fmt.Sprintf("Undefined name %s", e.name)
}

func (e UndefinedNameError) pos() position {
	return e.position
}

type InvalidFormError struct {
	position position
}

func (e InvalidFormError) Error() string {
	return fmt.Sprintf("Invalid form")
}

func (e InvalidFormError) pos() position {
	return e.position
}

type InvalidBindError struct {
	nodes    []*astNode
	position position
}

func (e InvalidBindError) Error() string {
	ss := make([]string, len(e.nodes))
	for i, n := range e.nodes {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Invalid bind error: %s", strings.Join(ss, " "))
}

func (e InvalidBindError) pos() position {
	return e.position
}

type InvalidIfError struct {
	nodes    []*astNode
	position position
}

func (e InvalidIfError) Error() string {
	ss := make([]string, len(e.nodes))
	for i, n := range e.nodes {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Invalid if error: %s", strings.Join(ss, " "))
}

func (e InvalidIfError) pos() position {
	return e.position
}

type InvalidIfConditionError struct {
	cond     Value
	position position
}

func (e InvalidIfConditionError) Error() string {
	return fmt.Sprintf("Invalid if condition: %s", e.cond)
}

func (e InvalidIfConditionError) pos() position {
	return e.position
}

type UnexpectedCharacterError struct {
	char     string
	position position
}

func (e UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("Unexpected character %s", e.char)
}

func (e UnexpectedCharacterError) pos() position {
	return e.position
}

type UnexpectedTokenError struct {
	token    token
	position position
}

func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Unexpected token %s", e.token)
}

func (e UnexpectedTokenError) pos() position {
	return e.position
}

type UnexpectedEndingError struct {
	position position
}

func (e UnexpectedEndingError) Error() string {
	return "Unexpected ending"
}

func (e UnexpectedEndingError) pos() position {
	return e.position
}

type RuntimeError struct {
	reason   string
	position position
}

func (e RuntimeError) Error() string {
	return "Runtime error:" + e.reason
}

func (e RuntimeError) pos() position {
	return e.position
}

type IncorrectNumberOfArgsError struct {
	node     *astNode
	required int
	given    int
}

func (e IncorrectNumberOfArgsError) Error() string {
	return fmt.Sprintf("Incorrect number of args in %s: requires %d but got %d",
		e.node, e.required, e.given)
}

func (e IncorrectNumberOfArgsError) pos() position {
	return e.node.position
}

type MismatchedArgumentsError struct {
	node *astNode
	args []Value
}

func (e MismatchedArgumentsError) Error() string {
	ss := make([]string, len(e.args))
	for i, n := range e.args {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Mismatched arguments to %s: %s", e.node, strings.Join(ss, " "))
}

func (e MismatchedArgumentsError) pos() position {
	return e.node.position
}

type InvalidStreamCallbackError struct {
	reason   string
	position position
}

func (e InvalidStreamCallbackError) Error() string {
	return "Invalid stream callback: " + e.reason
}

func (e InvalidStreamCallbackError) pos() position {
	return e.position
}
