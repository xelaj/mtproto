package main

import "github.com/dave/jennifer/jen"

func generateMethodCallerFunc(method *FuncObject, data *FileStructure) (*jen.Statement, error) {
	resp := method.ParamsStruct()
	maximumPositionalArguments := 0
	if resp.HaveOptionalArgs() {
		maximumPositionalArguments++
	}

	funcParameters := make([]jen.Code, 0)
	methodName := normalizeID(method.Name, false)
	typeName := methodName + "Params"

	argsAsSingleItem := false
	if len(resp.Fields) > maximumPositionalArguments {
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
		//firstErrorReturn = jen.False()
	}
	if assertedType == "Long" {
		assertedType = "*serialize.Long"
		//firstErrorReturn = jen.Lit(0)
	}
	if assertedType == "Int" {
		assertedType = "*serialize.Int"
		//firstErrorReturn = jen.Lit(0)
	}

	calls := make([]jen.Code, 0)
	calls = append(calls,
		jen.List(jen.Id("data"), jen.Err()).Op(":=").Id("c.MakeRequest").Call(requestStruct),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			//jen.Return(firstErrorReturn, jen.Qual("github.com/pkg/errors", "Wrap").Call(jen.Err(), jen.Lit("sedning "+methodName))),
			jen.Return(
				firstErrorReturn,
				jen.Qual("fmt", "Errorf").Call(jen.Lit(methodName+": %w"), jen.Id("err")),
			),
		),
		jen.Line(),
		jen.List(jen.Id("resp"), jen.Id("ok")).Op(":=").Id("data").Assert(jen.Id(assertedType)),
		jen.If(jen.Op("!").Id("ok")).Block(
			//jen.Panic(jen.Lit("got invalid response type: ").Op("+").Qual("reflect", "TypeOf").Call(jen.Id("data")).Dot("String").Call()),
			jen.Err().Op(":=").Qual("fmt", "Errorf").Call(jen.Lit(methodName+": got invalid response type: %T"), jen.Id("data")),
			jen.Comment(
				jen.Return(
					firstErrorReturn,
					jen.Qual("fmt", "Errorf").Call(jen.Lit(methodName+": got invalid response type: %T"), jen.Id("data")),
				).GoString(),
			),
			jen.Panic(jen.Err()),
		),
		jen.Line(),
		jen.Return(jen.Id("resp"), jen.Nil()),
	)

	return jen.Func().Params(jen.Id("c").Id("*Client")).Id(methodName).Params(funcParameters...).Params(jen.Id(assertedType), jen.Error()).Block(
		calls...,
	), nil
}
