package main

import (
	"bytes"
	"flag"
	"os"
)

func main() {
	omitNewLine := flag.Bool("n", false, "Do not print the trailing newline character.")

	flag.Parse()

	var buffer bytes.Buffer

	var sep string

	for i := 0; i < flag.NArg(); i++ {
		buffer.WriteString(sep)
		buffer.WriteString(flag.Arg(i))
		sep = " "
	}

	if !*omitNewLine {
		buffer.WriteString("\n")
	}

	buffer.WriteTo(os.Stdout)
}
