package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

import (
	"github.com/sternix/commands/lib/kenv"
)

var (
	hflag = flag.Bool("h", false, "Only list hints")
	Nflag = flag.Bool("N", false, "Only display Names")
	qflag = flag.Bool("q", false, "Quiet's the errors")
	uflag = flag.Bool("u", false, "Unset kenv variable")
	vflag = flag.Bool("v", false, "Verbose")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	args := flag.Args()
	if len(args) == 1 {
		if keyval := strings.Split(args[0], "="); len(keyval) == 2 {
			if err := set(keyval[0], keyval[1]); err != nil {
				logErrorAndExit(fmt.Sprintf("Unable to set %s to %s", keyval[0], keyval[1]), err)
			}
		} else {
			if *uflag {
				if err := unset(args[0]); err != nil {
					logErrorAndExit(fmt.Sprintf("Unable to unset %s", args[0]), err)
				}
			} else {
				if err := get(args[0]); err != nil {
					logErrorAndExit(fmt.Sprintf("Unable to get %s", args[0]), err)
				}
			}
		}
	} else {
		if err := dump(); err != nil {
			logErrorAndExit("Unable to dump kenv", err)
		}
	}
}

func dump() error {
	kenvs, err := kenv.Dump()
	if err != nil {
		return err
	}

	for _, item := range kenvs {
		if *hflag {
			if !strings.HasPrefix(item, "hint.") {
				continue
			}
		}

		keyval := strings.Split(item, "=")
		if len(keyval) != 2 {
			continue
		}

		if *Nflag {
			fmt.Printf("%s\n", keyval[0])
		} else {
			fmt.Printf("%s=\"%s\"\n", keyval[0], keyval[1])
		}
	}
	return nil
}

func get(name string) error {
	value, err := kenv.Get(name)
	if err != nil {
		return err
	}

	if *vflag {
		fmt.Printf("%s=\"%s\"", name, value)
	} else {
		fmt.Printf("%s\n", value)
	}
	return nil
}

func set(key string, value string) error {
	if err := kenv.Set(key, value); err != nil {
		return err
	}
	fmt.Printf("%s=\"%s\"\n", key, value)
	return nil
}

func unset(key string) error {
	return kenv.Unset(key)
}

func logErrorAndExit(msg string, err error) {
	if !*qflag {
		log.Printf("%s: %v", msg, err)
	}

	os.Exit(-1)
}
