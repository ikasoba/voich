package low

type CSpaSupport struct {
	TypeName *byte
	Data     uintptr
}

type CSpaDict struct {
	Flags  uint32
	NItems uint32
	Items  *CSpaDictItem
}

type CSpaDictItem struct {
	Key   *byte
	Value *byte
}

type CSpaHandle struct {
	Version      uint32
	GetInterface uintptr
	Clear        uintptr
}
