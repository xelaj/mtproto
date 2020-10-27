package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func GenerateInterfaces(f *jen.File, data *FileStructure) error {
	keys := make([]string, 0, len(data.Types))
	for key := range data.Types {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, i := range keys {
		structs := data.Types[i]

		iface := jen.Type().Id(normalizeID(i, true)).Interface(
			jen.Qual("github.com/xelaj/mtproto/serialize", "TL"),
			jen.Id("Implements"+normalizeID(i, true)).Params(),
			jen.Id("Validate").Params().Id("error"),
		)
		f.Add(
			iface,
			jen.Line(),
		)

		//fmt.Println("data.Types_key[interface{}]: ", i)

		for _, _struct := range structs {
			str, err := generateStruct(_struct, data)
			if err != nil {
				return err
			}
			structName := normalizeID(_struct.Name, false)

			crcFunc := jen.Func().Params(jen.Id("*" + structName)).Id("CRC").Params().Uint32().Block(
				jen.Return(jen.Lit(_struct.CRCCode)),
			)

			implFunc := jen.Func().Params(jen.Id("*" + structName)).Id("Implements" + normalizeID(i, true)).Params().Block()

			encoderFunc, err := generateEncodeFunc(_struct, data)
			if err != nil {
				return err
			}

			encoderNonReflectFunc, err := generateEncodeNonreflectFunc(_struct, data)
			if err != nil {
				return err
			}

			validateFn, err := generateStructValidatorFunc(_struct, data)
			if err != nil {
				return err
			}

			f.Add(
				str,
				jen.Line(),
				jen.Line(),
				validateFn,
				jen.Line(),
				jen.Line(),
				crcFunc,
				jen.Line(),
				jen.Line(),
				implFunc,
				jen.Line(),
				jen.Line(),
				encoderFunc,
				jen.Line(),
				jen.Line(),
				encoderNonReflectFunc,
				jen.Line(),
				jen.Line(),
			)

			// part of DecodeFram(d *serialize.Decoder)
			// don't touch, we try to make it more cool
			//* calls = make([]jen.Code, 0)
			//* calls = append(calls,
			//* 	jen.Id("crc").Op(":=").Id("buf.PopUint").Call(),
			//* 	jen.If(jen.Id("crc").Op("!=").Id("e.CRC").Call()).Block(
			//* 		jen.Panic(jen.Lit("wrong type: ").Op("+").Qual("fmt", "Sprintf").Call(jen.Lit("%#v"), jen.Id("crc"))),
			//* 	),
			//* )
			//*
			//* for _, field := range _struct.Fields {
			//* 	name := strcase.ToCamel(field.Name)
			//* 	typ := field.Type
			//*
			//* 	funcCall := jen.Nil()
			//* 	listType := ""
			//*
			//* 	switch typ {
			//* 	case "true", "Bool":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopBool").Call()
			//* 		listType = "bool"
			//* 	case "long":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopLong").Call()
			//* 		listType = "int64"
			//* 	case "double":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopDouble").Call()
			//* 		listType = "float64"
			//* 	case "int":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopInt").Call()
			//* 		listType = "int32"
			//* 	case "string":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopString").Call()
			//* 		listType = "string"
			//* 	case "bytes":
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopMessage").Call()
			//* 		listType = "[]byte"
			//* 	case "bitflags":
			//* 		funcCall = jen.Id("flags").Op(":=").Id("buf.PopUint").Call()
			//* 		listType = "uint32"
			//* 	default:
			//* 		normalized := normalizeID(typ, false)
			//* 		if _, ok := data.Enums[typ]; ok {
			//* 			//*((buf.PopObj()).(*SecureValueType))
			//* 			funcCall = jen.Id("e." + name).Op("=").Id("*").Call(jen.Id("buf.PopObj").Call().Assert(jen.Id("*" + normalized)))
			//* 			listType = normalized
			//* 			break
			//* 		}
			//* 		if _, ok := data.Types[typ]; ok {
			//* 			funcCall = jen.Id("e." + name).Op("=").Id(normalized).Call(jen.Id("buf.PopObj").Call())
			//* 			listType = normalized
			//* 			break
			//* 		}
			//* 		if _, ok := data.SingleInterfaceCanonical[typ]; ok {
			//* 			funcCall = jen.Id("e." + name).Op("=").Id("buf.PopObj").Call().Assert(jen.Id("*" + normalized))
			//* 			listType = "*" + normalized
			//* 			break
			//* 		}
			//*
			//* 		pp.Fprintln(os.Stderr, data)
			//* 		panic("пробовали обработать '" + field.Type + "'")
			//* 	}
			//*
			//* 	if field.IsList {
			//* 		funcCall = jen.Id("e." + name).Op("=").Id("buf.PopVector").Call(jen.Qual("reflect", "TypeOf").Call(jen.Index().Id(listType).Values())).Assert(jen.Index().Id(listType))
			//* 	}
			//*
			//* 	if field.IsOptional {
			//* 		funcCall = jen.If(jen.Id("flags").Op("&").Lit(1).Op("<<").Lit(field.BitToTrigger).Op(">").Lit(0)).Block(
			//* 			funcCall,
			//* 		)
			//* 	}
			//*
			//* 	calls = append(calls,
			//* 		funcCall,
			//* 	)
			//* }
		}
	}

	return nil
}
