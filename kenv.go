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

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		if keyval := strings.Split(args[0], "="); len(keyval) == 2 {
			if err := set(keyval[0], keyval[1]); err != nil {
				logErrorAndExit(err)
			}
		} else {
			if *uflag {
				if err := unset(args[0]); err != nil {
					logErrorAndExit(err)
				}
			} else {
				if err := get(args[0]); err != nil {
					logErrorAndExit(err)
				}
			}
		}
	} else {
		if err := dump(); err != nil {
			logErrorAndExit(err)
		}
	}
}

func dump() error {
	kenvs, err := kenv.Dump()
	if err != nil {
		return fmt.Errorf("Unable to dump kenv: %v", err)
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
		return fmt.Errorf("Unable to get %s: %v", name, err)
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
		return fmt.Errorf("Unable to set %s to %s: %v", key, value, err)
	}
	fmt.Printf("%s=\"%s\"\n", key, value)
	return nil
}

func unset(key string) error {
	if err := kenv.Unset(key); err != nil {
		return fmt.Errorf("Unable to unset %s: %v", key, err)
	}
	return nil
}

func logErrorAndExit(err error) {
	if !*qflag {
		log.Println(err)
	}
	os.Exit(-1)
}
