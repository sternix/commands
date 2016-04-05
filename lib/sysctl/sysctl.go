// +build freebsd

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the
// https://github.com/golang/go/blob/master/LICENSE file.
package sysctl

import (
	"C"
	"syscall"
	"unsafe"
)

// some functions copied from
// https://github.com/golang/sys/blob/master/unix/syscall_bsd.go

type (
	_C_int C.int
)

var (
	ZeroPtr uintptr
	EIO     = syscall.Errno(0x5)
)

func NameToMIB(name string) (mib []_C_int, err error) {
	const siz = unsafe.Sizeof(mib[0])

	var buf [syscall.CTL_MAXNAME + 2]_C_int
	n := uintptr(syscall.CTL_MAXNAME) * siz

	p := (*byte)(unsafe.Pointer(&buf[0]))
	bytes, err := syscall.ByteSliceFromString(name)
	if err != nil {
		return nil, err
	}

	if err = sysctl([]_C_int{0, 3}, p, &n, &bytes[0], uintptr(len(name))); err != nil {
		return nil, err
	}
	return buf[0 : n/siz], nil
}

func sysctl(mib []_C_int, oldb *byte, oldblen *uintptr, newb *byte, newblen uintptr) error {
	var (
		mibptr unsafe.Pointer
		err    error
	)

	if len(mib) > 0 {
		mibptr = unsafe.Pointer(&mib[0])
	} else {
		mibptr = unsafe.Pointer(&ZeroPtr)
	}
	_, _, e := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(mibptr), uintptr(len(mib)), uintptr(unsafe.Pointer(oldb)), uintptr(unsafe.Pointer(oldblen)), uintptr(unsafe.Pointer(newb)), uintptr(newblen))
	if e != 0 {
		err = syscall.Errno(e)
	}
	return err
}

func sysctlmib(name string, args ...int) ([]_C_int, error) {
	// Translate name to mib number.
	mib, err := NameToMIB(name)
	if err != nil {
		return nil, err
	}

	for _, a := range args {
		mib = append(mib, _C_int(a))
	}

	return mib, nil
}

func ByName(name string) (string, error) {
	return Args(name)
}

func Args(name string, args ...int) (string, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return "", err
	}

	// Find size.
	n := uintptr(0)
	if err := sysctl(mib, nil, &n, nil, 0); err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}

	// Read into buffer of that size.
	buf := make([]byte, n)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return "", err
	}

	// Throw away terminating NUL.
	if n > 0 && buf[n-1] == '\x00' {
		n--
	}
	return string(buf[0:n]), nil
}

func Uint32(name string) (uint32, error) {
	return Uint32Args(name)
}

func Uint32Args(name string, args ...int) (uint32, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return 0, err
	}

	n := uintptr(4)
	buf := make([]byte, 4)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, EIO
	}
	return *(*uint32)(unsafe.Pointer(&buf[0])), nil
}

func Uint64(name string, args ...int) (uint64, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return 0, err
	}

	n := uintptr(8)
	buf := make([]byte, 8)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, EIO
	}
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}

func Raw(name string, args ...int) ([]byte, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return nil, err
	}

	// Find size.
	n := uintptr(0)
	if err := sysctl(mib, nil, &n, nil, 0); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	// Read into buffer of that size.
	buf := make([]byte, n)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return nil, err
	}

	// The actual call may return less than the original reported required
	// size so ensure we deal with that.
	return buf[:n], nil
}

func SetString(name string, value string) error {
	mib, err := NameToMIB(name)
	if err != nil {
		return err
	}

	valueslc, err := syscall.ByteSliceFromString(value)
	if err != nil {
		return err
	}

	if err := sysctl(mib, nil, nil, &valueslc[0], uintptr(len(value))); err != nil {
		return err
	}

	return nil

}

func SetUint32(name string, value uint32) error {
	mib, err := NameToMIB(name)
	if err != nil {
		return err
	}

	if err := sysctl(mib, nil, nil, (*byte)(unsafe.Pointer(&value)), 4); err != nil {
		return err
	}

	return nil
}

func SetUint64(name string, value uint64) error {
	mib, err := NameToMIB(name)
	if err != nil {
		return err
	}

	if err := sysctl(mib, nil, nil, (*byte)(unsafe.Pointer(&value)), 8); err != nil {
		return err
	}

	return nil
}
