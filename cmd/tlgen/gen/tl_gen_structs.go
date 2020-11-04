package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/xelaj/mtproto/cmd/tlgen/typelang"
)

func (g *Generator) generateSpecificStructs(f *jen.File) error {
	sigKeys := make([]string, 0, len(g.schema.SingleInterfaceCanonical))
	for key := range g.schema.SingleInterfaceCanonical {
		sigKeys = append(sigKeys, key)
	}
	sort.Strings(sigKeys)

	for _, _type := range g.schema.SingleInterfaceTypes {
		interfaceName := ""
		for _, k := range sigKeys {
			v := g.schema.SingleInterfaceCanonical[k]
			if v == _type.Name {
				interfaceName = k
			}
		}

		if interfaceName == "" {
			panic("не нашли каноничное имя")
		}

		_structWithIfaceName := typelang.Object{
			Name:       interfaceName,
			CRC:        _type.CRC,
			Parameters: _type.Parameters,
		}

		str, err := g.generateStruct(_structWithIfaceName)
		if err != nil {
			return err
		}

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + noramlizeIdentificator(interfaceName))).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(_structWithIfaceName.CRC)),
		)

		validatorFunc, err := g.generateStructValidatorFunc(_structWithIfaceName)
		if err != nil {
			return err
		}

		encoderFunc, err := g.generateEncodeFunc(_structWithIfaceName)
		if err != nil {
			return err
		}

		encoderNonreflectFunc, err := g.generateEncodeNonreflectFunc(_structWithIfaceName)
		if err != nil {
			return err
		}

		f.Add(
			// jen.Commentf("interface name: %s", interfaceName),
			// jen.Line(),
			// jen.Commentf("original  name: %s", _type.Name),
			// jen.Line(),
			str,
			jen.Line(),
			jen.Line(),
			validatorFunc,
			jen.Line(),
			jen.Line(),
			crcFunc,
			jen.Line(),
			jen.Line(),
			encoderFunc,
			jen.Line(),
			jen.Line(),
			encoderNonreflectFunc,
			jen.Line(),
			jen.Line(),
		)
	}

	return nil
}
