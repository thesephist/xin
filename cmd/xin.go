package main

import (
	"fmt"
	"strings"

	"github.com/thesephist/xin/pkg/xin"
)

const cliVersion = "0.1"

func main() {
	// fmt.Printf("Xin v%s\n", cliVersion)

	testProgram := "(: limit 7) (: (add a b c)\n\t(+ a (+ b limit))) (add 1 limit 3)"

	fmt.Println("Running test program ", testProgram)

	vm := xin.NewVm()
	result, err := vm.Eval(strings.NewReader(testProgram))
	if err != nil {
		fmt.Println("Eval error:", err.Error())
		return
	}

	fmt.Println(result)
}
