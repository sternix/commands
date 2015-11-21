package kenv

import (
	"syscall"
	"unsafe"
	"fmt"
)

/*
390     AUE_NULL        STD     { int kenv(int what, const char *name, char *value, int len); }
*/

const (
	SYS_KENV   = 390
	GET        = 0
	SET        = 1
	UNSET      = 2
	DUMP       = 3

	MNAMELEN   = 128     /* Maximum name length (for the syscall) */
	MVALLEN    = 128     /* Maximum value length (for the syscall) */
)

func Get(name string) (string , error) {
	kenv(GET,name,"")
	return "",nil
}

func Set(name string , value string) error {
	kenv(SET,name,"")
	return nil
}

func Unset(name string) error {
	kenv(UNSET,name,"")
	return nil
}

func Dump() string {
	kenv(DUMP,"","")
	return ""
}

func kenv(action int , name string, value string) (err error) {
	/*
	var (
		_p0 *byte
		_p1 *byte
		e error
	)
	*/

/*
	if action == GET {
	
	}

	if action == SET {
	
	}

	if action == UNSET {
	
	}

*/
	if action == DUMP {
		_null := unsafe.Pointer(uintptr(0))
		if envlen , err := kenv_sys(DUMP, _null, _null, 0); err != nil {
			fmt.Printf("%v\n",err)
		} else {
			buf := make([]byte , envlen + 1)
//			fmt.Printf("len : %d : \n",len(buf))
			kenv_sys(DUMP,_null,unsafe.Pointer(&buf[0]) , envlen + 1)
			fmt.Printf("%s\n",string(buf))
		}
	}

	/*
	_p0 , e = syscall.BytePtrFromString(name)
	if e != nil {
		err = e
		return
	}

	_p1 = syscall.BytePtrFromString(value)
	if e != nil {
		err = e
		return
	}
	*/

	return nil
}

/*
// func Syscall(trap int64, a1, a2, a3 int64) (r1, r2, err int64);
// func Syscall6(trap int64, a1, a2, a3, a4, a5, a6 int64) (r1, r2, err int64);
// func Syscall9(trap int64, a1, a2, a3, a4, a5, a6, a7, a8, a9 int64) (r1, r2, err int64)
*/

func kenv_sys(what int ,name unsafe.Pointer, value unsafe.Pointer, size int) (ret int , err error) {
	r1 , _ , e := syscall.Syscall6(SYS_KENV,uintptr(what),uintptr(name),uintptr(value),uintptr(size),0,0)
	if e != 0 {
		err = syscall.Errno(e)
		ret = -1
	} else {
		err = nil
		ret = int(r1)
	}

	return
}
