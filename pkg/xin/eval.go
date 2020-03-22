package xin

import (
	"fmt"
	"strings"
)

// these values are interned
const zeroValue = IntValue(0)
const trueValue = IntValue(1)
const falseValue = zeroValue

type Frame struct {
	Vm     *Vm
	Scope  map[string]Value
	Parent *Frame

	cwd *string
}

func newFrame(parent *Frame) *Frame {
	if parent == nil {
		return &Frame{
			Scope: make(map[string]Value),
		}
	}

	return &Frame{
		Vm:     parent.Vm,
		Scope:  make(map[string]Value),
		Parent: parent,
		cwd:    parent.cwd,
	}
}

func (fr *Frame) String() string {
	ss := make([]string, 0, len(fr.Scope))
	for name, val := range fr.Scope {
		ss = append(ss, name+"\t"+val.String())
	}

	parent := ""
	if fr.Parent != nil {
		parent = fr.Parent.String()
	} else {
		parent = "(<frame root>)"
	}

	return fmt.Sprintf("(<frame dump>)\n  %s\n%s",
		strings.Join(ss, "\n  "),
		parent,
	)
}

func (fr *Frame) Get(name string, pos position) (Value, InterpreterError) {
	if val, prs := fr.Scope[name]; prs {
		return val, nil
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

func (fr *Frame) Up(name string, val Value, pos position) InterpreterError {
	if _, prs := fr.Scope[name]; prs {
		fr.Put(name, val)
		return nil
	} else if fr.Parent != nil {
		return fr.Parent.Up(name, val, pos)
	}
	return UndefinedNameError{
		name:     name,
		position: pos,
	}
}

func eval(fr *Frame, node *astNode) (Value, InterpreterError) {
	if node.isForm {
		return evalForm(fr, node)
	}

	return evalAtom(fr, node)
}

func evalForm(fr *Frame, node *astNode) (Value, InterpreterError) {
	formNode := node.leaves[0]

	switch formNode.token.kind {
	case tkBindForm:
		return evalBindForm(fr, node.leaves[1:])
	case tkIfForm:
		return evalIfForm(fr, node.leaves[1:])
	case tkDoForm:
		return evalDoForm(fr, node.leaves[1:])
	case tkImportForm:
		return evalImportForm(fr, node.leaves[1:])
	case tkName:
		maybeForm, err := fr.Get(formNode.token.value, formNode.position)
		if err != nil {
			return nil, err
		}

		return evalFormValue(fr, node, maybeForm)
	default:
		maybeForm, err := unlazyEval(fr, formNode)
		if err != nil {
			return nil, err
		}

		return evalFormValue(fr, node, maybeForm)
	}
}

func evalFormValue(fr *Frame, node *astNode, maybeForm Value) (Value, InterpreterError) {
	args := make([]Value, len(node.leaves)-1)

	for i, n := range node.leaves[1:] {
		val, err := unlazyEval(fr, n)
		if err != nil {
			return nil, err
		}

		args[i] = val
	}

	return evalFormWithArgs(fr, maybeForm, args, node)
}

func evalFormWithArgs(fr *Frame, maybeForm Value, args []Value, node *astNode) (Value, InterpreterError) {
	switch form := maybeForm.(type) {
	case FormValue:
		localFrame := newFrame(form.frame)

		nargs := len(*form.arguments)
		if len(args) < nargs {
			nargs = len(args)
		}
		for i, v := range args[:nargs] {
			localFrame.Put((*form.arguments)[i], v)
		}

		return LazyValue{
			frame: localFrame,
			node:  form.definition,
		}, nil
	case NativeFormValue:
		return form.evaler(fr, args, node)
	}

	return nil, InvalidFormError{
		position: node.position,
	}
}

// unlazyEvalFormValue is used to evaluate-and-coerce a Xin form invocation,
// usually for use when calling callbacks from builtin native forms.
func unlazyEvalFormValue(fr *Frame, node *astNode, maybeForm Value) (Value, InterpreterError) {
	val, err := evalFormValue(fr, node, maybeForm)
	if err != nil {
		return nil, err
	}

	return unlazy(val)
}

func unlazyEvalFormWithArgs(fr *Frame, maybeForm Value, args []Value, node *astNode) (Value, InterpreterError) {
	val, err := evalFormWithArgs(fr, maybeForm, args, node)
	if err != nil {
		return nil, err
	}

	return unlazy(val)
}

func evalAtom(fr *Frame, node *astNode) (Value, InterpreterError) {
	tok := node.token
	switch tok.kind {
	case tkName:
		return fr.Get(tok.value, node.position)
	case tkNumberLiteralInt, tkNumberLiteralHex:
		return tok.intv, nil
	case tkNumberLiteralDecimal:
		return tok.fracv, nil
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

	specimen, body := args[0], args[1]

	if specimen.isForm {
		if len(specimen.leaves) < 1 {
			return nil, InvalidBindError{nodes: args}
		}

		formNameNode := specimen.leaves[0]
		formName := ""
		if !formNameNode.isForm && formNameNode.token.kind == tkName {
			formName = formNameNode.token.value
		}

		args := specimen.leaves[1:]
		argNames := make(argList, len(args))
		for i, n := range args {
			if !n.isForm && n.token.kind == tkName {
				argNames[i] = n.token.value
			} else {
				return nil, InvalidBindError{nodes: args}
			}
		}

		form := FormValue{
			frame:      fr,
			arguments:  &argNames,
			definition: body,
		}
		fr.Put(formName, form)

		return form, nil
	}

	if specimen.token.kind != tkName {
		return nil, InvalidBindError{nodes: args}
	}

	val, err := unlazyEval(fr, body)
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
	case trueValue:
		return eval(fr, ifTrueNode)
	case falseValue:
		return eval(fr, ifFalseNode)
	default:
		return nil, InvalidIfConditionError{cond: cond}
	}
}

func evalDoForm(fr *Frame, exprs []*astNode) (Value, InterpreterError) {
	if len(exprs) == 0 {
		return zeroValue, nil
	}

	lastIndex := len(exprs) - 1

	for _, node := range exprs[:lastIndex] {
		_, err := unlazyEval(fr, node)
		if err != nil {
			return nil, err
		}
	}

	return eval(fr, exprs[lastIndex])
}
