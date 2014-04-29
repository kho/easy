package easy

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// AddFlags adds flags to fs from ptr. ptr must be a pointer to a
// struct type. Every exported field in the struct is added as a
// flag. Only built-in flag variable types except Duration (bool,
// float64, int64, int, string, uint64, uint) are supported. The name
// (or usage, respectively) of a field can be customized with a field
// tag named "name" (or "usage", respectively).
func AddFlags(ptr interface{}, fs *flag.FlagSet) {
	forEachField(ptr, func(field reflect.StructField, value reflect.Value) error {
		if field.PkgPath == "" {
			addFieldFlag(field, value, fs)
		}
		return nil
	})
}

// SetArgs sets the given ptr from command line arguments. ptr must be
// a pointer to a struct type. Every exported field is processed in
// the order of declaration. Only bool, float64, int64, int, string,
// uint64, uint and slices of these types can be set. The name (or
// usage, respectively) of a field can be customized with a field tag
// named "name" (or "usage", respectively).
func SetArgs(ptr interface{}, args []string) error {
	i := 0
	if err := forEachField(ptr, func(field reflect.StructField, value reflect.Value) error {
		if field.PkgPath == "" {
			n, err := setField(field, value, args[i:])
			if err != nil {
				return err
			}
			i += n
		}
		return nil
	}); err != nil {
		return err
	}
	if i != len(args) {
		return errors.New("extra arguments: " + fmt.Sprintf("%q", args[i:]))
	}
	return nil
}

// ParseFlagsAndArgs parses the standard flags and then sets the
// arguments to ptr. When ptr is nil, flags are still parsed but no
// arguments will be procssed. This does certain magic with flags as
// well (e.g. tweaking glog).
func ParseFlagsAndArgs(ptr interface{}) {
	flag.Usage = func() {
		CombinedUsage(os.Args[0], ptr, flag.PrintDefaults)
	}
	// When using glog, I would like to log to stderr by default.
	if f := flag.Lookup("logtostderr"); f != nil {
		f.DefValue = "true"
	}
	flag.Parse()
	if ptr == nil {
		return
	}
	if err := SetArgs(ptr, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}
}

// ParseFlagsAndArgsWith parses a specified flagset and the sets the
// arguments to ptr. The flagset may be nil, in which case no flags
// except -h are processed. ptr may also be nil, in which case no
// arguments will be processed.
func ParseFlagsAndArgsWith(name string, ptr interface{}, fs *flag.FlagSet, args []string) error {
	if fs == nil {
		fs = flag.NewFlagSet("", 0)
	}
	fs.Usage = func() {
		CombinedUsage(name, ptr, fs.PrintDefaults)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if ptr == nil {
		return nil
	}
	return SetArgs(ptr, fs.Args())
}

func CombinedUsage(name string, ptr interface{}, printDefaults func()) {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] %s\n", name, strings.Join(FieldNames(ptr), " "))
	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	PrintArguments(ptr)
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	printDefaults()
}

func FieldNames(ptr interface{}) (s []string) {
	forEachField(ptr, func(field reflect.StructField, _ reflect.Value) error {
		if field.PkgPath == "" {
			name, _ := getFieldNameUsage(field)
			if field.Type.Kind() == reflect.Slice {
				name = "[" + name + " ...]"
			}
			s = append(s, name)
		}
		return nil
	})
	return
}

func PrintArguments(ptr interface{}) {
	forEachField(ptr, func(field reflect.StructField, _ reflect.Value) error {
		if field.PkgPath == "" {
			name, usage := getFieldNameUsage(field)
			fmt.Fprintf(os.Stderr, "  %s: %s\n", name, usage)
		}
		return nil
	})
}

func forEachField(ptr interface{}, action func(reflect.StructField, reflect.Value) error) error {
	ptrValue := reflect.ValueOf(ptr)
	if ptrValue.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("required pointer to struct; got %v", ptrValue))
	}
	elemValue := ptrValue.Elem()
	if elemValue.Kind() != reflect.Struct {
		panic(fmt.Sprintf("required pointer to struct; got %v", ptrValue))
	}
	elemType := elemValue.Type()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		value := elemValue.Field(i)
		if err := action(field, value); err != nil {
			return err
		}
	}
	return nil
}

func addFieldFlag(field reflect.StructField, value reflect.Value, fs *flag.FlagSet) {
	name, usage := getFieldNameUsage(field)
	ptr := unsafe.Pointer(value.UnsafeAddr())
	switch field.Type.Kind() {
	case reflect.Bool:
		fs.BoolVar((*bool)(ptr), name, *(*bool)(ptr), usage)
	case reflect.Float64:
		fs.Float64Var((*float64)(ptr), name, *(*float64)(ptr), usage)
	case reflect.Int64:
		fs.Int64Var((*int64)(ptr), name, *(*int64)(ptr), usage)
	case reflect.Int:
		fs.IntVar((*int)(ptr), name, *(*int)(ptr), usage)
	case reflect.String:
		fs.StringVar((*string)(ptr), name, *(*string)(ptr), usage)
	case reflect.Uint64:
		fs.Uint64Var((*uint64)(ptr), name, *(*uint64)(ptr), usage)
	case reflect.Uint:
		fs.UintVar((*uint)(ptr), name, *(*uint)(ptr), usage)
	default:
		panic(fmt.Sprintf("unsupported type: %v", field.Type))
	}
}

func getFieldNameUsage(field reflect.StructField) (name, usage string) {
	name = field.Tag.Get("name")
	if name == "" {
		name = strings.ToLower(field.Name)
	}
	usage = field.Tag.Get("usage")
	return
}

type nameError struct {
	Name string
	Err  error
}

func (e *nameError) Error() string {
	return e.Name + ": " + e.Err.Error()
}

func newNameError(name string, err error) error {
	if err == nil {
		return nil
	}
	return &nameError{name, err}
}

func setField(field reflect.StructField, value reflect.Value, args []string) (int, error) {
	name, _ := getFieldNameUsage(field)
	if field.Type.Kind() == reflect.Slice {
		n, err := setSlice(value, args)
		return n, newNameError(name, err)
	}
	if len(args) == 0 {
		return 0, newNameError(name, errors.New("missing argument"))
	}
	return 1, newNameError(name, setValue(value, args[0]))
}

func setValue(value reflect.Value, s string) error {
	ptr := unsafe.Pointer(value.UnsafeAddr())
	switch value.Type().Kind() {
	case reflect.Bool:
		return setBool(ptr, s)
	case reflect.Float64:
		return setFloat64(ptr, s)
	case reflect.Int64:
		return setInt64(ptr, s)
	case reflect.Int:
		return setInt(ptr, s)
	case reflect.String:
		return setString(ptr, s)
	case reflect.Uint64:
		return setUint64(ptr, s)
	case reflect.Uint:
		return setUint(ptr, s)
	default:
		panic(fmt.Sprintf("unsupported type: %v", value.Type()))
	}
}

func setSlice(value reflect.Value, args []string) (n int, err error) {
	value.SetLen(0)
	slice := value
	for _, i := range args {
		x := reflect.New(value.Type().Elem())
		if err = setValue(x.Elem(), i); err != nil {
			return
		}
		slice = reflect.Append(slice, x.Elem())
		n++
	}
	value.Set(slice)
	return
}

func setBool(ptr unsafe.Pointer, s string) error {
	b, err := strconv.ParseBool(s)
	if err == nil {
		*(*bool)(ptr) = b
	}
	return err
}

func setFloat64(ptr unsafe.Pointer, s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*(*float64)(ptr) = f
	}
	return err
}

func setInt64(ptr unsafe.Pointer, s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		*(*int64)(ptr) = i
	}
	return err
}

func setInt(ptr unsafe.Pointer, s string) error {
	i, err := strconv.ParseInt(s, 10, 0)
	if err == nil {
		*(*int)(ptr) = int(i)
	}
	return err
}

func setString(ptr unsafe.Pointer, s string) error {
	*(*string)(ptr) = s
	return nil
}

func setUint64(ptr unsafe.Pointer, s string) error {
	u, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		*(*uint64)(ptr) = u
	}
	return err
}

func setUint(ptr unsafe.Pointer, s string) error {
	u, err := strconv.ParseUint(s, 10, 0)
	if err == nil {
		*(*uint)(ptr) = uint(u)
	}
	return err
}
