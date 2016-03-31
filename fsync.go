package main

import (
	"log"
	"os"
	"syscall"
)

import (
	"github.com/sternix/commands/lib/sysexits"
)


func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	var rval int = sysexits.OK

	args := os.Args
	if len(args) < 2 {
		usage()
	}

	for _, f := range args[1:] {
		if fd, err := syscall.Open(f, syscall.O_RDONLY, 0777); err != nil {
			log.Printf("open %s: %v\n", f, err)
			if rval == sysexits.OK {
				rval = sysexits.NOINPUT
			}
			continue
		} else {
			if errf := syscall.Fsync(fd); errf != nil {
				log.Printf("fsync %s:%v\n", f, errf)
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
	log.Println("usage: fsync file ...")
	os.Exit(sysexits.USAGE)
}
