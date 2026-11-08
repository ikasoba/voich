package low

type CMainLoop struct{}

type CSpaHook struct {
	Link    CSpaList[CSpaHook]
	Cb      CSpaCallbacks
	Removed uintptr
	Priv    uintptr
}

type CSpaList[T any] struct {
	Next *T
	Prev *T
}

type CSpaCallbacks struct {
	Funcs uintptr
	Data  uintptr
}

const PW_VERSION_MAIN_LOOP_EVENTS = 0

type CMainLoopEvents struct {
	Version uint32
	Destroy uintptr
}
