package rt

import (
	"unsafe"
)

type mutex struct {
	_ struct{} //lockRankStruct
	_ uintptr  // key
}

//go:linkname allpptr runtime.allp
var allpptr uintptr

//go:linkname allpLock runtime.allpLock
var allpLock mutex

//go:linkname lock2 runtime.lock2
func lock2(*mutex)

//go:linkname unlock2 runtime.unlock2
func unlock2(*mutex)

func getgptr() unsafe.Pointer

func allp() []unsafe.Pointer {
	pptr := &allpptr
	allp := *(*[]unsafe.Pointer)(unsafe.Pointer(pptr))
	return allp
}
