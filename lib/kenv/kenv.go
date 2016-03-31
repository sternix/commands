package kenv

import (
	"syscall"
	"unsafe"
)

import (
	"github.com/sternix/commands/lib"
)

const (
	SYS_KENV = 390

	GET   = 0
	SET   = 1
	UNSET = 2
	DUMP  = 3
)

func Get(name string) (string, error) {
	var (
		buf  = make([]byte, 1024)
		bptr *byte
		err  error
		size int
	)

	if bptr, err = syscall.BytePtrFromString(name); err != nil {
		return "", err
	}

	if size, err = kenv_sys(GET, unsafe.Pointer(bptr), unsafe.Pointer(&buf[0]), len(buf)); err != nil {
		return "", err
	}

	buf = buf[:size]
	//remove C null term
	if buf[len(buf)-1] == '\x00' {
		buf = buf[:len(buf)-1]
	}

	return string(buf), nil
}

func Set(name string, value string) error {
	var (
		bptrname  *byte
		bptrvalue *byte
		err       error
	)

	if bptrname, err = syscall.BytePtrFromString(name); err != nil {
		return err
	}

	if bptrvalue, err = syscall.BytePtrFromString(value); err != nil {
		return err
	}

	if _, err = kenv_sys(SET, unsafe.Pointer(bptrname), unsafe.Pointer(bptrvalue), len(value)+1); err != nil {
		return err
	}

	return nil
}

func Unset(name string) error {
	var (
		bptr *byte
		err  error
	)

	if bptr, err = syscall.BytePtrFromString(name); err != nil {
		return err
	}

	if _, err = kenv_sys(UNSET, unsafe.Pointer(bptr), unsafe.Pointer(nil), 0); err != nil {
		return err
	}

	return nil
}

func Dump() ([]string, error) {
	var (
		envlen int
		err    error
	)

	// get size of the buffer
	if envlen, err = kenv_sys(DUMP, unsafe.Pointer(nil), unsafe.Pointer(nil), 0); err != nil {
		return nil, err
	}

	// get envs to buffer
	buf := make([]byte, envlen+1)
	if _, err = kenv_sys(DUMP, unsafe.Pointer(nil), unsafe.Pointer(&buf[0]), len(buf)); err != nil {
		return nil, err
	}

	// remove last null term
	if buf[len(buf)-1] == '\x00' {
		buf = buf[:len(buf)-1]
	}

	// remove all c null term character (\x00) from string
	return lib.CStrToStrSlice(buf), nil
}

func kenv_sys(what int, name unsafe.Pointer, value unsafe.Pointer, size int) (int, error) {
	ret, _, err := syscall.Syscall6(SYS_KENV, uintptr(what), uintptr(name), uintptr(value), uintptr(size), 0, 0)
	if err != 0 {
		return -1, syscall.Errno(err)
	}
	return int(ret), nil
}
