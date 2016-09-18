package main

/*
#include <unistd.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

import (
	"github.com/sternix/commands/lib/sysctl"
)

type Timeval struct {
	Sec  int64
	Usec int64
}

func main() {
	var clk_tck int64 = clk_tck_from_sysconf()
	fmt.Printf("%d\n", clk_tck)

	bt := boottime()
	fmt.Printf("%d - %d\n", bt.Sec, bt.Usec)

}

func clk_tck_from_sysconf() int64 {
	return int64(C.sysconf(C._SC_CLK_TCK))
}

func boottime() Timeval {
	// boottime
	var boottv Timeval

	tv, err := sysctl.Raw("kern.boottime")
	if err != nil {
		log.Fatal(err)
	}

	br := bytes.NewReader(tv)
	binary.Read(br, binary.LittleEndian, &boottv)
	return boottv
}
