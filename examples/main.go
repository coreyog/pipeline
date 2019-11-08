package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coreyog/pipeline"
)

func main() {
	BasicChain()

	fmt.Println()

	VariadicChain()

	fmt.Println("DONE")
}

func BasicChain() {
	pipe := pipeline.New()
	fmt.Println("Basic chain with error checking:")

	err := pipe.PushFunc(os.Open, ioutil.ReadAll, toString)
	if err != nil {
		panic(err)
	}

	fmt.Println(pipe.String())
	fmt.Println()

	args := []interface{}{"example.txt"}

	fmt.Println("Input: ")
	for i, arg := range args {
		fmt.Printf("%d) %v\n", i, arg)
	}

	fmt.Println()

	results, err := pipe.Call(args...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Output: ")

	for i, r := range results {
		fmt.Printf("%d) %v\n", i, r)
	}
}

func VariadicChain() {
	pipe := pipeline.New()

	fmt.Println("Variadic chain:")

	err := pipe.PushFunc(abc, fmt.Println)
	if err != nil {
		panic(err)
	}

	fmt.Println(pipe.String())
	fmt.Println()

	fmt.Println("Input:")
	fmt.Println("None")
	fmt.Println()

	results, err := pipe.Call()
	if err != nil {
		panic(err)
	}

	fmt.Println("Output: ")

	for i, r := range results {
		fmt.Printf("%d) %v\n", i, r)
	}
}

func toString(data []byte) (str string) {
	return string(data)
}

func abc() (a, b, c string) {
	return "A", "B", "C"
}
