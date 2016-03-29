package kenv_test

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

import (
	"github.com/sternix/commands/lib/kenv"
)

func kenvsFromCmd(t *testing.T) []string {
	out, err := exec.Command("/bin/kenv").Output()
	if err != nil {
		t.Error(err)
	}

	kenvs := strings.Split(string(out), "\n")
	kenvs = kenvs[:len(kenvs)-1] //remove last newline
	return kenvs
}

func TestKenvGet(t *testing.T) {
	// get key and values from os commands output
	kenvs := make(map[string]string)
	for _, ke := range kenvsFromCmd(t) {
		// items in the form key="value"
		strs := strings.Split(ke, "=")
		if len(strs) != 2 {
			t.Errorf("%s", ke)
		}
		val, err := strconv.Unquote(strs[1])
		if err != nil {
			t.Error(err)
		}

		kenvs[strs[0]] = val
	}

	// for each key call kenv.Get and check value
	for k, v := range kenvs {
		get, err := kenv.Get(k)
		if err != nil {
			t.Error(err)
		}

		if get != v {
			t.Errorf("k:%q = get:%q , expected v:%q", k, get, v)
		}
	}
}

func TestKenvDump(t *testing.T) {
	fromCmd := kenvsFromCmd(t)

	kenvs, err := kenv.Dump()
	if err != nil {
		t.Error(err)
	}

	var fromLib []string
	for _, item := range kenvs {
		keyval := strings.Split(item, "=")
		if len(keyval) != 2 {
			t.Errorf("%s has different format than key=val", keyval)
		}
		fromLib = append(fromLib, fmt.Sprintf("%s=\"%s\"", keyval[0], keyval[1]))
	}

	if len(fromLib) != len(fromCmd) {
		t.Errorf("cmd has %d item but lib has %d", len(fromCmd), len(fromLib))
	}

	for i, ke := range fromCmd {
		if fromLib[i] != ke {
			t.Errorf("lib %s different than %s", fromLib[i], ke)
		}
	}
}

/*
func TestKenvSet(t *testing.T) {

}

func TestKenvUnset(t *testing.T) {

}

*/
