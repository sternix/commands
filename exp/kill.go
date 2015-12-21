package main

/*
#include <signal.h>
*/
import "C"

import (
	"fmt"
	//	"syscall"
	"flag"
	"io"
	"os"
	"strings"
)

/*
const char *const sys_signame[NSIG] = {...}
const char *const sys_siglist[NSIG] = {...}
*/
func main() {
	listFlag := flag.Bool("l", false, "")
	flag.Parse()

	if *listFlag {
		printSignals(os.Stdout)
	}

	//	syscall.Kill()
}

func signameToSignum(sig string) int {
	if len(sig) > 3 {
		if strings.ToUpper(sig[:3]) == "SIG" {
			sig = sig[3:]
		}
	}

	sig = strings.ToUpper(sig)

	for i := 1; i < C.NSIG; i++ {
		if strings.Compare(C.GoString(C.sys_signame[i]), sig) == 0 {
			return i
		}
	}

	return -1
}

func printSignals(w io.Writer) {
	for i := 1; i < C.NSIG; i++ {
		fmt.Fprintf(w, "%s", C.GoString(C.sys_signame[i]))
		if (i == C.NSIG/2) || (i == C.NSIG-1) {
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, " ")
		}
	}
}

func signalNotKnown(name string) {
	fmt.Fprintf(os.Stderr, "Unknown signal %s; valid signals:", name)
	printSignals(os.Stderr)
	os.Exit(2)
}

func usage() {

}
