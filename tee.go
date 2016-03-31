package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	var (
		flags   int = (os.O_WRONLY | os.O_CREATE)
		exitval int
		files   []io.Writer = []io.Writer{os.Stdout}
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

	for _, arg := range flag.Args() {
		if arg == "-" {
			continue
		}

		if f, err := os.OpenFile(arg, flags, os.ModePerm); err != nil {
			log.Printf("%s - %v", arg, err)
			exitval = 1
		} else {
			defer f.Close()
			files = append(files, f)
		}
	}

	if _, err := io.Copy(io.MultiWriter(files...), os.Stdin); err != nil {
		log.Printf("%v", err)
		exitval = 1
	}

	os.Exit(exitval)
}

func usage() {
	log.Fatalln("usage: tee [-ai] [file ...]")
}
