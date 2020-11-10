package tl

import (
	"reflect"
)

//! мне не нравится, непонятно зачем
func haveFlag(v any) bool {
	typ := reflect.TypeOf(v)
	for i := 0; i < typ.NumField(); i++ {
		tag, found := typ.Field(i).Tag.Lookup("tl")
		if found {
			info, err := parseTag(tag)
			if err != nil {
				continue
			}

			if info.ignore {
				continue
			}

			return true
		}
	}

	return false
}

//! слайстрикс
func sliceToInterfaceSlice(in any) []any {
	if in == nil {
		return nil
	}

	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := make([]any, ival.Len())
	for i := 0; i < ival.Len(); i++ {
		res[i] = ival.Index(i).Interface()
	}

	return res
}
