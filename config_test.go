package easy

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

type justBool bool

type goodConfig struct {
	// Supported types.
	Bool     bool
	Float64  float64
	Int64    int64
	Int      int
	String   string
	Uint64   uint64
	Uint     uint
	JustBool justBool `name:"just-bool"`
	// Name and usage.
	Flag int `name:"sound" usage:"quack quack"`
	// Unexported fields.
	hidden int `usage:"you don't see me"`
}

type badConfig struct {
	// OK.
	Bool bool
	// Not OK.
	Struct struct{ N, M int }
}

type sliceConfig struct {
	First int
	Rest  []int
}

func TestAddFlags(t *testing.T) {
	// Only take pointer to structs.
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		var x goodConfig
		AddFlags(x, flag.NewFlagSet("", 0))
	}()
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		var x int
		AddFlags(&x, flag.NewFlagSet("", 0))
	}()
	// All exported fields must be of supported types.
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		var x badConfig
		AddFlags(&x, flag.NewFlagSet("", 0))
	}()
	// Name and usage.
	func() {
		x := goodConfig{}
		fs := flag.NewFlagSet("", 0)
		AddFlags(&x, fs)
		// Flag with default name and usage.
		if f := fs.Lookup("bool"); f == nil {
			t.Errorf("-bool is not defined")
		} else if f.Usage != "" {
			t.Errorf("-bool has non-empty usage")
		}
		// Flag with custom name.
		if f := fs.Lookup("just-bool"); f == nil {
			t.Errorf("-just-bool is not defined")
		}
		// Flag with custom name and usage.
		if f := fs.Lookup("sound"); f == nil {
			t.Errorf("-sound is not defined")
		} else if f.Usage != "quack quack" {
			t.Errorf("incorrect usage for -sound: %q", f.Usage)
		}
		// hidden is not a flag.
		if f := fs.Lookup("hidden"); f != nil {
			t.Errorf("-hidden is defined")
		}
	}()
	// Zero default.
	func() {
		x := goodConfig{}
		fs := flag.NewFlagSet("", 0)
		AddFlags(&x, fs)
		// Parse and set values.
		if err := fs.Parse(strings.Fields("-bool=true -float64=1 -int64=2 -int=3 -string=a -uint64=4 -uint=5 -just-bool=true -sound=6")); err != nil {
			t.Errorf("error in parsing flags: %v", err)
		} else {
			y := goodConfig{true, 1, 2, 3, "a", 4, 5, true, 6, 0}
			if x != y {
				t.Errorf("after parsing got %+v; expected %+v", x, y)
			}
		}
	}()
	// Non-zero default.
	func() {
		x := goodConfig{String: "ok"}
		fs := flag.NewFlagSet("", 0)
		AddFlags(&x, fs)
		if f := fs.Lookup("string"); f == nil {
			t.Errorf("-string is not defined")
		} else if f.DefValue != "ok" {
			t.Errorf("default value of -string is %q; expected %q", f.DefValue, "ok")
		}
		// Parse and change default value.
		if err := fs.Parse([]string{"-string=not ok"}); err != nil {
			t.Errorf("error in parsing flags: %v", err)
		} else if x.String != "not ok" {
			t.Errorf("-string is set to %q; expected %q", x.String, "not ok")
		}
	}()
}

func TestSetArgs(t *testing.T) {
	// Simple case.
	func() {
		var c goodConfig
		if err := SetArgs(&c, []string{"true", "1", "2", "3", "s", "5", "6", "true", "8"}); err != nil {
			t.Error("unexpected error: ", err)
		} else {
			y := goodConfig{true, 1, 2, 3, "s", 5, 6, true, 8, 0}
			if c != y {
				t.Errorf("expected %+v; got %+v", y, c)
			}
		}
	}()
	// Slice.
	func() {
		var c sliceConfig
		if err := SetArgs(&c, []string{"1", "2", "3"}); err != nil {
			t.Error("unexpected error: ", err)
		} else {
			y := sliceConfig{1, []int{2, 3}}
			if !reflect.DeepEqual(c, y) {
				t.Errorf("expected %+v; got %+v", y, c)
			}
		}
		if err := SetArgs(&c, []string{"4", "5"}); err != nil {
			t.Error("unexpected error: ", err)
		} else {
			y := sliceConfig{4, []int{5}}
			if !reflect.DeepEqual(c, y) {
				t.Errorf("expected %+v; got %+v", y, c)
			}
		}
		if err := SetArgs(&c, []string{"6"}); err != nil {
			t.Error("unexpected error: ", err)
		} else {
			y := sliceConfig{6, []int{}}
			if !reflect.DeepEqual(c, y) {
				t.Errorf("expected %+v; got %+v", y, c)
			}
		}
	}()
}
