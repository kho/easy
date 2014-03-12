package easy

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// AddFlags adds flags to fs from ptr. ptr must be a pointer to a
// struct type. Every exported field in the struct is added as
// flags. Only built-in flag variable types execept Duration (bool,
// float64, int64, int, string, uint64, uint) are supported. The name
// (or usage, respectively) of a field can be customized with a field
// tag named "name" (or "usage", respectively).
func AddFlags(ptr interface{}, fs *flag.FlagSet) {
	forEachField(ptr, func(field reflect.StructField, value reflect.Value) {
		if field.PkgPath == "" {
			addFieldFlag(field, value, fs)
		}
	})
}

func forEachField(ptr interface{}, action func(reflect.StructField, reflect.Value)) {
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
		action(field, value)
	}
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
