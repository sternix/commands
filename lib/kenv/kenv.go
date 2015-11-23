package kenv

import (
	"syscall"
	"unsafe"
)

const (
	SYS_KENV = 390

	GET   = 0
	SET   = 1
	UNSET = 2
	DUMP  = 3

//	MNAMELEN = 128 /* Maximum name length (for the syscall) */
//	MVALLEN  = 128 /* Maximum value length (for the syscall) */
)

func Get(name string) ([]byte, error) {
	buf := make([]byte, 1024)

	if err := kenv(GET, name, buf); err != nil {
		return nil, err
	} else {
		return buf, nil
	}

}

func Set(name string, value string) error {
	return kenv(SET, name, []byte(value))
}

func Unset(name string) error {
	return kenv(UNSET, name, nil)
}

func Dump() ([]byte, error) {
	_null := unsafe.Pointer(nil)

	if envlen, e := kenv_sys(DUMP, _null, _null, 0); e != nil {
		return nil, e
	} else {
		buf := make([]byte, envlen+1)
		if e2 := kenv(DUMP, "", buf); e2 != nil {
			return nil, e2
		} else {
			return buf, nil
		}
	}
}

func kenv(action int, name string, value []byte) (err error) {

	if action == GET {
		if _name, e := syscall.BytePtrFromString(name); e != nil {
			err = e
			return
		} else {
			if _, e2 := kenv_sys(GET, unsafe.Pointer(_name), unsafe.Pointer(&value[0]), len(value)); e2 != nil {
				err = e2
				return
			} else {
				return nil
			}
		}
	}

	if action == SET {
		if _name, e := syscall.BytePtrFromString(name); e != nil {
			err = e
			return
		} else {
			if _, e2 := kenv_sys(SET, unsafe.Pointer(_name), unsafe.Pointer(&value[0]), len(value)); e2 != nil {
				err = e2
				return
			} else {
				return nil
			}
		}
	}
	if action == UNSET {
		_null := unsafe.Pointer(nil)

		if _name, e := syscall.BytePtrFromString(name); e != nil {
			err = e
			return
		} else {
			if _, e2 := kenv_sys(UNSET, unsafe.Pointer(_name), _null, 0); e2 != nil {
				err = e2
				return
			} else {
				return nil
			}
		}
	}

	if action == DUMP {
		_null := unsafe.Pointer(nil)

		if _, e := kenv_sys(DUMP, _null, unsafe.Pointer(&value[0]), len(value)); e != nil {
			err = e
			return
		} else {
			return nil
		}
	}

	return nil
}

func kenv_sys(what int, name unsafe.Pointer, value unsafe.Pointer, size int) (ret int, err error) {
	r1, _, e := syscall.Syscall6(SYS_KENV, uintptr(what), uintptr(name), uintptr(value), uintptr(size), 0, 0)
	if e != 0 {
		err = syscall.Errno(e)
		ret = -1
	} else {
		err = nil
		ret = int(r1)
	}

	return
}
