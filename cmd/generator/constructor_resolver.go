package main

import "github.com/dave/jennifer/jen"

func GenerateConstructorRouter(file *jen.File, data *FileStructure) error {
	structs, enums := data.GetAllConstructors()

	cases := make([]jen.Code, 0)
	for constructorID, id := range structs {
		cases = append(cases, jen.Case(jen.Lit(constructorID)).Block(jen.Return(jen.Op("&").Id(id).Values(), jen.False())))
	}
	for constructorID, id := range enums {
		cases = append(cases, jen.Case(jen.Lit(constructorID)).Block(jen.Return(jen.Id(id), jen.True())))
	}

	cases = append(cases, jen.Default().Block(jen.Panic(jen.Lit("unknown constructor id: ").Op("+").Qual("fmt", "Sprintf").Call(jen.Lit("%#v"), jen.Id("constructorID")))))

	f := jen.Func().Id("GenerateStructByConstructor").Params(jen.Id("constructorID").Uint32()).Params(jen.Id("object").Qual("github.com/xelaj/mtproto", "TLNEW"), jen.Id("isEnum").Bool()).Block(
		jen.Switch(jen.Id("constructorID")).Block(cases...),
	)

	file.Add(f)

	return nil
}
