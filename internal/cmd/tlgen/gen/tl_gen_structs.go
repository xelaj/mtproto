package gen

import (
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/k0kubun/pp"
	
	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

func (g *Generator) generateSpecificStructs(f *jen.File) {
	sort.Slice(g.schema.SingleInterfaceTypes, func(i, j int) bool {
		return g.schema.SingleInterfaceTypes[i].Name < g.schema.SingleInterfaceTypes[j].Name
	})

	for _, _type := range g.schema.SingleInterfaceTypes {
		f.Add(g.generateStructTypeAndMethods(_type, nil))
		f.Line()
	}
}

func (g *Generator) generateStructTypeAndMethods(definition tlparser.Object, implementsMethods []string) jen.Code {
	structName := goify(definition.Name, true)
	containsOptionalParameters := false
	var typeDefinition jen.Code

	fields := make([]jen.Code, len(definition.Parameters))
	for i, param := range definition.Parameters {
		if param.IsOptional {
			containsOptionalParameters = true
		}

		if param.Type == "bitflags" {
			continue
		}
		fields[i] = g.generateStructParameter(&param)
	}
	typeDefinition = jen.Type().Id(structName).Struct(fields...)

	// func (*T) CRC() uint32 { return 0x89abcdef }
	var crcFunc jen.Code
	f := jen.Func().Params(jen.Op("*").Id(structName)).Id("CRC").Params().Uint32().Block(
		jen.Return(jen.Id(fmt.Sprintf("%#v", definition.CRC))),
	)
	crcFunc = f

	// can be nil
	var fieldIndexFunc jen.Code
	if containsOptionalParameters {
		flagBitIndex := -1 // -1 это если нет опциональных полей вообще
		if containsOptionalParameters {
			for i, param := range definition.Parameters {
				if param.Name == "flags" && param.Type == "bitflags" {
					flagBitIndex = i
				}
			}
		}
		if flagBitIndex == -1 {
			pp.Println(definition)
			panic("optional bitflag not found!")
		}

		f := jen.Func().Params(jen.Op("*").Id(structName)).Id("FlagIndex").Params().Int().Block(
			jen.Return(jen.Lit(flagBitIndex)),
		)
		fieldIndexFunc = f
	}

	// func (*T) Implements<InterfaceName>() {}
	var implementFuncs []jen.Code
	if implementsMethods != nil {
		implementFuncs = make([]jen.Code, len(implementsMethods))
		for i, suffixName := range implementsMethods {
			f := jen.Func().Params(jen.Op("*").Id(structName)).Id("Implements" + suffixName).Params().Block()
			implementFuncs[i] = f
		}
	}

	result := &jen.Statement{}
	if definition.Comment != "" {
		result = result.Comment(definition.Comment)
		result = result.Line()
	}

	result = result.Add(typeDefinition, jen.Line(), jen.Line())
	result = result.Add(crcFunc, jen.Line(), jen.Line())
	if fieldIndexFunc != nil {
		result = result.Add(fieldIndexFunc, jen.Line(), jen.Line())
	}

	if implementFuncs != nil {
		result = result.Add(implementFuncs...)
	}
	result = result.Add(jen.Line(), jen.Line())

	return result
}

func (g *Generator) generateStructParameter(param *tlparser.Parameter) *jen.Statement {
	goifiedName := goify(param.Name, true)
	tag := "" // то что в `tl:"..."` находится
	f := jen.Id(goifiedName)
	if param.IsVector {
		f = f.Index()
	}

	if param.IsOptional {
		tag = fmt.Sprintf("flag:%v", param.BitToTrigger)
	}

	if param.Type == "true" {
		tag += ",encoded_in_bitflags"
	}

	f = f.Add(g.typeIdFromSchemaType(param.Type))

	if tag != "" {
		f = f.Tag(map[string]string{"tl": tag})
	}

	if param.Comment != "" {
		f = f.Comment(param.Comment)
	}

	return f
}
