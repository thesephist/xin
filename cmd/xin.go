package main

import (
	"fmt"
	"strings"

	"github.com/thesephist/xin/pkg/xin"
)

const cliVersion = "0.1"

func main() {
	// fmt.Printf("Xin v%s\n", cliVersion)

	vm := xin.NewVm()
	result, err := vm.Eval(strings.NewReader("(+ 1 14)"))
	if err != nil {
		fmt.Println("Eval error:", err.Error())
		return
	}

	fmt.Println(result)
}
