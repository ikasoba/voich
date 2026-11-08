package pipewire

import (
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ikasoba/voich/pipewire/low"
)

type SpaHandle struct {
	pinner runtime.Pinner

	api    *InitializationAPI
	handle *low.CSpaHandle

	version      uint32
	getInterface func(*low.CSpaHandle, string, *unsafe.Pointer) int
	clear        func(*low.CSpaHandle) int
}

func NewSpaHandle(api *InitializationAPI, h *low.CSpaHandle) *SpaHandle {
	var s SpaHandle

	s.api = api
	s.handle = h
	s.pinner.Pin(h)
	s.version = h.Version

	purego.RegisterFunc(&s.getInterface, h.GetInterface)
	purego.RegisterFunc(&s.clear, h.Clear)

	return &s
}

func (s *SpaHandle) Version() uint32 {
	return s.version
}

func (s *SpaHandle) GetInterface(typeName SpaTypeName, dst *unsafe.Pointer) int {
	return s.getInterface(s.handle, string(typeName), dst)
}

func (s *SpaHandle) Clear() int {
	return s.clear(s.handle)
}

func (s *SpaHandle) Unload() int {
	defer s.pinner.Unpin()

	return s.api.pw_unload_spa_handle(s.handle)
}
