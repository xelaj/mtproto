package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func GenerateSpecificStructs(f *jen.File, data *FileStructure) error {
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

		_structWithIfaceName := new(StructObject)
		_structWithIfaceName.Name = interfaceName
		_structWithIfaceName.CRCCode = _type.CRCCode
		_structWithIfaceName.Fields = _type.Fields

		str, err := generateStruct(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + normalizeID(interfaceName, false))).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(_structWithIfaceName.CRCCode)),
		)

		//fmt.Println("SingleInterfaceTypes_name[struct]:", _structWithIfaceName.Name)
		validatorFunc, err := generateStructValidatorFunc(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		encoderFunc, err := generateEncodeFunc(_structWithIfaceName, data)
		if err != nil {
			return err
		}

		encoderNonreflectFunc, err := generateEncodeNonreflectFunc(_structWithIfaceName, data)
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
