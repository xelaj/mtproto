package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

func (g *Generator) generateSpecificStructs(f *jen.File, data *internalSchema) error {
	sort.Slice(data.SingleInterfaceTypes, func(i, j int) bool {
		return data.SingleInterfaceTypes[i].Name < data.SingleInterfaceTypes[j].Name
	})

	sigKeys := make([]string, 0, len(data.SingleInterfaceCanonical))
	for key := range data.SingleInterfaceCanonical {
		sigKeys = append(sigKeys, key)

	}
	sort.Strings(sigKeys)

	// for _, key := range sigKeys {
	// 	fmt.Println("data.SingleInterfaceCanonical[interface]: ", key)
	// }

	for _, _type := range data.SingleInterfaceTypes {
		interfaceName := ""
		for _, k := range sigKeys {
			v := data.SingleInterfaceCanonical[k]
			if v == _type.Name {
				interfaceName = k
			}
		}

		if interfaceName == "" {
			panic("не нашли каноничное имя")
		}

		_structWithIfaceName := tlparser.Object{
			Name:       interfaceName,
			CRC:        _type.CRC,
			Parameters: _type.Parameters,
		}

		str, err := g.generateStruct(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + g.goify(interfaceName))).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(_structWithIfaceName.CRC)),
		)

		//fmt.Println("SingleInterfaceTypes_name[struct]:", _structWithIfaceName.Name)
		validatorFunc, err := g.generateStructValidatorFunc(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		encoderFunc, err := g.generateEncodeFunc(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		encoderNonreflectFunc, err := g.generateEncodeNonreflectFunc(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		_ = crcFunc
		_ = encoderFunc
		_ = encoderNonreflectFunc
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
