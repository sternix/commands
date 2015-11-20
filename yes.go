package main

import (
	"fmt"
	"os"
)

func main() {
	var str string = "y"

	args := os.Args

	if len(args) > 1 {
		str = args[1]
	}

	for {
		fmt.Println(str)
	}
}
