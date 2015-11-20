package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		flags   int = (os.O_WRONLY | os.O_CREATE)
		files   []io.Writer
		exitval int
	)

	appendFlag := flag.Bool("a", false, "Append the output to the files rather than overwriting them.")
	interruptFlag := flag.Bool("i", false, "Ignore the SIGINT signal.")
	flag.Usage = usage

	flag.Parse()

	if *interruptFlag {
		signal.Ignore(syscall.SIGINT)
	}

	if *appendFlag {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	files = append(files, os.Stdout)

	for _, arg := range flag.Args() {
		if arg == "-" {
			continue
		}

		if f, err := os.OpenFile(arg, flags, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "%s - %v", arg, err)
			exitval = 1
		} else {
			defer f.Close()
			files = append(files, f)
		}
	}

	if _, err := io.Copy(io.MultiWriter(files...), os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		exitval = 1
	}

	os.Exit(exitval)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tee [-ai] [file ...]")
	os.Exit(1)
}
