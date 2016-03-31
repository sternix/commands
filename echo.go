package main

import (
	"flag"
	"fmt"
)

func main() {
	omitNewLine := flag.Bool("n", false, "Do not print the trailing newline character.")
	flag.Parse()

	var sep string
	for i := 1; i < flag.NArg(); i++ {
		fmt.Print(sep)
		fmt.Print(flag.Arg(i))
		sep = " "
	}

	if !*omitNewLine {
		fmt.Println()
	}
}
