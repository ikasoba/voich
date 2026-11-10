package pipewire

import (
	"runtime"
	"unsafe"

	"github.com/ikasoba/voich/pipewire/low"
	"github.com/ikasoba/voich/pipewire/utils"
)

type SpaDict struct {
	Flags uint32
	Items map[string]string
}

func (s *SpaDict) MakeCSpaDict(pinner *runtime.Pinner) low.CSpaDict {
	n := len(s.Items)
	info := low.CSpaDict{
		Flags:  s.Flags,
		NItems: uint32(n),
		Items:  nil,
	}

	items := make([]low.CSpaDictItem, n)

	i := 0
	for k, v := range s.Items {
		if i >= n {
			break
		}

		item := low.CSpaDictItem{
			Key:   utils.NewCString(pinner, k),
			Value: utils.NewCString(pinner, v),
		}

		items[i] = item
		i++
	}

	info.Items = unsafe.SliceData(items)
	pinner.Pin(items)

	return info
}
