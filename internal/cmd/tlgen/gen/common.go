package gen

import (
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

func (g *Generator) generateStruct(str tlparser.Object) *jen.Statement {
	fields := make([]jen.Code, 0, len(str.Parameters))

	for _, field := range str.Parameters {
		name := strcase.ToCamel(field.Name)
		typ := field.Type
		valueInsideFlag := false

		if name == "Flags" && typ == "bitflags" {
			continue
		}

		f := jen.Id(name)
		if field.IsVector {
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
			valueInsideFlag = true
		default:
			if _, ok := g.schema.Enums[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := g.schema.Types[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := g.schema.SingleInterfaceCanonical[typ]; ok {
				f = f.Id("*" + g.goify(typ))
				break
			}

			//pp.Fprintln(os.Stderr, g.schema)
			panic("пробовали обработать '" + field.Type + "'")
		}

		tags := map[string]string{}
		if field.IsOptional {
			tags["tl"] = "flag:" + strconv.Itoa(field.BitToTrigger)
			if valueInsideFlag {
				tags["tl"] += ",encoded_in_bitflags"
			}
		}

		f.Tag(tags)
		fields = append(fields, f)
	}

	structName := g.goify(str.Name)

	return jen.Type().Id(structName).Struct(
		fields...,
	)
}

func (g *Generator) generateMethodCallerFunc(method tlparser.Method) *jen.Statement {
	resp := createParamsStructFromMethod(method)
	maximumPositionalArguments := 0
	if haveOptionalParams(resp.Parameters) {
		maximumPositionalArguments++
	}

	funcParameters := make([]jen.Code, 0)
	methodName := g.goify(method.Name)
	typeName := methodName + "Params"

	argsAsSingleItem := false
	if len(resp.Parameters) > maximumPositionalArguments {
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

	assertedType := g.goify(method.Response.Type)
	//firstErrorReturn := jen.Code(jen.Nil())
	if assertedType == "Bool" {
		assertedType = "bool"
		//firstErrorReturn = jen.False()
	}
	if assertedType == "Long" {
		assertedType = "int64"
		//firstErrorReturn = jen.Lit(0)
	}
	if assertedType == "Int" {
		assertedType = "int"
		//firstErrorReturn = jen.Lit(0)
	}

	if method.Response.IsList {
		assertedType = "[]" + assertedType
	}

	return jen.Func().Params(jen.Id("c").Id("*Client")).Id(methodName).Params(funcParameters...).Params(jen.Id(assertedType), jen.Error()).Block(
		jen.Var().Id("resp").Id(assertedType),
		jen.Err().Op(":=").Id("c.MakeRequest").Call(requestStruct, jen.Id("&resp")),
		jen.Return(jen.Id("resp"), jen.Err()),
	)
}

func createCrcFunc(typ string, crc uint32) *jen.Statement {
	hex := fmt.Sprintf("0x%x", crc)
	return jen.Func().Params(jen.Id(typ)).Id("CRC").Params().Uint32().
		Id("{" + jen.Return(jen.Id(hex)).GoString() + "}")
}
