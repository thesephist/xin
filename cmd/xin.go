package main

import (
	"fmt"
	"strings"

	"github.com/thesephist/xin/pkg/xin"
)

const cliVersion = "0.1"

func main() {
	// fmt.Printf("Xin v%s\n", cliVersion)

	toks, err := xin.Parse(xin.Lex(strings.NewReader("(if 1 (sum 10 2) (* 4 5))")))
	if err != nil {
		fmt.Printf("there was an error: %s", err.Error())
	}
	fmt.Println(toks.String())

	// vm := xin.NewVm()
	// result := vm.Eval(strings.NewReader("(+ 1 14)"))
	// fmt.Println(result)
}
