package pipewire

import (
	"errors"
	"fmt"
	"runtime"
	"slices"

	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ikasoba/voich/pipewire/low"
	"github.com/ikasoba/voich/pipewire/utils"
)

var PipewireLibraryNames = []string{
	"libpipewire-0.3.so",
	"libpipewire-0.3.so.0",
}

// reference spa/include/spa/utils/defs.h
type SpaDirection int

func (d SpaDirection) String() string {
	switch d {
	case SpaDirectionInput:
		return "(SpaDirectionInput 0)"

	case SpaDirectionOutput:
		return "(SpaDirectionOutput 1)"

	default:
		return fmt.Sprintf("(SpaDirection!UNK %d)", int(d))
	}
}

const (
	SpaDirectionInput  SpaDirection = 0
	SpaDirectionOutput              = 1
)

type Pipewire struct {
	libptr uintptr

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

func (s *SpaSupport) MakeCSpaSupport(pinner *runtime.Pinner) low.CSpaSupport {
	return low.CSpaSupport{
		TypeName: utils.NewCString(pinner, string(s.TypeName)),
		Data:     s.Data,
	}
}

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
	}

	info.Items = unsafe.SliceData(items)
	pinner.Pin(items)

	return info
}

type SpaHandle struct {
	pinner runtime.Pinner

	pw     *Pipewire
	handle *low.CSpaHandle

	version      uint32
	getInterface func(*low.CSpaHandle, string, *unsafe.Pointer) int
	clear        func(*low.CSpaHandle) int
}

func NewSpaHandle(pw *Pipewire, h *low.CSpaHandle) *SpaHandle {
	var s SpaHandle

	s.pw = pw
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

	return s.pw.pw_unload_spa_handle(s.handle)
}

func loadLibptr() (uintptr, error) {
	var errs []error

	for _, name := range PipewireLibraryNames {
		lib, err := purego.Dlopen(name, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if lib != 0 {
			return lib, nil
		}
	}

	return 0, errors.Join(errs...)
}

func LoadPipewire() (pw *Pipewire, err error) {
	pw = &Pipewire{}
	pw.libptr, err = loadLibptr()
	if err != nil {
		return nil, err
	}

	purego.RegisterLibFunc(&pw.pw_init, pw.libptr, "pw_init")
	purego.RegisterLibFunc(&pw.pw_deinit, pw.libptr, "pw_deinit")
	purego.RegisterLibFunc(&pw.pw_debug_is_category_enabled, pw.libptr, "pw_debug_is_category_enabled")
	purego.RegisterLibFunc(&pw.pw_get_application_name, pw.libptr, "pw_get_application_name")
	purego.RegisterLibFunc(&pw.pw_get_prgname, pw.libptr, "pw_get_prgname")
	purego.RegisterLibFunc(&pw.pw_get_user_name, pw.libptr, "pw_get_user_name")
	purego.RegisterLibFunc(&pw.pw_get_host_name, pw.libptr, "pw_get_host_name")
	purego.RegisterLibFunc(&pw.pw_get_client_name, pw.libptr, "pw_get_client_name")
	purego.RegisterLibFunc(&pw.pw_check_option, pw.libptr, "pw_check_option")
	purego.RegisterLibFunc(&pw.pw_direction_reverse, pw.libptr, "pw_direction_reverse")
	purego.RegisterLibFunc(&pw.pw_set_domain, pw.libptr, "pw_set_domain")
	purego.RegisterLibFunc(&pw.pw_get_domain, pw.libptr, "pw_get_domain")
	purego.RegisterLibFunc(&pw.pw_get_support, pw.libptr, "pw_get_support")
	purego.RegisterLibFunc(&pw.pw_load_spa_handle, pw.libptr, "pw_load_spa_handle")
	purego.RegisterLibFunc(&pw.pw_unload_spa_handle, pw.libptr, "pw_unload_spa_handle")

	return pw, nil
}

// Initialize PipeWire.
func (p *Pipewire) Init(args []string) []string {
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

	p.pw_init(&n, &argvp)

	argv = unsafe.Slice(argvp, n)
	modified := make([]string, n)
	for i := range n {
		modified[i] = utils.CloneFromCString(argv[i], 1024)
	}

	return modified
}

// Deinitialize PipeWire.
func (p *Pipewire) Deinit() {
	p.pw_deinit()
}

// Check if a debug category is enabled.
func (p *Pipewire) DebugIsCategoryEnabled(name string) bool {
	return p.pw_debug_is_category_enabled(name)
}

// Get the application name.
func (p *Pipewire) GetApplicationName() string {
	return p.pw_get_application_name()
}

// Get the program name.
func (p *Pipewire) GetPrgName() string {
	return p.pw_get_prgname()
}

// Get the user name.
func (p *Pipewire) GetUserName() string {
	return p.pw_get_user_name()
}

// Get the host name.
func (p *Pipewire) GetHostName() string {
	return p.pw_get_host_name()
}

// Get the client name.
func (p *Pipewire) GetClientName() string {
	return p.pw_get_client_name()
}

func (p *Pipewire) CheckOption(name string, value string) bool {
	return p.pw_check_option(name, value)
}

// Reverse the direction.
func (p *Pipewire) DirectionReverse(d SpaDirection) SpaDirection {
	return p.pw_direction_reverse(d)
}

func (p *Pipewire) SetDomain(domain string) int {
	return p.pw_set_domain(domain)
}

func (p *Pipewire) GetDomain() string {
	return p.pw_get_domain()
}

func (p *Pipewire) GetSupport(dst []SpaSupport) uint32 {
	var pin runtime.Pinner
	defer pin.Unpin()

	tmp := make([]low.CSpaSupport, len(dst))
	dptr := unsafe.SliceData(tmp)
	pin.Pin(dptr)

	count := p.pw_get_support(dptr, uint32(len(dst)))

	tmp = tmp[:count]

	for i, s := range tmp {
		dst[i] = SpaSupport{
			TypeName: SpaTypeName(utils.CloneFromCString(unsafe.Pointer(s.TypeName), 1024)),
			Data:     s.Data,
		}
	}

	return count
}

func (p *Pipewire) LoadSpaHandle(lib string, factoryName string, info SpaDict, supports []SpaSupport) *SpaHandle {
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

	handle := p.pw_load_spa_handle(lib, factoryName, cinfo, uint32(numSupports), unsafe.SliceData(csupports))
	if handle == nil {
		return nil
	}

	return NewSpaHandle(p, handle)
}
