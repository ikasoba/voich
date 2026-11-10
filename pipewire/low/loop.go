package low

type CLoop struct {
	System  *CSpaSystem
	Loop    *CSpaLoop
	Control *CSpaLoopControl
}

type CSpaSystem struct {
	Interface CSpaInterface
}

type CSpaInterface struct {
	TypeName *byte
	Version  uint32
	Cb       CSpaCallbacks
}

type CSpaLoop struct {
	Interface CSpaInterface
}

type CSpaLoopControl struct {
	Interface CSpaInterface
}

type CSpaSourceFunc uintptr
type CSpaInvokeFunc uintptr

type CSpaSource struct {
	Loop  *CSpaLoop
	Func  CSpaSourceFunc
	Data  uintptr
	Fd    int
	Mask  uint32
	RMask uint32
	Priv  uintptr
}

type CSpaLoopControlHooks struct {
	Version uint32
	Before  uintptr
	After   uintptr
}

type Timespec struct{}

type SpaSourceIoFunc uintptr
type SpaSourceIdleFunc uintptr
type SpaSourceEventFunc uintptr
type SpaSourceTimerFunc uintptr
type SpaSourceSignalFunc uintptr
