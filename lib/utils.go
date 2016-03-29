package lib

import (
	"bytes"
)

// return string slice from null terminating C strings
func CStrToStrSlice(b []byte) []string {
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

// remove null terminator from a C string
func TrimStrTerm(val []byte) string {
	i := bytes.IndexByte(val, '\x00')
	if i != -1 {
		val = val[:i]
	}

	return string(val)
}
