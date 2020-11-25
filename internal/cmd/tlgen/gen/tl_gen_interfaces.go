package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateInterfaces(f *jen.File) {
	keys := make([]string, 0, len(g.schema.Types))
	for key := range g.schema.Types {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, i := range keys {
		structs := g.schema.Types[i]

		iface := jen.Type().Id(g.goify(i)).Interface(
			jen.Qual("github.com/xelaj/mtproto/encoding/tl", "Object"),
			jen.Id("Implements"+g.goify(i)).Params(),
		)
		f.Add(
			iface,
			jen.Line(),
		)

		for _, _struct := range structs {
			implFunc := jen.Func().Params(jen.Id("*" + g.goify(_struct.Name))).Id("Implements" + g.goify(i)).Params().Block()

			f.Add(
				g.generateStruct(_struct),
				jen.Line(),
				jen.Line(),
				createCrcFunc("*"+g.goify(_struct.Name), _struct.CRC),
				jen.Line(),
				jen.Line(),
				implFunc,
				jen.Line(),
				jen.Line(),
			)
		}
	}
}
