package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateEnumDefinitions(file *jen.File) {
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
}

func (g *Generator) generateSpecificEnum(enumType string, enumValues []enum) []jen.Code {
	total := make([]jen.Code, 0)

	typeID := g.goify(enumType)

	enumDef := jen.Type().Id(typeID).Uint32()
	total = append(total, enumDef, jen.Line())

	opc := make([]jen.Code, len(enumValues))
	cases := make([]jen.Code, len(enumValues))
	for i, id := range enumValues {
		name := g.goify(id.Name)

		opc[i] = jen.Id(name).Id(typeID).Op("=").Lit(int(id.CRC))
		cases[i] = jen.Case(jen.Id(typeID).Call(jen.Lit(int(id.CRC)))).Block(jen.Return(jen.Lit(id.Name)))
	}

	total = append(total, jen.Const().Defs(opc...), jen.Line())

	cases = append(cases, jen.Default().Block(jen.Return(jen.Lit("<UNKNOWN "+enumType+">"))))

	stringFunc := jen.Func().Params(jen.Id("e").Id(typeID)).Id("String").Params().String().Block(
		jen.Switch(jen.Id("e")).Block(cases...),
	)

	crcFunc := jen.Func().Params(jen.Id("e").Id(typeID)).Id("CRC").Params().Uint32().
		Id("{ return uint32(e) }")

	total = append(total,
		stringFunc,
		jen.Line(),
		jen.Line(),
		crcFunc,
		jen.Line(),
		jen.Line(),
	)

	return total
}
