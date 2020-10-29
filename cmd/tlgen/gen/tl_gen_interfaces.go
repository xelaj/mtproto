package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateInterfaces(f *jen.File) error {
	keys := make([]string, 0, len(g.schema.Types))
	for key := range g.schema.Types {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, i := range keys {
		structs := g.schema.Types[i]

		iface := jen.Type().Id(g.goify(i)).Interface(
			jen.Qual("github.com/xelaj/mtproto/serialize", "TL"),
			jen.Id("Implements"+g.goify(i)).Params(),
			jen.Id("Validate").Params().Id("error"),
		)
		f.Add(
			iface,
			jen.Line(),
		)

		for _, _struct := range structs {
			str, err := g.generateStruct(_struct)
			if err != nil {
				return err
			}
			structName := g.goify(_struct.Name)

			crcFunc := jen.Func().Params(jen.Id("*" + structName)).Id("CRC").Params().Uint32().Block(
				jen.Return(jen.Lit(_struct.CRC)),
			)

			implFunc := jen.Func().Params(jen.Id("*" + structName)).Id("Implements" + g.goify(i)).Params().Block()

			encoderFunc, err := g.generateEncodeFunc(_struct)
			if err != nil {
				return err
			}

			encoderNonReflectFunc, err := g.generateEncodeNonreflectFunc(_struct)
			if err != nil {
				return err
			}

			validateFn, err := g.generateStructValidatorFunc(_struct)
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

			// старые комменты, не ебу
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
