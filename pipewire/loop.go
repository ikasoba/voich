package pipewire

import (
	"runtime"

	"github.com/ikasoba/voich/pipewire/low"
)

type Loop struct {
	pinner runtime.Pinner

	api  *MainLoopAPI
	cptr *low.CLoop
}
