package main

import (
	"os"
	"sort"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/k0kubun/pp"
)

func generateStruct(str *StructObject, data *FileStructure) (*jen.Statement, error) {
	fields := make([]jen.Code, 0, len(str.Fields))

	sort.Slice(str.Fields, func(i, j int) bool {
		return str.Fields[i].Name < str.Fields[j].Name
	})

	for _, field := range str.Fields {
		name := strcase.ToCamel(field.Name)
		typ := field.Type
		ЗНАЧЕНИЕ_В_ФЛАГЕ := false

		if name == "Flags" && typ == "bitflags" {
			name = "__flagsPosition"
		}

		f := jen.Id(name)
		if field.IsList {
			f = f.Index()
		}

		switch typ {
		case "Bool":
			f = f.Bool()
		case "long":
			f = f.Int64()
		case "double":
			f = f.Float64()
		case "int":
			f = f.Int32()
		case "string":
			f = f.String()
		case "bytes":
			f = f.Index().Byte()
		case "bitflags":
			f = f.Struct().Comment("flags param position")
		case "true":
			f = f.Bool()
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			ЗНАЧЕНИЕ_В_ФЛАГЕ = true
		default:
			if _, ok := data.Enums[typ]; ok {
				f = f.Id(normalizeID(typ, false))
				break
			}
			if _, ok := data.Types[typ]; ok {
				f = f.Id(normalizeID(typ, false))
				break
			}
			if _, ok := data.SingleInterfaceCanonical[typ]; ok {
				f = f.Id("*" + normalizeID(typ, false))
				break
			}

			pp.Fprintln(os.Stderr, data)
			panic("пробовали обработать '" + field.Type + "'")
		}

		tags := map[string]string{}
		if !field.IsOptional {
			tags["validate"] = "required"
		} else {
			tags["flag"] = strconv.Itoa(field.BitToTrigger)
			if ЗНАЧЕНИЕ_В_ФЛАГЕ {
				tags["flag"] += ",encoded_in_bitflags"
			}
		}

		f.Tag(tags)
		fields = append(fields, f)
	}

	structName := normalizeID(str.Name, false)

	return jen.Type().Id(structName).Struct(
		fields...,
	), nil
}
