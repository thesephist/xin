package xin

import (
	"strings"
)

type astNode struct {
	isForm bool
	token  token
	leaves []*astNode
	position
}

func (n astNode) String() string {
	if n.isForm {
		parts := make([]string, len(n.leaves))
		for i, leaf := range n.leaves {
			parts[i] = leaf.String()
		}
		return "(" + strings.Join(parts, " ") + ")"
	} else {
		return n.token.String()
	}
}

func parse(toks tokenStream) (astNode, InterpreterError) {
	root := astNode{
		isForm: true,
		leaves: []*astNode{},
	}
	// top level of a parse tree is a big do loop
	root.leaves = append(root.leaves, &astNode{
		isForm: false,
		token: token{
			kind: tkDoForm,
		},
	})

	index := 0
	max := len(toks)
	for index < max {
		n, delta, err := parseGeneric(toks[index:])
		if err != nil {
			return astNode{}, err
		}
		root.leaves = append(root.leaves, &n)
		index += delta
	}

	return root, nil
}

func parseGeneric(toks tokenStream) (astNode, int, InterpreterError) {
	if toks[0].kind == tkOpenParen {
		return parseForm(toks)
	} else if toks[0].kind == tkCloseParen {
		return astNode{}, 0, UnexpectedTokenError{
			token:    toks[0],
			position: toks[0].position,
		}
	} else {
		return parseAtom(toks)
	}
}

func parseForm(toks tokenStream) (astNode, int, InterpreterError) {
	root := astNode{
		isForm:   true,
		leaves:   []*astNode{},
		position: toks[0].position,
	}

	index := 1 // Skip open paren
	max := len(toks)

	if max == 1 {
		return astNode{}, 0, UnexpectedEndingError{
			position: toks[0].position,
		}
	}

	if toks[index].kind == tkCloseParen {
		return astNode{}, 0, InvalidFormError{
			position: toks[index].position,
		}
	}

	for toks[index].kind != tkCloseParen {
		n, delta, err := parseGeneric(toks[index:])
		if err != nil {
			return astNode{}, 0, err
		}
		root.leaves = append(root.leaves, &n)
		index += delta

		if index >= max {
			return astNode{}, 0, UnexpectedEndingError{
				position: toks[index-1].position,
			}
		}
	}

	return root, index + 1, nil
}

func parseAtom(toks tokenStream) (astNode, int, InterpreterError) {
	if len(toks) > 0 {
		return astNode{
			isForm:   false,
			token:    toks[0],
			position: toks[0].position,
		}, 1, nil
	} else {
		return astNode{}, 0, UnexpectedEndingError{}
	}
}
