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
	result := vm.Eval(strings.NewReader("(+ 1 14)"))
	fmt.Println(result)
}
