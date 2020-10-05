package main

import (
	"os"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/k0kubun/pp"
)

func GenerateInterfaces(file *jen.File, data *FileStructure) error {
	for i, structs := range data.Types {
		t := jen.Type().Id(normalizeID(i, true)).Interface(
			jen.Qual("github.com/xelaj/mtproto/serialize", "TL"),
			jen.Id("Implements"+normalizeID(i, true)).Params(),
		)
		file.Add(t)
		file.Line()

		for _, _struct := range structs {
			fields := make([]jen.Code, len(_struct.Fields))
			atLeastOneFieldOptional := false
			maxFlagBit := 0
			putFuncs := make([]jen.Code, len(_struct.Fields))

			for i, field := range _struct.Fields {
				name := strcase.ToCamel(field.Name)
				typ := field.Type
				ЗНАЧЕНИЕ_В_ФЛАГЕ := false

				if name == "Flags" && typ == "bitflags" {
					name = "__flagsPosition"
				}

				f := jen.Id(name)
				putFuncId := ""
				if field.IsList {
					f = f.Index()
				}

				switch typ {
				case "Bool":
					f = f.Bool()
					putFuncId = "buf.PutBool"
				case "long":
					f = f.Int64()
					putFuncId = "buf.PutLong"
				case "double":
					f = f.Float64()
					putFuncId = "buf.PutDouble"
				case "int":
					f = f.Int32()
					putFuncId = "buf.PutInt"
				case "string":
					f = f.String()
					putFuncId = "buf.PutString"
				case "bytes":
					f = f.Index().Byte()
					putFuncId = "buf.PutMessage"
				case "bitflags":
					f = f.Struct().Comment("flags param position")
				case "true":
					f = f.Bool()
					//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
					//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
					ЗНАЧЕНИЕ_В_ФЛАГЕ = true
				default:
					putFuncId = "buf.PutRawBytes"

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

				if field.IsList {
					putFuncId = "buf.PutVector"
				}

				putFunc := jen.Null()
				if putFuncId == "buf.PutRawBytes" {
					putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name).Dot("Encode").Call())
				} else if putFuncId != "" {
					putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name))
				}

				tags := map[string]string{}
				if !field.IsOptional {
					tags["validate"] = "required"
				} else {
					tags["flag"] = strconv.Itoa(field.BitToTrigger)
					if ЗНАЧЕНИЕ_В_ФЛАГЕ {
						tags["flag"] += ",encoded_in_bitflags"
					}
					atLeastOneFieldOptional = true
					if maxFlagBit < field.BitToTrigger {
						maxFlagBit = field.BitToTrigger
					}

					putFunc = jen.If(jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name)))).Block(
						putFunc,
					)
				}

				f.Tag(tags)

				fields[i] = f
				putFuncs[i] = putFunc
			}

			structName := normalizeID(_struct.Name, false)

			t := jen.Type().Id(structName).Struct(
				fields...,
			)
			file.Add(t)
			file.Add(jen.Line())

			// CRC() uint23
			file.Add(jen.Func().Params(jen.Id("*" + structName)).Id("CRC").Params().Uint32().Block(
				jen.Return(jen.Lit(_struct.CRCCode)),
			))
			file.Line()

			// Implements<InterfaceName>()
			file.Add(jen.Func().Params(jen.Id("*" + structName)).Id("Implements" + normalizeID(i, true)).Params().Uint32().Block())
			file.Line()

			// Ecncode() []byte
			calls := make([]jen.Code, 0)
			calls = append(calls,
				jen.Id("err").Op(":=").Qual("github.com/go-playground/validator", "New").Call().Dot("Struct").Call(jen.Id("e")),
				jen.Qual("github.com/xelaj/go-dry", "PanicIfErr").Call(jen.Id("err")),
				jen.Line(),
			)

			if atLeastOneFieldOptional {
				// string это fieldname
				sortedOptionalValues := make([][]*Param, maxFlagBit+1)
				for _, field := range _struct.Fields {
					if !field.IsOptional {
						continue
					}
					if sortedOptionalValues[field.BitToTrigger] == nil {
						sortedOptionalValues[field.BitToTrigger] = make([]*Param, 0)
					}

					sortedOptionalValues[field.BitToTrigger] = append(sortedOptionalValues[field.BitToTrigger], &Param{
						Name: field.Name,
						Type: field.Type,
					})
				}

				flagchecks := make([]jen.Code, len(sortedOptionalValues))
				for i, fields := range sortedOptionalValues {
					if len(fields) == 0 {
						continue
					}

					statements := jen.Null()
					for j, field := range fields {
						if j != 0 {
							statements.Add(jen.Op("||"))
						}
						//? zero.IsZeroVal(e.Fieldname)
						statements.Add(jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name))))
					} //?               if !zero.IsZeroVal(n) || !zer.IsZeroVal(m)...
					flagchecks[i] = jen.If(statements).Block(
						//? flag |= 1 << n
						jen.Id("flag").Op("|=").Lit(1).Op("<<").Lit(i),
					)
				}

				calls = append(calls, jen.Var().Id("flag").Uint32())
				calls = append(calls,
					flagchecks...,
				)
			}

			calls = append(calls,
				jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto", "NewEncoder").Call(),
				jen.Id("buf.PutUint").Call(jen.Id("e.CRC").Call()),
			)

			calls = append(calls,
				putFuncs...,
			)

			calls = append(calls,
				jen.Return(jen.Id("buf.Result").Call()),
			)

			f := jen.Func().Params(jen.Id("e").Id("*" + structName)).Id("Encode").Params().Index().Byte().Block(
				calls...,
			)

			file.Add(f)
			file.Add(jen.Line())

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

			file.Add(jen.Line())
		}
	}

	return nil
}
