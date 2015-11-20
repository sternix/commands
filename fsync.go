package main

import (
	"fmt"
	"github.com/sternix/commands/lib"
	"os"
	"syscall"
)

func main() {
	var rval int = sysexits.OK

	args := os.Args
	if len(args) < 2 {
		usage()
	}

	for _, f := range args[1:] {
		if fd, err := syscall.Open(f, syscall.O_RDONLY, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v\n", f, err)
			if rval == sysexits.OK {
				rval = sysexits.NOINPUT
			}

			continue
		} else {
			if errf := syscall.Fsync(fd); errf != nil {
				fmt.Fprintf(os.Stderr, "fsync %s:%v\n", f, errf)
				if rval == sysexits.OK {
					rval = sysexits.OSERR
				}
			}

			syscall.Close(fd)
		}

	}

	os.Exit(rval)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: fsync file ...\n")
	os.Exit(sysexits.USAGE)
}
