package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateEnumDefinitions(file *jen.File) error {
	enumTypes := make([]string, len(g.schema.Enums))
	enumIndex := 0
	for _type := range g.schema.Enums {
		enumTypes[enumIndex] = _type
		enumIndex++
	}

	sort.Strings(enumTypes)

	for _, enumType := range enumTypes {
		values := g.schema.Enums[enumType]
		sort.Slice(values, func(i, j int) bool {
			return values[i].Name < values[j].Name
		})

		file.Add(g.generateSpecificEnum(enumType, values)...)
	}
	return nil
}

func (g *Generator) generateSpecificEnum(enumType string, enumValues []enum) []jen.Code {
	total := make([]jen.Code, 0)

	typeID := noramlizeIdentificator(enumType)

	enumDef := jen.Type().Id(typeID).Uint32()
	total = append(total, enumDef, jen.Line())

	opc := make([]jen.Code, len(enumValues))
	cases := make([]jen.Code, len(enumValues))
	for i, id := range enumValues {
		name := noramlizeIdentificator(id.Name)

		opc[i] = jen.Id(name).Id(typeID).Op("=").Lit(int(id.CRC))
		cases[i] = jen.Case(jen.Id(typeID).Call(jen.Lit(int(id.CRC)))).Block(jen.Return(jen.Lit(id.Name)))
	}

	total = append(total, jen.Const().Defs(opc...), jen.Line())

	cases = append(cases, jen.Default().Block(jen.Return(jen.Lit("<UNKNOWN "+enumType+">"))))

	stringFunc := jen.Func().Params(jen.Id("e").Id(typeID)).Id("String").Params().String().Block(
		jen.Switch(jen.Id("e")).Block(cases...),
	)

	crcFunc := jen.Func().Params(jen.Id("e").Id(typeID)).Id("CRC").Params().Uint32().Block(
		jen.Return(jen.Uint32().Call(jen.Id("e"))),
	)

	encoderFunc := jen.Func().Params(jen.Id("e").Id(typeID)).Id("Encode").Params().Index().Byte().Block(
		jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto/serialize", "NewEncoder").Call(),
		jen.Id("buf.PutCRC").Call(jen.Uint32().Call(jen.Id("e"))),
		jen.Line(),
		jen.Return(jen.Id("buf.Result").Call()),
	)

	total = append(total,
		stringFunc,
		jen.Line(),
		jen.Line(),
		crcFunc,
		jen.Line(),
		jen.Line(),
		encoderFunc,
		jen.Line(),
		jen.Line(),
	)

	return total

	// старые комменты, не ебу
	/*
		!type InputWallPaper uint32
		!
		!const (
		!	InputWallPaperNoFile InputWallPaper = 2217196460
		!)
		!
		!func (e *InputWallPaper) String() string {
		!	switch e {
		!	case 2217196460:
		!		return "inputWallPaperNoFile"
		!	default:
		!		return "<UNKNOWN input.WallPaper>"
		!	}
		!}
		!
		!func (e *InputWallPaper) Encode() []byte {
		!	buf := mtproto.NewEncoder()
		!	buf.PutUint(uint32(e))
		!
		!	return buf.Result()
		!}
	*/
}
