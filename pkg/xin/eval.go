package xin

import (
	"fmt"
	"strconv"
	"strings"
)

type Frame struct {
	Scope  map[string]Value
	Parent *Frame
}

func newFrame() *Frame {
	return &Frame{
		Scope: make(map[string]Value),
	}
}

type UndefinedNameError struct {
	name string
}

func (e UndefinedNameError) Error() string {
	return fmt.Sprintf("Undefined name %s", e.name)
}

type InvalidFormError struct {
	node *astNode
}

func (e InvalidFormError) Error() string {
	return fmt.Sprintf("Invalid form: %s", e.node)
}

func (fr *Frame) Get(name string) (Value, error) {
	if val, prs := fr.Scope[name]; prs {
		return val, nil
	}
	if fr.Parent != nil {
		return fr.Parent.Get(name)
	}
	return nil, UndefinedNameError{name: name}
}

func (fr *Frame) Put(name string, val Value) {
	fr.Scope[name] = val
}

func eval(fr *Frame, node *astNode) (Value, error) {
	if node.isForm {
		return evalForm(fr, node)
	}

	return evalAtom(fr, node)
}

func evalForm(fr *Frame, node *astNode) (Value, error) {
	if len(node.leaves) == 0 {
		return nil, InvalidFormError{node: node}
	}

	formNode := node.leaves[0]

	// Evaluate special forms
	if formNode.token.kind == tkName {
		switch formNode.token.value {
		case ":":
			return evalBindForm(fr, node.leaves[1:])
		case "if":
			return evalIfForm(fr, node.leaves[1:])
		case "do":
			return evalDoForm(fr, node.leaves[1:])
		}
	}

	formHead, err := eval(fr, formNode)
	if err != nil {
		return nil, err
	}

	form, ok := formHead.(FormValue)
	if ok {
		localFrame := newFrame()
		for i, n := range node.leaves[1:] {
			localFrame.Scope[form.arguments[i]] = LazyValue{node: n}
		}

		return eval(localFrame, form.definition)
	}

	dForm, ok := formHead.(DefaultFormValue)
	if ok {
		args := make([]Value, len(node.leaves)-1)
		for i, n := range node.leaves[1:] {
			args[i] = LazyValue{node: n}
		}

		return dForm.eval(fr, args)
	}

	return nil, InvalidFormError{node: node}

}

func evalAtom(fr *Frame, node *astNode) (Value, error) {
	tok := node.token
	switch tok.kind {
	case tkName:
		return fr.Get(tok.value)
	case tkNumberLiteralInt:
		iv, _ := strconv.ParseInt(tok.value, 10, 64)
		return IntValue(iv), nil
	case tkNumberLiteralDecimal:
		fv, _ := strconv.ParseFloat(tok.value, 64)
		return FracValue(fv), nil
	case tkNumberLiteralHex:
		iv, _ := strconv.ParseInt(tok.value, 16, 64)
		return IntValue(iv), nil
	case tkStringLiteral:
		return StringValue(tok.value), nil
	default:
		panic(fmt.Sprintf("Unrecognized token type: %d", node.token.kind))
	}
}

type InvalidBindError struct {
	nodes []*astNode
}

func (e InvalidBindError) Error() string {
	ss := make([]string, len(e.nodes))
	for i, n := range e.nodes {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Invalid bind error: %s", strings.Join(ss, " "))
}

type InvalidIfError struct {
	nodes []*astNode
}

func (e InvalidIfError) Error() string {
	ss := make([]string, len(e.nodes))
	for i, n := range e.nodes {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Invalid if error: %s", strings.Join(ss, " "))
}

type InvalidIfConditionError struct {
	cond Value
}

func (e InvalidIfConditionError) Error() string {
	return fmt.Sprintf("Invalid if condition: %s", e.cond)
}

func evalBindForm(fr *Frame, args []*astNode) (Value, error) {
	if len(args) != 2 {
		return nil, InvalidBindError{nodes: args}
	}

	specimen := args[0]
	body := args[1]

	if specimen.isForm {
		if len(specimen.leaves) < 1 {
			return nil, InvalidBindError{nodes: args}
		}

		argList := specimen.leaves[1:]
		argNames := make([]string, len(argList))
		for i, n := range argList {
			if !n.isForm && n.token.kind == tkName {
				argNames[i] = n.token.value
			}

			return nil, InvalidBindError{nodes: args}
		}

		form := FormValue{
			arguments:  argNames,
			definition: body,
		}
		fr.Put("fn", form)

		return form, nil
	}

	if specimen.token.kind != tkName {
		return nil, InvalidBindError{nodes: args}
	}

	val, err := eval(fr, body)
	if err != nil {
		return nil, err
	}

	fr.Put(specimen.token.value, val)
	return val, nil
}

func evalIfForm(fr *Frame, args []*astNode) (Value, error) {
	if len(args) != 3 {
		return nil, InvalidIfError{nodes: args}
	}

	condNode := args[0]
	ifTrueNode := args[1]
	ifFalseNode := args[2]

	cond, err := eval(fr, condNode)
	if err != nil {
		return nil, err
	}

	switch cond {
	case IntValue(1):
		return eval(fr, ifTrueNode)
	case IntValue(0):
		return eval(fr, ifFalseNode)
	default:
		return nil, InvalidIfConditionError{cond: cond}
	}
}

func evalDoForm(fr *Frame, args []*astNode) (Value, error) {
	var final Value
	var err error

	for _, node := range args {
		final, err = eval(fr, node)
		if err != nil {
			return nil, err
		}
	}

	return final, nil
}
