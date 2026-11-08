package pipewire

import (
	"runtime"
	"slices"

	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ikasoba/voich/pipewire/low"
	"github.com/ikasoba/voich/pipewire/utils"
)

type InitializationAPI struct {
	pw_init                      func(*int, **unsafe.Pointer)
	pw_deinit                    func()
	pw_debug_is_category_enabled func(string) bool
	pw_get_application_name      func() string
	pw_get_prgname               func() string
	pw_get_user_name             func() string
	pw_get_host_name             func() string
	pw_get_client_name           func() string
	pw_check_option              func(string, string) bool
	pw_direction_reverse         func(SpaDirection) SpaDirection
	pw_set_domain                func(string) int
	pw_get_domain                func() string
	pw_get_support               func(*low.CSpaSupport, uint32) uint32
	pw_load_spa_handle           func(string, string, *low.CSpaDict, uint32, *low.CSpaSupport) *low.CSpaHandle
	pw_unload_spa_handle         func(*low.CSpaHandle) int
}

func (i *InitializationAPI) Bind(libptr uintptr) {
	purego.RegisterLibFunc(&i.pw_init, libptr, "pw_init")
	purego.RegisterLibFunc(&i.pw_deinit, libptr, "pw_deinit")
	purego.RegisterLibFunc(&i.pw_debug_is_category_enabled, libptr, "pw_debug_is_category_enabled")
	purego.RegisterLibFunc(&i.pw_get_application_name, libptr, "pw_get_application_name")
	purego.RegisterLibFunc(&i.pw_get_prgname, libptr, "pw_get_prgname")
	purego.RegisterLibFunc(&i.pw_get_user_name, libptr, "pw_get_user_name")
	purego.RegisterLibFunc(&i.pw_get_host_name, libptr, "pw_get_host_name")
	purego.RegisterLibFunc(&i.pw_get_client_name, libptr, "pw_get_client_name")
	purego.RegisterLibFunc(&i.pw_check_option, libptr, "pw_check_option")
	purego.RegisterLibFunc(&i.pw_direction_reverse, libptr, "pw_direction_reverse")
	purego.RegisterLibFunc(&i.pw_set_domain, libptr, "pw_set_domain")
	purego.RegisterLibFunc(&i.pw_get_domain, libptr, "pw_get_domain")
	purego.RegisterLibFunc(&i.pw_get_support, libptr, "pw_get_support")
	purego.RegisterLibFunc(&i.pw_load_spa_handle, libptr, "pw_load_spa_handle")
	purego.RegisterLibFunc(&i.pw_unload_spa_handle, libptr, "pw_unload_spa_handle")
}

// Initialize PipeWire.
func (i *InitializationAPI) Init(args []string) []string {
	var pin runtime.Pinner
	defer pin.Unpin()

	n := int(len(args))
	args = slices.Clone(args)
	argv := make([]unsafe.Pointer, len(args))

	for i := range argv {
		args[i] += "\x00"
		argv[i] = unsafe.Pointer(unsafe.StringData(args[i]))
		pin.Pin(argv[i])
	}

	argvp := unsafe.SliceData(argv)

	i.pw_init(&n, &argvp)

	argv = unsafe.Slice(argvp, n)
	modified := make([]string, n)
	for i := range n {
		modified[i] = utils.CloneFromCString(argv[i], 1024)
	}

	return modified
}

// Deinitialize PipeWire.
func (i *InitializationAPI) Deinit() {
	i.pw_deinit()
}

// Check if a debug category is enabled.
func (i *InitializationAPI) DebugIsCategoryEnabled(name string) bool {
	return i.pw_debug_is_category_enabled(name)
}

// Get the application name.
func (i *InitializationAPI) GetApplicationName() string {
	return i.pw_get_application_name()
}

// Get the program name.
func (i *InitializationAPI) GetPrgName() string {
	return i.pw_get_prgname()
}

// Get the user name.
func (i *InitializationAPI) GetUserName() string {
	return i.pw_get_user_name()
}

// Get the host name.
func (i *InitializationAPI) GetHostName() string {
	return i.pw_get_host_name()
}

// Get the client name.
func (i *InitializationAPI) GetClientName() string {
	return i.pw_get_client_name()
}

func (i *InitializationAPI) CheckOption(name string, value string) bool {
	return i.pw_check_option(name, value)
}

// Reverse the direction.
func (i *InitializationAPI) DirectionReverse(d SpaDirection) SpaDirection {
	return i.pw_direction_reverse(d)
}

func (i *InitializationAPI) SetDomain(domain string) int {
	return i.pw_set_domain(domain)
}

func (i *InitializationAPI) GetDomain() string {
	return i.pw_get_domain()
}

func (i *InitializationAPI) GetSupport(dst []SpaSupport) uint32 {
	var pin runtime.Pinner
	defer pin.Unpin()

	tmp := make([]low.CSpaSupport, len(dst))
	dptr := unsafe.SliceData(tmp)
	pin.Pin(dptr)

	count := i.pw_get_support(dptr, uint32(len(dst)))

	tmp = tmp[:count]

	for i, s := range tmp {
		dst[i] = SpaSupport{
			TypeName: SpaTypeName(utils.CloneFromCString(unsafe.Pointer(s.TypeName), 1024)),
			Data:     s.Data,
		}
	}

	return count
}

func (i *InitializationAPI) LoadSpaHandle(lib string, factoryName string, info SpaDict, supports []SpaSupport) *SpaHandle {
	var pinner runtime.Pinner
	defer pinner.Unpin()

	numSupports := len(supports)
	csupports := make([]low.CSpaSupport, numSupports)
	pinner.Pin(csupports)

	for i, v := range supports {
		csupports[i] = v.MakeCSpaSupport(&pinner)
	}

	cinfo := new(low.CSpaDict)

	*cinfo = info.MakeCSpaDict(&pinner)
	pinner.Pin(cinfo)

	handle := i.pw_load_spa_handle(lib, factoryName, cinfo, uint32(numSupports), unsafe.SliceData(csupports))
	if handle == nil {
		return nil
	}

	return NewSpaHandle(i, handle)
}
