package main

import (
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
)

func GenerateMethods(f *jen.File, data *FileStructure) error {
	sort.Slice(data.Methods, func(i, j int) bool {
		return data.Methods[i].Name < data.Methods[j].Name
	})

	for _, method := range data.Methods {
		paramsStruct, err := generateStruct(method.ParamsStruct(), data)
		if err != nil {
			return fmt.Errorf("generate method params struct: %s: %w", method.Name, err)
		}

		methodName := normalizeID(method.Name, false)
		typeName := methodName + "Params"

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(method.CRCCode)),
		)

		encodeFunc, err := generateEncodeFunc(method.ParamsStruct(), data)
		if err != nil {
			return fmt.Errorf("generate encode func for method %s: %w", method.Name, err)
		}

		callerFunc, err := generateMethodCallerFunc(method, data)
		if err != nil {
			return fmt.Errorf("generate caller func for method %s: %w", method.Name, err)
		}

		f.Add(
			paramsStruct,
			jen.Line(),
			jen.Line(),
			crcFunc,
			jen.Line(),
			jen.Line(),
			encodeFunc,
			jen.Line(),
			jen.Line(),
			callerFunc,
			jen.Line(),
		)
	}

	return nil
}
