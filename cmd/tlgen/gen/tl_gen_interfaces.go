package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateInterfaces(f *jen.File) error {
	keys := make([]string, 0, len(g.schema.Types))
	for key := range g.schema.Types {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, i := range keys {
		structs := g.schema.Types[i]

		iface := jen.Type().Id(noramlizeIdentificator(i)).Interface(
			jen.Qual("github.com/xelaj/mtproto/serialize", "TL"),
			jen.Id("Implements"+noramlizeIdentificator(i)).Params(),
			jen.Id("Validate").Params().Id("error"),
		)
		f.Add(
			iface,
			jen.Line(),
		)

		for _, _struct := range structs {
			str, err := g.generateStruct(_struct)
			if err != nil {
				return err
			}
			structName := noramlizeIdentificator(_struct.Name)

			crcFunc := jen.Func().Params(jen.Id("*" + structName)).Id("CRC").Params().Uint32().Block(
				jen.Return(jen.Lit(_struct.CRC)),
			)

			implFunc := jen.Func().Params(jen.Id("*" + structName)).Id("Implements" + noramlizeIdentificator(i)).Params().Block()

			encoderFunc, err := g.generateEncodeFunc(_struct)
			if err != nil {
				return err
			}

			encoderNonReflectFunc, err := g.generateEncodeNonreflectFunc(_struct)
			if err != nil {
				return err
			}

			validateFn, err := g.generateStructValidatorFunc(_struct)
			if err != nil {
				return err
			}

			f.Add(
				str,
				jen.Line(),
				jen.Line(),
				validateFn,
				jen.Line(),
				jen.Line(),
				crcFunc,
				jen.Line(),
				jen.Line(),
				implFunc,
				jen.Line(),
				jen.Line(),
				encoderFunc,
				jen.Line(),
				jen.Line(),
				encoderNonReflectFunc,
				jen.Line(),
				jen.Line(),
			)
		}
	}

	return nil
}
