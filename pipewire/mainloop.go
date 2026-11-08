package pipewire

import (
	"runtime"

	"github.com/ebitengine/purego"
	"github.com/ikasoba/voich/pipewire/low"
)

type MainLoopAPI struct {
	pw_main_loop_new          func(*low.CSpaDict) *low.CMainLoop
	pw_main_loop_add_listener func(*low.CMainLoop, *low.CSpaHook, *low.CMainLoopEvents, uintptr)
	pw_main_loop_get_loop     func(*low.CMainLoop) *low.CLoop
}

func (m *MainLoopAPI) Bind(libptr uintptr) {
	purego.RegisterLibFunc(&m.pw_main_loop_new, libptr, "pw_main_loop_new")
	purego.RegisterLibFunc(&m.pw_main_loop_add_listener, libptr, "pw_main_loop_add_listener")
	purego.RegisterLibFunc(&m.pw_main_loop_get_loop, libptr, "pw_main_loop_get_loop")
}

func (m *MainLoopAPI) NewMainLoop(props *SpaDict) *MainLoop {
	var pinner runtime.Pinner
	defer pinner.Unpin()

	var cprops *low.CSpaDict

	if props != nil {
		cprops = new(low.CSpaDict)

		*cprops = props.MakeCSpaDict(&pinner)
		pinner.Pin(cprops)
	}

	var l MainLoop

	l.api = m
	l.cptr = m.pw_main_loop_new(cprops)

	return &l
}

type MainLoop struct {
	pinner runtime.Pinner

	api  *MainLoopAPI
	cptr *low.CMainLoop
	loop *Loop
}

func (m *MainLoop) GetLoop() *Loop {
	if m.loop != nil {
		return m.loop
	}

	cloop := m.api.pw_main_loop_get_loop(m.cptr)

	var loop Loop
	loop.api = m.api
	loop.cptr = cloop

	m.loop = &loop

	return m.loop
}
