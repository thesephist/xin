package xin

import (
	"fmt"
	"strconv"
	"strings"
)

const zeroValue = IntValue(0)

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
		if lzv, ok := val.(LazyValue); ok {
			tmp, err := unlazy(lzv)
			if err != nil {
				return nil, err
			}

			fr.Scope[name] = tmp
			return tmp, nil
		}

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

func eval(fr *Frame, node *astNode) (Value, InterpreterError) {
	if node.isForm {
		return evalForm(fr, node)
	}

	return evalAtom(fr, node)
}

func evalForm(fr *Frame, node *astNode) (Value, InterpreterError) {
	formNode := node.leaves[0]

	var maybeForm Value
	var err InterpreterError
	if formNode.token.kind == tkName {
		// Evaluate special forms
		switch formNode.token.value {
		case ":":
			return evalBindForm(fr, node.leaves[1:])
		case "if":
			return evalIfForm(fr, node.leaves[1:])
		case "do":
			return evalDoForm(fr, node.leaves[1:])
		case "import":
			return evalImportForm(fr, node.leaves[1:])
		}

		// In the common case that the form head is a named reference,
		// take this shortcut (Get) in this hot path.
		maybeForm, err = fr.Get(formNode.token.value, formNode.position)
	} else {
		maybeForm, err = unlazyEval(fr, formNode)
	}
	if err != nil {
		return nil, err
	}

	switch form := maybeForm.(type) {
	case FormValue:
		localFrame := newFrame(form.frame)

		nargs := len(*form.arguments)
		if len(node.leaves)-1 < nargs {
			nargs = len(node.leaves) - 1
		}
		for i, n := range node.leaves[1 : nargs+1] {
			if n.isLiteral() {
				val, err := evalAtom(fr, n)
				if err != nil {
					return nil, err
				}

				localFrame.Put((*form.arguments)[i], val)
			} else {
				localFrame.Put((*form.arguments)[i], LazyValue{
					frame: fr,
					node:  n,
				})
			}
		}

		return LazyValue{
			frame: localFrame,
			node:  form.definition,
		}, nil
	case NativeFormValue:
		args := make([]Value, len(node.leaves)-1)
		for i, n := range node.leaves[1:] {
			if n.isLiteral() {
				val, err := evalAtom(fr, n)
				if err != nil {
					return nil, err
				}

				args[i] = val
			} else {
				args[i] = LazyValue{
					frame: fr,
					node:  n,
				}
			}
		}

		return form.eval(fr, args, node)
	}

	return nil, InvalidFormError{
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
	case IntValue(1):
		return eval(fr, ifTrueNode)
	case zeroValue:
		return eval(fr, ifFalseNode)
	default:
		return nil, InvalidIfConditionError{cond: cond}
	}
}

func evalDoForm(fr *Frame, args []*astNode) (Value, InterpreterError) {
	if len(args) == 0 {
		return zeroValue, nil
	}

	lastIndex := len(args) - 1

	for _, node := range args[:lastIndex] {
		_, err := unlazyEval(fr, node)
		if err != nil {
			return nil, err
		}
	}

	return eval(fr, args[lastIndex])
}
