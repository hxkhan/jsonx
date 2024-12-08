package main

import (
	"fmt"
	"os"

	"github.com/hk-32/jsonx"
)

func main() {
	file, err := os.ReadFile("./input.json")
	if err != nil {
		panic(err)
	}

	structure, err := jsonx.Decode(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(structure)
}
