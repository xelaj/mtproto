package gen

import (
	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateMethods(f *jen.File) {
	for _, method := range g.schema.Methods {
		params := createParamsStructFromMethod(method)

		f.Add(
			g.generateStruct(params),
			jen.Line(),
			jen.Line(),
			createCrcFunc("*"+g.goify(params.Name), method.CRC),
			jen.Line(),
			jen.Line(),
			g.generateMethodCallerFunc(method),
			jen.Line(),
		)
	}
}
