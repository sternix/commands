package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	var (
		seconds time.Duration
		orjSec  uint
	)

	args := os.Args

	if len(args) != 2 {
		usage()
	}

	if sec, err := strconv.ParseUint(args[1], 10, 32); err != nil {
		usage()
	} else {
		orjSec = uint(sec)
		str := fmt.Sprintf("%ds", orjSec)
		seconds, _ = time.ParseDuration(str)
	}

	sigInfo := make(chan os.Signal, 1)
	signal.Notify(sigInfo, syscall.SIGINFO)

	start := time.Now()

	go func() {
		for {
			<-sigInfo
			fmt.Printf("sleep: about %d second(s) left out of the original %d\n", orjSec-uint((time.Since(start).Seconds())), orjSec)
		}
	}()

	time.Sleep(seconds)
}

func usage() {
	log.Fatalln("usage: sleep seconds")
}
