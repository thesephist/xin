package xin

import (
	"fmt"
	"strings"
)

type astNode struct {
	isForm bool
	token  token
	leaves []*astNode
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

type UnexpectedTokenError struct {
	token token
}

func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Unexpected token %s", e.token)
}

type UnexpectedEndingError struct{}

func (e UnexpectedEndingError) Error() string {
	return "Unexpected ending"
}

func parse(toks tokenStream) (astNode, error) {
	root := astNode{
		isForm: true,
		leaves: []*astNode{},
	}
	// top level of a parse tree is a big do loop
	root.leaves = append(root.leaves, &astNode{
		isForm: false,
		token: token{
			kind:  tkName,
			value: "do",
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

func parseGeneric(toks tokenStream) (astNode, int, error) {
	if toks[0].kind == tkOpenParen {
		return parseForm(toks)
	} else if toks[0].kind == tkCloseParen {
		return astNode{}, 0, UnexpectedTokenError{token: toks[0]}
	} else {
		return parseAtom(toks)
	}
}

func parseForm(toks tokenStream) (astNode, int, error) {
	root := astNode{
		isForm: true,
		leaves: []*astNode{},
	}

	index := 1 // Skip open paren
	max := len(toks)

	if max == 1 {
		return astNode{}, 0, UnexpectedEndingError{}
	}

	for toks[index].kind != tkCloseParen {
		n, delta, err := parseGeneric(toks[index:])
		if err != nil {
			return astNode{}, 0, err
		}
		root.leaves = append(root.leaves, &n)
		index += delta

		if index >= max {
			return astNode{}, 0, UnexpectedEndingError{}
		}
	}

	return root, index + 1, nil
}

func parseAtom(toks tokenStream) (astNode, int, error) {
	if len(toks) > 0 {
		return astNode{
			isForm: false,
			token:  toks[0],
		}, 1, nil
	} else {
		return astNode{}, 0, UnexpectedEndingError{}
	}
}
