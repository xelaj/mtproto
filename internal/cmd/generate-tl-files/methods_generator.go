package main

import (
	"os"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/k0kubun/pp"
)

func GenerateMethods(file *jen.File, data *FileStructure) error {
	for _, method := range data.Methods {
		fields := make([]jen.Code, len(method.Arguments))
		atLeastOneFieldOptional := false
		maxFlagBit := 0
		putFuncs := make([]jen.Code, len(method.Arguments))
		funcParameters := make([]jen.Code, 0)
		for i, field := range method.Arguments {
			name := strcase.ToCamel(field.Name)
			typ := field.Type
			ЗНАЧЕНИЕ_В_ФЛАГЕ := false

			if name == "Flags" && typ == "bitflags" {
				name = "__flagsPosition"
			} else {
				// для вызова метода
				funcParameters = append(funcParameters, jen.Id(name))
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
			switch {
			case name == "__flagsPosition":
				putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
			case ЗНАЧЕНИЕ_В_ФЛАГЕ:
				// не делаем ничего, значение уже заложили в флаг
			case putFuncId == "buf.PutRawBytes":
				putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name).Dot("Encode").Call())
			case putFuncId != "":
				putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name))
			default:
				panic("putFincID is empty!")
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
				if !ЗНАЧЕНИЕ_В_ФЛАГЕ {
					putFunc = jen.If(jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name)))).Block(
						putFunc,
					)
				}
			}

			f.Tag(tags)

			fields[i] = f
			putFuncs[i] = putFunc
			//arg := jen.Id(name)
			//if i+1 < len(method.Arguments) && method.Arguments[i+1].Type == field.Type {
			//	// type will be described next
			//} else {
			//	arg = f // описание энивей одинаковое
			//}
			//
			//funcParameters = append(funcParameters, arg)
		}

		methodName := normalizeID(method.Name, false)
		typeName := methodName + "Params"
		t := jen.Type().Id(typeName).Struct(
			fields...,
		)
		file.Add(t)
		file.Add(jen.Line())

		// CRC() uint32
		file.Add(jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(method.CRCCode)),
		))
		file.Add(jen.Line())

		// Ecncode() []byte
		calls := make([]jen.Code, 0)
		if len(method.Arguments) > 0 {
			calls = append(calls,
				jen.Id("err").Op(":=").Qual("github.com/go-playground/validator", "New").Call().Dot("Struct").Call(jen.Id("e")),
				jen.Qual("github.com/xelaj/go-dry", "PanicIfErr").Call(jen.Id("err")),
				jen.Line(),
			)
		}

		if atLeastOneFieldOptional {
			// string это fieldname
			sortedOptionalValues := make([][]*Param, maxFlagBit+1)
			for _, field := range method.Arguments {
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
			jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto/serialize", "NewEncoder").Call(),
			jen.Id("buf.PutUint").Call(jen.Id("e.CRC").Call()),
		)

		calls = append(calls,
			putFuncs...,
		)

		calls = append(calls,
			jen.Return(jen.Id("buf.Result").Call()),
		)

		f := jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("Encode").Params().Index().Byte().Block(
			calls...,
		)

		file.Add(f)
		file.Add(jen.Line())

		maximumPositionalArguments := 0
		if atLeastOneFieldOptional {
			maximumPositionalArguments++
		}

		argsAsSingleItem := false
		if len(method.Arguments) > maximumPositionalArguments {
			argsAsSingleItem = true
			funcParameters = []jen.Code{jen.Id("params").Id("*" + typeName)}
		}

		parameters := jen.Dict{}
		for _, arg := range funcParameters {
			parameters[arg] = arg
		}

		requestStruct := jen.Op("&").Id(typeName).Values(parameters)
		if argsAsSingleItem {
			requestStruct = jen.Id("params")
		}

		assertedType := ""
		assertedType = normalizeID(method.Returns.Type, false)

		if _, ok := data.SingleInterfaceCanonical[method.Returns.Type]; ok {
			assertedType = "*" + assertedType
		}
		firstErrorReturn := jen.Code(jen.Nil())
		if assertedType == "Bool" {
			assertedType = "*serialize.Bool"
			// firstErrorReturn = jen.False()
		}
		if assertedType == "Long" {
			assertedType = "*serialize.Long"
			// firstErrorReturn = jen.Lit(0)
		}
		if assertedType == "Int" {
			assertedType = "*serialize.Int"
			// firstErrorReturn = jen.Lit(0)
		}

		calls = make([]jen.Code, 0)
		calls = append(calls,
			jen.List(jen.Id("data"), jen.Err()).Op(":=").Id("c.MakeRequest").Call(requestStruct),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(firstErrorReturn, jen.Qual("github.com/pkg/errors", "Wrap").Call(jen.Err(), jen.Lit("sedning "+methodName))),
			),
			jen.Line(),
			jen.List(jen.Id("resp"), jen.Id("ok")).Op(":=").Id("data").Assert(jen.Id(assertedType)),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Panic(jen.Lit("got invalid response type: ").Op("+").Qual("reflect", "TypeOf").Call(jen.Id("data")).Dot("String").Call()),
			),
			jen.Line(),
			jen.Return(jen.Id("resp"), jen.Nil()),
		)

		f = jen.Func().Params(jen.Id("c").Id("*Client")).Id(methodName).Params(funcParameters...).Params(jen.Id(assertedType), jen.Error()).Block(
			calls...,
		)

		file.Add(f)
		file.Add(jen.Line())

	}

	return nil
}

/* //! example method:
type AuthSendCodeParams struct {
	PhoneNumber string
	ApiID       int
	ApiHash     string
	Settings    *CodeSettings
}

func (_ *AuthSendCodeParams) CRC() uint32 {
	return 0xa677244f
}

func (t *AuthSendCodeParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutUint(t.CRC())
	buf.PutString(t.PhoneNumber)
	buf.PutInt(int32(t.ApiID))
	buf.PutString(t.ApiHash)
	buf.PutRawBytes(t.Settings.Encode())
	return buf.Result()
}

func (c *Client) AuthSendCode(PhoneNumber string, ApiID int, ApiHash string, Settings *CodeSettings) (*AuthSentCode, error) {
	data, err := c.MakeRequest(&AuthSendCodeParams{
		PhoneNumber: PhoneNumber,
		ApiID:       ApiID,
		ApiHash:     ApiHash,
		Settings:    Settings,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending AuthSendCode")
	}

	resp, ok := data.(*AuthSentCode)
	if !ok {
		panic(errors.New("got invalid response type: " + reflect.TypeOf(data).String()))
	}

	return resp, nil

}
*/
