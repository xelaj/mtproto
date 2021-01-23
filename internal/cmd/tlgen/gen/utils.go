package gen

import (
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/xelaj/go-dry"
	
	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

func createParamsStructFromMethod(method tlparser.Method) tlparser.Object {
	return tlparser.Object{
		Name:       method.Name + "Params",
		CRC:        method.CRC,
		Parameters: method.Parameters,
	}
}

func haveOptionalParams(params []tlparser.Parameter) bool {
	for _, param := range params {
		if param.IsOptional {
			return true
		}
	}

	return false
}

func maxBitflag(params []tlparser.Parameter) int {
	max := 0
	for _, param := range params {
		if param.BitToTrigger > max {
			max = param.BitToTrigger
		}
	}

	return max
}

func goify(name string, public bool) string {
	delim := strcase.ToDelimited(name, '|')
	delim = strings.ReplaceAll(delim, ".", "|") // strace не видит точки!!
	splitted := strings.Split(delim, "|")
	for i, item := range splitted {
		item = strings.ToLower(item)
		if dry.SliceContains(capitalizePatterns, item) {
			item = strings.ToUpper(item)
		}

		itemRunes := []rune(item)

		if i == 0 && !public {
			// потому что aPI, uRL, это криворуко
			itemRunes = []rune(strings.ToLower(item))
		} else {
			itemRunes[0] = unicode.ToUpper(itemRunes[0])
		}

		splitted[i] = string(itemRunes)
	}

	return strings.Join(splitted, "")
}

func (g *Generator) typeIdFromSchemaType(t string) *jen.Statement {
	item := &jen.Statement{}
	switch t {
	case "Bool":
		item = jen.Bool()
	case "long":
		item = jen.Int64()
	case "double":
		item = jen.Float64()
	case "int":
		item = jen.Int32()
	case "string":
		item = jen.String()
	case "bytes":
		item = jen.Index().Byte()
	case "bitflags":
		panic("bitflags cant be generated or even cath in this part")
	case "true":
		item = jen.Bool()
	default:
		if _, ok := g.schema.Enums[t]; ok {
			item = jen.Id(goify(t, true))
			break
		}
		if _, ok := g.schema.Types[t]; ok {
			item = jen.Id(goify(t, true))
			break
		}
		found := false
		for _, _struct := range g.schema.SingleInterfaceTypes {
			if _struct.Interface == t {
				item = jen.Id("*" + goify(_struct.Name, true))
				found = true
				break
			}
		}
		if found {
			break
		}
		//pp.Fprintln(os.Stderr, g.schema)
		panic("пробовали обработать '" + t + "'")
	}

	return item
}
