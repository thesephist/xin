package xin

import (
	"fmt"
	"strings"
)

func FormatError(e InterpreterError) string {
	return e.Error() + "\n\t" + e.Position().String()
}

type InterpreterError interface {
	error
	Position() position
}

type UndefinedNameError struct {
	name     string
	position position
}

func (e UndefinedNameError) Error() string {
	return fmt.Sprintf("Undefined name %s", e.name)
}

func (e UndefinedNameError) Position() position {
	return e.position
}

type InvalidFormError struct {
	node     *astNode
	position position
}

func (e InvalidFormError) Error() string {
	return fmt.Sprintf("Invalid form: %s", e.node)
}

func (e InvalidFormError) Position() position {
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

func (e InvalidBindError) Position() position {
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

func (e InvalidIfError) Position() position {
	return e.position
}

type InvalidIfConditionError struct {
	cond     Value
	position position
}

func (e InvalidIfConditionError) Error() string {
	return fmt.Sprintf("Invalid if condition: %s", e.cond)
}

func (e InvalidIfConditionError) Position() position {
	return e.position
}

type UnexpectedCharacterError struct {
	char     string
	position position
}

func (e UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("Unexpected character %s", e.char)
}

func (e UnexpectedCharacterError) Position() position {
	return e.position
}

type UnexpectedTokenError struct {
	token    token
	position position
}

func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Unexpected token %s", e.token)
}

func (e UnexpectedTokenError) Position() position {
	return e.position
}

type UnexpectedEndingError struct {
	position position
}

func (e UnexpectedEndingError) Error() string {
	return "Unexpected ending"
}

func (e UnexpectedEndingError) Position() position {
	return e.position
}

type RuntimeError struct {
	reason   string
	position position
}

func (e RuntimeError) Error() string {
	return "Runtime error:" + e.reason
}

func (e RuntimeError) Position() position {
	return e.position
}

type IncorrectNumberOfArgsError struct {
	required int
	given    int
	position position
}

func (e IncorrectNumberOfArgsError) Error() string {
	return fmt.Sprintf("Incorrect number of arguments: requires %d but got %d",
		e.required, e.given)
}

func (e IncorrectNumberOfArgsError) Position() position {
	return e.position
}

type MismatchedArgumentsError struct {
	args     []Value
	position position
}

func (e MismatchedArgumentsError) Error() string {
	ss := make([]string, len(e.args))
	for i, n := range e.args {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Mismatched arguments: %s", strings.Join(ss, " "))
}

func (e MismatchedArgumentsError) Position() position {
	return e.position
}

type InvalidStreamCallbackError struct {
	reason   string
	position position
}

func (e InvalidStreamCallbackError) Error() string {
	return "Invalid stream callback:" + e.reason
}

func (e InvalidStreamCallbackError) Position() position {
	return e.position
}
