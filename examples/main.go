package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coreyog/pipeline"
)

func main() {
	pipe := pipeline.New()

	err := pipe.PushFunc(os.Open)
	if err != nil {
		panic(err)
	}

	err = pipe.PushFunc(ioutil.ReadAll)
	if err != nil {
		panic(err)
	}

	err = pipe.PushFunc(toString)
	if err != nil {
		panic(err)
	}

	results, err := pipe.Call("example.txt")
	if err != nil {
		panic(err)
	}

	for i, r := range results {
		fmt.Printf("%d) %v\n", i, r)
	}

	fmt.Println("DONE")
}

func toString(data []byte) (str string) {
	return string(data)
}

func x() {}
