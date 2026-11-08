package utils

import (
	"runtime"
	"unsafe"
)

func CloneFromCString(sptr unsafe.Pointer, max int) string {
	pz := unsafe.Sizeof(byte(0))
	ln := 0

	if uintptr(sptr) == 0 {
		return ""
	}

	for i := range max {
		b := *(*byte)(unsafe.Add(sptr, pz*uintptr(i)))
		if b == 0 {
			break
		}

		ln++
	}

	return unsafe.String((*byte)(sptr), uintptr(ln))
}

func NewCString(pinner *runtime.Pinner, s string) *byte {
	tmp := s + "\x00"
	pinner.Pin(&tmp)

	return unsafe.StringData(tmp)
}
