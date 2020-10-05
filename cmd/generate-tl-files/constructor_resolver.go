package main

import "github.com/dave/jennifer/jen"

func GenerateConstructorRouter(file *jen.File, data *FileStructure) error {
	structs, enums := data.GetAllConstructors()

	cases := make([]jen.Code, 0)
	for constructorID, id := range structs {
		cases = append(cases, jen.Case(jen.Lit(constructorID)).Block(jen.Return(jen.Op("&").Id(id).Values(), jen.False(), jen.Nil())))
	}
	for constructorID, id := range enums {
		cases = append(cases, jen.Case(jen.Lit(constructorID)).Block(jen.Return(jen.Id(id), jen.True(), jen.Nil())))
	}

	cases = append(cases, jen.Default().Block(jen.Return(jen.Nil(), jen.False(), jen.Qual("errors", "New").Call(jen.Lit("constructor not found")))))

	f := jen.Func().Id("GenerateStructByConstructor").Params(jen.Id("constructorID").Uint32()).Params(jen.Id("object").Qual("github.com/xelaj/mtproto", "TLNEW"), jen.Id("isEnum").Bool()).Block(
		jen.Switch(jen.Id("constructorID")).Block(cases...),
	)

	file.Add(f)

	return nil
}
