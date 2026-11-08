package pipewire

import (
	"runtime"

	"github.com/ikasoba/voich/pipewire/low"
	"github.com/ikasoba/voich/pipewire/utils"
)

type SpaTypeName string

const (
	SpaTypeInfoBase SpaTypeName = "Spa:"
)

type SpaSupportDescriptor interface {
	Type() SpaTypeName
}

type SpaSupport struct {
	TypeName SpaTypeName
	Data     uintptr
}

func (s *SpaSupport) Type() SpaTypeName {
	return s.TypeName
}

func (s *SpaSupport) MakeCSpaSupport(pinner *runtime.Pinner) low.CSpaSupport {
	return low.CSpaSupport{
		TypeName: utils.NewCString(pinner, string(s.TypeName)),
		Data:     s.Data,
	}
}
