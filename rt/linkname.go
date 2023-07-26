package rt

import (
	"unsafe"
)

//go:linkname allpptr runtime.allp
var allpptr uintptr

//go:linkname allpLock runtime.allpLock
var allpLock uintptr

//go:linkname lock2 runtime.lock2
func lock2(*uintptr)

//go:linkname unlock2 runtime.unlock2
func unlock2(*uintptr)

func lockallp() {
	lock2(&allpLock)
}

func unlockallp() {
	unlock2(&allpLock)
}

func getgptr() unsafe.Pointer

func allp() []unsafe.Pointer {
	pptr := &allpptr
	allp := *(*[]unsafe.Pointer)(unsafe.Pointer(pptr))
	return allp
}
