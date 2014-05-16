package easy

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
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

type stringChoice struct {
	Valid []string
	Value string
}

func (c *stringChoice) String() string {
	return fmt.Sprintf("%q", c.Value)
}

func (c *stringChoice) Set(v string) error {
	for _, i := range c.Valid {
		if i == v {
			c.Value = v
			return nil
		}
	}
	return errors.New(fmt.Sprintf("not one from %q", c.Valid))
}

func StringChoice(name string, value string, usage string, choices []string) *string {
	f := &stringChoice{choices, value}
	flag.Var(f, name, fmt.Sprintf("%s (one of %q)", usage, choices))
	return &f.Value
}

type intChoice struct {
	Valid []int
	Value int
}

func (c *intChoice) String() string {
	return strconv.Itoa(c.Value)
}

func (c *intChoice) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	for _, j := range c.Valid {
		if i == j {
			c.Value = i
			return nil
		}
	}
	return errors.New(fmt.Sprintf("not one of %q", c.Valid))
}

func IntChoice(name string, value int, usage string, choices []int) *int {
	f := &intChoice{choices, value}
	flag.Var(f, name, fmt.Sprintf("%s (one of %q)", usage, choices))
	return &f.Value
}
