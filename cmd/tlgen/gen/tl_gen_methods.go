package gen

import (
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateMethods(f *jen.File, data *internalSchema) error {
	sort.Slice(data.Methods, func(i, j int) bool {
		return data.Methods[i].Name < data.Methods[j].Name
	})

	for _, method := range data.Methods {
		paramsStruct, err := g.generateStruct(createParamsStructFromMethod(method), data)
		if err != nil {
			return fmt.Errorf("generate method params struct: %s: %w", method.Name, err)
		}

		methodName := g.goify(method.Name)
		typeName := methodName + "Params"

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(method.CRC)),
		)

		validatorFunc, err := g.generateStructValidatorFunc(createParamsStructFromMethod(method), data)
		if err != nil {
			return fmt.Errorf("generate validator func for method params %s: %w", method.Name, err)
		}

		encodeFunc, err := g.generateEncodeFunc(createParamsStructFromMethod(method), data)
		if err != nil {
			return fmt.Errorf("generate encode func for method %s: %w", method.Name, err)
		}

		encodeNonreflectFunc, err := g.generateEncodeNonreflectFunc(createParamsStructFromMethod(method), data)
		if err != nil {
			return fmt.Errorf("generate encode func for method %s: %w", method.Name, err)
		}

		callerFunc, err := g.generateMethodCallerFunc(method, data)
		if err != nil {
			return fmt.Errorf("generate caller func for method %s: %w", method.Name, err)
		}
		_ = callerFunc
		_ = encodeFunc
		_ = encodeNonreflectFunc
		_ = crcFunc
		f.Add(
			paramsStruct,
			jen.Line(),
			jen.Line(),
			crcFunc,
			jen.Line(),
			jen.Line(),
			validatorFunc,
			jen.Line(),
			jen.Line(),
			encodeFunc,
			jen.Line(),
			jen.Line(),
			encodeNonreflectFunc,
			jen.Line(),
			jen.Line(),
			callerFunc,
			jen.Line(),
		)
	}

	return nil
}
