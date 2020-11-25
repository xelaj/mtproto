package main

import (
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
)

func GenerateConstructorRouter(file *jen.File, data *FileStructure) error {
	structs, enums := data.GetAllConstructors()

	sortedCrcs := make([]int, 0)
	for crc := range structs {
		sortedCrcs = append(sortedCrcs, int(crc))
	}
	for crc := range enums {
		sortedCrcs = append(sortedCrcs, int(crc))
	}
	sort.Ints(sortedCrcs)

	cases := make([]jen.Code, 0)
	for _, crc := range sortedCrcs {
		var obj jen.Code
		var isEnum jen.Code
		if id, ok := structs[uint32(crc)]; ok {
			obj = jen.Op("&").Id(id).Values()
			isEnum = jen.False()
		} else if id, ok := enums[uint32(crc)]; ok {
			obj = jen.Id(id)
			isEnum = jen.True()
		} else {
			panic(fmt.Sprintf("where did you find that crc?? %d", crc))
		}

		cases = append(cases, jen.Case(jen.Lit(uint32(crc))).Block(jen.Return(obj, isEnum, jen.Nil())))
	}

	cases = append(cases, jen.Default().Block(
		jen.Return(jen.Nil(), jen.False(), jen.Qual("github.com/xelaj/errs", "NotFound").Call(jen.Lit("constructorID"), jen.Id("constructorID"))),
	))

	f := jen.Func().Id("GenerateStructByConstructor").Params(jen.Id("constructorID").Uint32()).
		Params(
			jen.Id("object").Qual("github.com/xelaj/mtproto/serialize", "TL"),
			jen.Id("isEnum").Bool(),
			jen.Id("err").Error(),
		).
		Block(
			jen.Switch(jen.Id("constructorID")).Block(cases...),
		)

	file.Add(f)

	return nil
}
