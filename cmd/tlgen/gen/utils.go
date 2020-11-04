package gen

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/xelaj/mtproto/cmd/tlgen/typelang"
)

func createParamsStructFromMethod(method *typelang.Method) typelang.Object {
	return typelang.Object{
		Name:       method.Name + "Params",
		CRC:        method.CRC,
		Parameters: method.Parameters,
	}
}

func haveOptionalParams(params []typelang.Parameter) bool {
	for _, param := range params {
		if param.IsOptional {
			return true
		}
	}

	return false
}

func maxBitflag(params []typelang.Parameter) int {
	max := 0
	for _, param := range params {
		if param.BitToTrigger > max {
			max = param.BitToTrigger
		}
	}

	return max
}

var mustChange = map[string]string{
	"id":  "ID",
	"p2p": "P2P",
	"url": "URL",
	// may want to change abbriweations like "information"->"Info" or something like this
}

func noramlizeIdentificator(name string) string {
	delimeted := strcase.ToDelimited(name, '.')

	splitted := strings.Split(delimeted, ".")
	for i, elem := range splitted {
		if newOne, ok := mustChange[strings.ToLower(elem)]; ok {
			splitted[i] = newOne
		} else {
			splitted[i] = strings.Title(elem)
		}
	}

	return strings.Join(splitted, "")
}
