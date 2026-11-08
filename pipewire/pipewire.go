package pipewire

import (
	"errors"
	"github.com/ebitengine/purego"
)

var PipewireLibraryNames = []string{
	"libpipewire-0.3.so",
	"libpipewire-0.3.so.0",
}

type Pipewire struct {
	libptr uintptr

	// Misc

	pw_get_library_version func() string

	// Initialization

	InitializationAPI

	// Main Loop

	MainLoopAPI
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

	purego.RegisterLibFunc(&pw.pw_get_library_version, pw.libptr, "pw_get_library_version")

	pw.InitializationAPI.Bind(pw.libptr)

	pw.MainLoopAPI.Bind(pw.libptr)

	return pw, nil
}

func (p *Pipewire) GetLibraryVersion() string {
	return p.pw_get_library_version()
}
