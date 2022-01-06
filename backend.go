package goplayground

import (
	"flag"
	"fmt"
)

// Backend indicates run Go environment.
type Backend string

const (
	// BackendDefault indicates default environment (latest release version).
	BackendDefault Backend = ""
	// BackendGotip indicates using the develop branch.
	BackendGotip Backend = "gotip"
)

// String implements flag.Value.String.
func (b Backend) String() string {
	return string(b)
}

// Set implements flag.Value.Set.
func (b *Backend) Set(s string) error {
	switch Backend(s) {
	case BackendDefault, BackendGotip:
		*b = Backend(s)
		return nil
	default:
		return fmt.Errorf("unexpected backend: %s", s)
	}
}

var _ flag.Value = (*Backend)(nil)
