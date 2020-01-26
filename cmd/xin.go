package main

import (
	"fmt"
	"os"

	"github.com/thesephist/xin/pkg/xin"
)

const cliVersion = "0.1"

func main() {
	// fmt.Printf("Xin v%s\n", cliVersion)

	file, err := os.Open("samples/first.xin")
	defer file.Close()

	vm := xin.NewVm()
	result, err := vm.Eval(file)
	if err != nil {
		fmt.Println("Eval error:", err.Error())
		return
	}

	fmt.Println(result)
}
