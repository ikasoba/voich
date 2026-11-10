package pipewire

import (
	"runtime"

	"github.com/ebitengine/purego"
	"github.com/ikasoba/voich/pipewire/low"
)

type Size = uint

type LoopAPI struct {
	pw_loop_new            func(*low.CSpaDict) *low.CLoop
	pw_loop_destroy        func(*low.CLoop)
	pw_loop_set_name       func(*low.CLoop, string) int
	pw_loop_add_source     func(*low.CLoop, *low.CSpaSource) int
	pw_loop_update_source  func(*low.CLoop, *low.CSpaSource) int
	pw_loop_remove_source  func(*low.CLoop, *low.CSpaSource) int
	pw_loop_invoke         func(*low.CLoop, low.CSpaInvokeFunc, uint32, uintptr, uint, bool, uintptr) int
	pw_loop_locked         func(*low.CLoop, low.CSpaInvokeFunc, uint32, uintptr, Size, uintptr) int
	pw_loop_get_fd         func(*low.CLoop) int
	pw_loop_add_hook       func(*low.CLoop, *low.CSpaHook, *low.CSpaLoopControlHooks, uintptr)
	pw_loop_enter          func(*low.CLoop)
	pw_loop_leave          func(*low.CLoop)
	pw_loop_iterate        func(*low.CLoop, int) int
	pw_loop_check          func(*low.CLoop) int
	pw_loop_lock           func(*low.CLoop) int
	pw_loop_unlock         func(*low.CLoop) int
	pw_loop_get_time       func(*low.CLoop, *low.Timespec, int64) int
	pw_loop_wait           func(*low.CLoop, *low.Timespec) int
	pw_loop_signal         func(*low.CLoop, bool) int
	pw_loop_accept         func(*low.CLoop) int
	pw_loop_add_io         func(*low.CLoop, int, uint32, bool, low.SpaSourceIoFunc, uintptr) *low.CSpaSource
	pw_loop_update_io      func(*low.CLoop, *low.CSpaSource, uint32) int
	pw_loop_add_idle       func(*low.CLoop, bool, low.SpaSourceIdleFunc, uintptr) *low.CSpaSource
	pw_loop_enable_idle    func(*low.CLoop, *low.CSpaSource, bool) int
	pw_loop_add_event      func(*low.CLoop, low.SpaSourceEventFunc, uintptr) *low.CSpaSource
	pw_loop_signal_event   func(*low.CLoop, *low.CSpaSource) int
	pw_loop_add_timer      func(*low.CLoop, low.SpaSourceTimerFunc, uintptr) *low.CSpaSource
	pw_loop_update_timer   func(*low.CLoop, *low.CSpaSource, *low.Timespec, *low.Timespec, bool) int
	pw_loop_add_signal     func(*low.CLoop, int, low.SpaSourceSignalFunc, uintptr) *low.CSpaSource
	pw_loop_destroy_source func(*low.CLoop, *low.CSpaSource)
}

func (l *LoopAPI) Bind(libptr uintptr) {
	purego.RegisterLibFunc(&l.pw_loop_new, libptr, "pw_loop_new")
	purego.RegisterLibFunc(&l.pw_loop_destroy, libptr, "pw_loop_destroy")
	purego.RegisterLibFunc(&l.pw_loop_set_name, libptr, "pw_loop_set_name")
	purego.RegisterLibFunc(&l.pw_loop_add_source, libptr, "pw_loop_add_source")
	purego.RegisterLibFunc(&l.pw_loop_update_source, libptr, "pw_loop_update_source")
	purego.RegisterLibFunc(&l.pw_loop_remove_source, libptr, "pw_loop_remove_source")
	purego.RegisterLibFunc(&l.pw_loop_invoke, libptr, "pw_loop_invoke")
	purego.RegisterLibFunc(&l.pw_loop_locked, libptr, "pw_loop_locked")
	purego.RegisterLibFunc(&l.pw_loop_get_fd, libptr, "pw_loop_get_fd")
	purego.RegisterLibFunc(&l.pw_loop_add_hook, libptr, "pw_loop_add_hook")
	purego.RegisterLibFunc(&l.pw_loop_enter, libptr, "pw_loop_enter")
	purego.RegisterLibFunc(&l.pw_loop_leave, libptr, "pw_loop_leave")
	purego.RegisterLibFunc(&l.pw_loop_iterate, libptr, "pw_loop_iterate")
	purego.RegisterLibFunc(&l.pw_loop_check, libptr, "pw_loop_check")
	purego.RegisterLibFunc(&l.pw_loop_lock, libptr, "pw_loop_lock")
	purego.RegisterLibFunc(&l.pw_loop_unlock, libptr, "pw_loop_unlock")
	purego.RegisterLibFunc(&l.pw_loop_get_time, libptr, "pw_loop_get_time")
	purego.RegisterLibFunc(&l.pw_loop_wait, libptr, "pw_loop_wait")
	purego.RegisterLibFunc(&l.pw_loop_signal, libptr, "pw_loop_signal")
	purego.RegisterLibFunc(&l.pw_loop_accept, libptr, "pw_loop_accept")
	purego.RegisterLibFunc(&l.pw_loop_add_io, libptr, "pw_loop_add_io")
	purego.RegisterLibFunc(&l.pw_loop_update_io, libptr, "pw_loop_update_io")
	purego.RegisterLibFunc(&l.pw_loop_add_idle, libptr, "pw_loop_add_idle")
	purego.RegisterLibFunc(&l.pw_loop_enable_idle, libptr, "pw_loop_enable_idle")
	purego.RegisterLibFunc(&l.pw_loop_add_event, libptr, "pw_loop_add_event")
	purego.RegisterLibFunc(&l.pw_loop_signal_event, libptr, "pw_loop_signal_event")
	purego.RegisterLibFunc(&l.pw_loop_add_timer, libptr, "pw_loop_add_timer")
	purego.RegisterLibFunc(&l.pw_loop_update_timer, libptr, "pw_loop_update_timer")
	purego.RegisterLibFunc(&l.pw_loop_add_signal, libptr, "pw_loop_add_signal")
	purego.RegisterLibFunc(&l.pw_loop_destroy_source, libptr, "pw_loop_destroy_source")
}

func (l *LoopAPI) NewLoop(props *SpaDict) *Loop {
	var pinner runtime.Pinner

	cprops := props.MakeCSpaDict(&pinner)

	return l.bindLoop(l.pw_loop_new(&cprops))
}

func (l *LoopAPI) bindLoop(cloop *low.CLoop) *Loop {
	var loop Loop

	loop.api = l
	loop.cptr = cloop

	return &loop
}

type Loop struct {
	pinner runtime.Pinner

	api  *LoopAPI
	cptr *low.CLoop
}

func (l *Loop) Destroy() {
	defer l.pinner.Unpin()

	l.api.pw_loop_destroy(l.cptr)
}

func (l *Loop) SetName(name string) int {
	return l.api.pw_loop_set_name(l.cptr, name)
}

func (l *Loop) AddSource(src *low.CSpaSource) int {
	return l.api.pw_loop_add_source(l.cptr, src)
}

func (l *Loop) UpdateSource(src *low.CSpaSource) int {
	return l.api.pw_loop_update_source(l.cptr, src)
}

func (l *Loop) RemoveSource(src *low.CSpaSource) int {
	return l.api.pw_loop_remove_source(l.cptr, src)
}

func (l *Loop) Invoke(f low.CSpaInvokeFunc, seq uint32, data uintptr, size uint, block bool, userData uintptr) int {
	return l.api.pw_loop_invoke(l.cptr, f, seq, data, size, block, userData)
}

func (l *Loop) Locked(f low.CSpaInvokeFunc, seq uint32, data uintptr, size Size, userData uintptr) int {
	return l.api.pw_loop_locked(l.cptr, f, seq, data, size, userData)
}

func (l *Loop) GetFd() int {
	return l.api.pw_loop_get_fd(l.cptr)
}

func (l *Loop) AddHook(hook *low.CSpaHook, hooks *low.CSpaLoopControlHooks, data uintptr) {
	l.api.pw_loop_add_hook(l.cptr, hook, hooks, data)
}

func (l *Loop) Enter() {
	l.api.pw_loop_enter(l.cptr)
}

func (l *Loop) Leave() {
	l.api.pw_loop_leave(l.cptr)
}

func (l *Loop) Iterate(timeout int) int {
	return l.api.pw_loop_iterate(l.cptr, timeout)
}

func (l *Loop) Check() int {
	return l.api.pw_loop_check(l.cptr)
}

func (l *Loop) Lock() int {
	return l.api.pw_loop_lock(l.cptr)
}

func (l *Loop) Unlock() int {
	return l.api.pw_loop_unlock(l.cptr)
}

func (l *Loop) GetTime(abstime *low.Timespec, timeout int64) int {
	return l.api.pw_loop_get_time(l.cptr, abstime, timeout)
}

func (l *Loop) Wait(abstime *low.Timespec) int {
	return l.api.pw_loop_wait(l.cptr, abstime)
}

func (l *Loop) Signal(waitForAccept bool) int {
	return l.api.pw_loop_signal(l.cptr, waitForAccept)
}

func (l *Loop) Accept() int {
	return l.api.pw_loop_accept(l.cptr)
}

func (l *Loop) AddIO(fd int, mask uint32, close bool, f low.SpaSourceIoFunc, data uintptr) *low.CSpaSource {
	return l.api.pw_loop_add_io(l.cptr, fd, mask, close, f, data)
}

func (l *Loop) UpdateIO(src *low.CSpaSource, mask uint32) int {
	return l.api.pw_loop_update_io(l.cptr, src, mask)
}

func (l *Loop) AddIdle(enabled bool, f low.SpaSourceIdleFunc, data uintptr) *low.CSpaSource {
	return l.api.pw_loop_add_idle(l.cptr, enabled, f, data)
}

func (l *Loop) EnableIdle(src *low.CSpaSource, enabled bool) int {
	return l.api.pw_loop_enable_idle(l.cptr, src, enabled)
}

func (l *Loop) AddEvent(f low.SpaSourceEventFunc, data uintptr) *low.CSpaSource {
	return l.api.pw_loop_add_event(l.cptr, f, data)
}

func (l *Loop) SignalEvent(src *low.CSpaSource) int {
	return l.api.pw_loop_signal_event(l.cptr, src)
}

func (l *Loop) AddTimer(f low.SpaSourceTimerFunc, data uintptr) *low.CSpaSource {
	return l.api.pw_loop_add_timer(l.cptr, f, data)
}

func (l *Loop) UpdateTimer(src *low.CSpaSource, value *low.Timespec, interval *low.Timespec, absolute bool) int {
	return l.api.pw_loop_update_timer(l.cptr, src, value, interval, absolute)
}

func (l *Loop) AddSignal(signalNumber int, f low.SpaSourceSignalFunc, data uintptr) *low.CSpaSource {
	return l.api.pw_loop_add_signal(l.cptr, signalNumber, f, data)
}

func (l *Loop) DestroySource(src *low.CSpaSource) {
	l.api.pw_loop_destroy_source(l.cptr, src)
}
