package kenv_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

import (
	"github.com/sternix/commands/lib/kenv"
)

const kenvCmd = "/bin/kenv"

func NullTermToStrings(b []byte) []string {
	var s []string

	for {
		i := bytes.IndexByte(b, '\x00')
		if i == -1 {
			break
		}
		s = append(s, string(b[0:i]))
		b = b[i+1:]
	}
	return s
}

func TrimStrTerm(val []byte) string {
	i := bytes.IndexByte(val, '\x00')
	if i != -1 {
		val = val[:i]
	}

	return string(val)
}

func TestKenvGet(t *testing.T) {
	// exec os command
	out, err := exec.Command(kenvCmd).Output()
	if err != nil {
		t.Error(err)
	}

	// get key and values from os commands output
	kenvs := make(map[string]string)
	for _, ke := range strings.Split(string(out), "\n") {
		if ke == "" {
			continue
		}
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

		if TrimStrTerm(get) != v {
			t.Errorf("k:%q = get:%q , expected v:%q", k, get, v)
		}
	}
}

func TestKenvDump(t *testing.T) {
	// exec os command
	out, err := exec.Command(kenvCmd).Output()
	if err != nil {
		t.Error(err)
	}

	fromCmd := strings.Split(string(out), "\n")
	fromCmd = fromCmd[:len(fromCmd)-1] //remove last newline

	buf, err := kenv.Dump()
	if err != nil {
		t.Error(err)
	}

	var fromLib []string
	for _, item := range NullTermToStrings(buf) {
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
