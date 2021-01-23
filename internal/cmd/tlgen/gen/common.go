package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	
	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

func (g *Generator) generateMethodCallerFunc(method tlparser.Method) *jen.Statement {
	resp := createParamsStructFromMethod(method)
	maximumPositionalArguments := 0
	if haveOptionalParams(resp.Parameters) {
		maximumPositionalArguments++
	}

	funcParameters := make([]jen.Code, 0)
	methodName := goify(method.Name, true)
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

	assertedType := goify(method.Response.Type, true)
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
