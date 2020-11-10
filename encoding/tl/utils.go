package tl

import (
	"reflect"
)

func haveFlag(v interface{}) bool {
	typ := reflect.TypeOf(v)
	for i := 0; i < typ.NumField(); i++ {
		tag, found := typ.Field(i).Tag.Lookup(tagName)
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

func sliceToInterfaceSlice(in interface{}) []interface{} {
	if in == nil {
		return nil
	}

	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := make([]interface{}, ival.Len())
	for i := 0; i < ival.Len(); i++ {
		res[i] = ival.Index(i).Interface()
	}

	return res
}
