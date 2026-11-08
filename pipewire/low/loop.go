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
