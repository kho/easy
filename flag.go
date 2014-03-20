package easy

import (
	"strings"
)

// Comma separated list of strings as flag.Value.
type Strings []string

func (s *Strings) String() string {
	return strings.Join([]string(*s), ",")
}

func (s *Strings) Set(v string) error {
	*s = nil
	for _, i := range strings.Split(v, ",") {
		if i != "" {
			*s = append(*s, i)
		}
	}
	return nil
}
