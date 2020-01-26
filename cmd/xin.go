package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thesephist/xin/pkg/xin"
)

const cliVersion = "0.1"

func main() {
	// fmt.Printf("Xin v%s\n", cliVersion)

	file, err := os.Open("samples/hello.xin")
	defer file.Close()

	vm := xin.NewVm()
	result, err := vm.Eval(file)
	if err != nil {
		fmt.Println("Eval error:", err.Error())
		return
	}

	fmt.Println(result)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Repl error:", err.Error())
		}

		result, err := vm.Eval(strings.NewReader(text))
		if err != nil {
			fmt.Println("Eval error:", err.Error())
		} else {
			fmt.Println(result)
		}
	}
}
