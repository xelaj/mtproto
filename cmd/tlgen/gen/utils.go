package gen

import "github.com/xelaj/mtproto/cmd/tlgen/tlparser"

func createParamsStructFromMethod(method tlparser.Method) tlparser.Object {
	return tlparser.Object{
		Name:       method.Name + "Params",
		CRC:        method.CRC,
		Parameters: method.Parameters,
		Interface:  "<???>",
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
