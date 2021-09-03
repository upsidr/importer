package errorsplus

import (
	"errors"
	"fmt"
)

// Errors represents a slice of errors.
type Errors []error

func (es Errors) Error() string {
	switch len(es) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", es[0])
	}

	rt := "more than one error occurred:"
	for _, e := range es {
		if e != nil {
			rt = fmt.Sprintf("%s\n\t%v", rt, e)
		}
	}
	return rt
}

func (es Errors) Is(target error) bool {
	if len(es) == 0 {
		return false
	}

	for _, e := range es {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}
