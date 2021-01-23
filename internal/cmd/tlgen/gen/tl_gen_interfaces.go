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
		f.Add(jen.Type().Id(goify(i, true)).Interface(
			jen.Qual(tlPackagePath, "Object"),
			jen.Id("Implements"+goify(i, true)).Params(),
		))

		structs := g.schema.Types[i]

		sort.Slice(structs, func(i, j int) bool {
			return structs[i].Name < structs[j].Name
		})

		for _, _type := range structs {
			if goify(_type.Name, true) == goify(i, true) {
				_type.Name += "Obj"
			}

			f.Add(g.generateStructTypeAndMethods(_type, []string{goify(i, true)}))
			f.Line()
		}
	}
}
