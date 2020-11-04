package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateMethods(f *jen.File) error {
	for _, method := range g.schema.Methods {
		params := createParamsStructFromMethod(&method)

		paramsStruct, err := g.generateStruct(params)
		if err != nil {
			return fmt.Errorf("generate method params struct: %s: %w", method.Name, err)
		}

		crcFunc := jen.Func().Params(jen.Id("e").Id("*" + noramlizeIdentificator(params.Name))).Id("CRC").Params().Uint32().Block(
			jen.Return(jen.Lit(method.CRC)),
		)

		validatorFunc, err := g.generateStructValidatorFunc(params)
		if err != nil {
			return fmt.Errorf("generate validator func for method params %s: %w", method.Name, err)
		}

		encoderFunc, err := g.generateEncodeFunc(params)
		if err != nil {
			return fmt.Errorf("generate encode func for method %s: %w", method.Name, err)
		}

		encoderNonreflectFunc, err := g.generateEncodeNonreflectFunc(params)
		if err != nil {
			return fmt.Errorf("generate encode func for method %s: %w", method.Name, err)
		}

		callerFunc, err := g.generateMethodCallerFunc(&method)
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
			validatorFunc,
			jen.Line(),
			jen.Line(),
			encoderFunc,
			jen.Line(),
			jen.Line(),
			encoderNonreflectFunc,
			jen.Line(),
			jen.Line(),
			callerFunc,
			jen.Line(),
		)
	}

	return nil
}
