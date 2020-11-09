package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

func (g *Generator) generateSpecificStructs(f *jen.File) {
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

		_structWithIfaceName := tlparser.Object{
			Name:       interfaceName,
			CRC:        _type.CRC,
			Parameters: _type.Parameters,
		}

		f.Add(
			// jen.Commentf("interface name: %s", interfaceName),
			// jen.Line(),
			// jen.Commentf("original  name: %s", _type.Name),
			// jen.Line(),
			g.generateStruct(_structWithIfaceName),
			jen.Line(),
			jen.Line(),
			createCrcFunc("*"+g.goify(interfaceName), _structWithIfaceName.CRC),
			jen.Line(),
			jen.Line(),
		)
	}
}
