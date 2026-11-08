package pipewire

import "fmt"

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
