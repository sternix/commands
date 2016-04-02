package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

// getting from /usr/include/syslog.h
const (
	LOG_KERN    = 0 << 3
	LOG_FACMASK = 0x03f8
)

var (
	cflag = flag.Bool("c", false, "Clear kernel message buffer")
)

func isKernSyslogEntry(fac int64) bool {
	return (fac & LOG_FACMASK) == LOG_KERN
}

func main() {
	flag.Parse()

	msgbuf, err := syscall.Sysctl("kern.msgbuf")
	if err != nil {
		log.Fatalf("sysctl kern.msgbuf: %v", err)
	}

	buf := []byte(msgbuf)
	//trim leading \0 chars ( \0 == \x00 )
	for i := 0; i < len(buf); i++ {
		if buf[i] != '\x00' {
			buf = buf[i:]
			break
		}
	}

	if buf[len(buf)-1] == '\n' {
		buf = buf[0 : len(buf)-1]
	}

	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`^<(\d+)>`) // search for syslog entries
		if matches := re.FindStringSubmatch(line); matches != nil {
			fac, err := strconv.ParseInt(matches[1], 10, 32)
			if err == nil {
				if !isKernSyslogEntry(fac) {
					continue
				}
				line = re.ReplaceAllString(line, "") // remove syslog facility
			}
		}
		fmt.Println(line)
	}

	/*
		if *cflag {
			TODO:implement syscall that can be set a var....
			kern.msgbuf_clear
		}
	*/
}
