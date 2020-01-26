package xin

import (
	"fmt"
	"strconv"
)

type Frame struct {
	Vm     *Vm
	Scope  map[string]Value
	Parent *Frame
}

func newFrame(parent *Frame) *Frame {
	if parent != nil {
		return &Frame{
			Vm:     parent.Vm,
			Scope:  make(map[string]Value),
			Parent: parent,
		}
	}

	return &Frame{
		Scope: make(map[string]Value),
	}
}

func (fr *Frame) Get(name string, pos position) (Value, InterpreterError) {
	if val, prs := fr.Scope[name]; prs {
		tmp, err := unlazy(val)
		if err != nil {
			return nil, err
		}

		fr.Scope[name] = tmp
		return tmp, nil
	} else if fr.Parent != nil {
		return fr.Parent.Get(name, pos)
	}
	return nil, UndefinedNameError{
		name:     name,
		position: pos,
	}
}

func (fr *Frame) Put(name string, val Value) {
	fr.Scope[name] = val
}

func eval(fr *Frame, node *astNode) (Value, InterpreterError) {
	if node.isForm {
		return evalForm(fr, node)
	}

	return evalAtom(fr, node)
}

func evalForm(fr *Frame, node *astNode) (Value, InterpreterError) {
	if len(node.leaves) == 0 {
		return nil, InvalidFormError{
			node:     node,
			position: node.position,
		}
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

	maybeForm, err := eval(fr, formNode)
	if err != nil {
		return nil, err
	}

	if form, ok := maybeForm.(FormValue); ok {
		localFrame := newFrame(fr)
		for i, n := range node.leaves[1:] {
			localFrame.Put(form.arguments[i], LazyValue{
				frame: fr,
				node:  n,
			})
		}

		return eval(localFrame, form.definition)
	} else if form, ok := maybeForm.(DefaultFormValue); ok {
		args := make([]Value, len(node.leaves)-1)
		for i, n := range node.leaves[1:] {
			args[i] = LazyValue{
				frame: fr,
				node:  n,
			}
		}

		return form.eval(fr, args)
	}

	return nil, InvalidFormError{
		node:     node,
		position: node.position,
	}

}

func evalAtom(fr *Frame, node *astNode) (Value, InterpreterError) {
	tok := node.token
	switch tok.kind {
	case tkName:
		return fr.Get(tok.value, node.position)
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

func evalBindForm(fr *Frame, args []*astNode) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, InvalidBindError{nodes: args}
	}

	specimen := args[0]
	body := args[1]

	if specimen.isForm {
		if len(specimen.leaves) < 1 {
			return nil, InvalidBindError{nodes: args}
		}

		formNameNode := specimen.leaves[0]
		formName := ""
		if !formNameNode.isForm && formNameNode.token.kind == tkName {
			formName = formNameNode.token.value
		}

		argList := specimen.leaves[1:]
		argNames := make([]string, len(argList))
		for i, n := range argList {
			if !n.isForm && n.token.kind == tkName {
				argNames[i] = n.token.value
			} else {
				return nil, InvalidBindError{nodes: args}
			}
		}

		form := FormValue{
			arguments:  argNames,
			definition: body,
		}
		fr.Put(formName, form)

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

func evalIfForm(fr *Frame, args []*astNode) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, InvalidIfError{nodes: args}
	}

	condNode := args[0]
	ifTrueNode := args[1]
	ifFalseNode := args[2]

	cond, err := unlazyEval(fr, condNode)
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

func evalDoForm(fr *Frame, args []*astNode) (Value, InterpreterError) {
	var final Value
	var err InterpreterError

	for _, node := range args {
		final, err = eval(fr, node)
		if err != nil {
			return nil, err
		}
	}

	return final, nil
}
