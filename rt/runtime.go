package rt

import (
	"unsafe"
)

// NumTimers return timer count in the process
func NumTimers() (int, error) {
	offset, err := Offsetof("runtime.p", "numTimers")
	if err != nil {
		return 0, err
	}

	lockallp()
	defer unlockallp()

	count := 0
	for _, p := range allp() {
		count += int(*(*int32)(unsafe.Pointer(uintptr(p) + offset)))
	}
	return count, nil
}

// Goid returns goid of current goroutine
func Goid() (uint64, error) {
	offset, err := Offsetof("runtime.g", "goid")
	if err != nil {
		return 0, err
	}

	goid := *(*uint64)(unsafe.Pointer(uintptr(getgptr()) + offset))
	return goid, nil
}

type Stack struct {
	Lo uintptr
	Hi uintptr
}

// GoStack returns the stack range of current goroutine
func GoStack() (Stack, error) {
	offset, err := Offsetof("runtime.g", "stack")
	if err != nil {
		return Stack{}, err
	}

	st := *(*Stack)(unsafe.Pointer(uintptr(getgptr()) + offset))
	return st, nil
}

// GoStackSize return stack size of current goroutine
func GoStackSize() (int, error) {
	stk, err := GoStack()
	if err != nil {
		return 0, err
	}

	size := stk.Hi - stk.Lo
	return int(size), nil
}
