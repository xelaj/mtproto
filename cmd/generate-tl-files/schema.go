package main

import (
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/willf/pad"
	"github.com/xelaj/go-dry"
)

type DefinitionObject struct {
	Constructor string
	CRC         uint32
	Interface   string
	Parameters  []*Param
}

type Param struct {
	Name         string
	Type         string
	IsList       bool
	IsOptional   bool // если в объявлении есть flags:#
	BitToTrigger int  // бит, который нужно триггернуть у flags что бы указать, что опциональное поле есть
}

type MethodResponse struct {
	Type   string
	IsList bool
}

type DefinitionMethod struct {
	MethodName string
	CRC        uint32
	Response   *MethodResponse
	Parameters []*Param
}


func GetCRCCode(constructorTrimmedLine string) (string, uint32, error) {
	var crcCode uint32
	crcCodeStr := regexp.MustCompilePOSIX("#[0-9a-f]{1,8}").FindString(constructorTrimmedLine)
	if crcCodeStr == "" {
		return "", 0, nil
	}

	crcCodeStr = strings.TrimPrefix(crcCodeStr, "#")
	crcCodeStr = pad.Left(crcCodeStr, 8, "0") // 8 потому что 8 hex символов на int32

	b, err := hex.DecodeString(crcCodeStr)
	if err != nil {
		return "", 0, errors.Wrap(err, "parsing CRC code")
	}
	crcCode = binary.BigEndian.Uint32(b)

	return "#" + crcCodeStr, crcCode, nil
}

func ParseFunction(line string) (*DefinitionMethod, error) {
	if strings.HasPrefix(line, "//") || line == "" {
		return nil, nil
	}

	methodName := regexp.MustCompilePOSIX("^[a-z][a-zA-Z0-9_.]+").FindString(line)
	switch methodName {
	// исключения, не ебу зачем вообще обозначены
	case "invokeAfterMsg",
		"invokeAfterMsgs",
		"initConnection",
		"invokeWithLayer",
		"invokeWithoutUpdates",
		"invokeWithMessagesRange",
		"invokeWithTakeout":
		return nil, nil
	}

	if methodName == "" {
		return nil, errors.New("method wasn't defined")
	}

	line = strings.TrimPrefix(line, methodName)

	str, crcCode, err := GetCRCCode(line)
	if err != nil {
		return nil, err
	}
	// line == "one:two {three:four} = SomeTypeToInterface;"
	line = strings.TrimLeft(line, str)
	line = strings.TrimLeft(line, " ")

	params := make([]*Param, 0)
	// mathing "= SomeTypeToInterface;"
	responseType := regexp.MustCompilePOSIX("=[ ]?(Vector<)?([a-z]+.)?[a-zA-Z0-9_.]+(>)?;").FindString(line)
	responseLeft := responseType == line
	if !responseLeft {
		hasOptionalFields := false

		line = strings.TrimSpace(strings.TrimSuffix(line, responseType))
		parameters := strings.Split(line, " ")
		for _, paramStr := range parameters {
			splitted := strings.Split(paramStr, ":")
			if len(splitted) != 2 {
				return nil, errors.New("incorrect parameter declaration: " + paramStr)
			}
			key := splitted[0]
			typ := splitted[1]

			if key == "flags" && typ == "#" {
				hasOptionalFields = true
				typ = "bitflags"
			}

			p := &Param{
				Name: key,
				Type: typ,
			}

			if hasOptionalFields {
				if strings.HasPrefix(typ, "flags.") {
					typ = strings.TrimPrefix(typ, "flags.")
					triggeringBitStr := regexp.MustCompilePOSIX("^[0-9]+").FindString(typ)
					if triggeringBitStr == "" {
						return nil, errors.New("expected number of bit for triggering")
					}

					typ = strings.TrimPrefix(typ, triggeringBitStr)

					triggeringBit, err := strconv.Atoi(triggeringBitStr)
					dry.PanicIfErr(err)

					realType := strings.TrimPrefix(typ, "?")
					p.Type = realType
					p.IsOptional = true
					p.BitToTrigger = triggeringBit
				}
			} else if strings.Contains(typ, "?") {
				return nil, errors.New("declaration didn't define flags parameter for optional values")
			}

			isVector := strings.HasPrefix(p.Type, "Vector<")
			if isVector {
				p.Type = strings.TrimSuffix(strings.TrimPrefix(p.Type, "Vector<"), ">")
				p.IsList = true
			}

			params = append(params, p)
		}
	}

	responseType = strings.TrimLeftFunc(responseType, func(r rune) bool { return r == '=' || r == ' ' })
	responseType = strings.TrimSuffix(responseType, ";")

	response := &MethodResponse{
		Type: responseType,
	}

	isVector := strings.HasPrefix(response.Type, "Vector<")
	if isVector {
		response.Type = strings.TrimSuffix(strings.TrimPrefix(response.Type, "Vector<"), ">")
		response.IsList = true
	}

	return &DefinitionMethod{
		CRC:        crcCode,
		MethodName: methodName,
		Response:   response,
		Parameters: params,
	}, nil
}

func ParseType(line string) (*DefinitionObject, error) {
	// TODO: СДЕЛАТЬ ЛЕКСЕР ПОД ВСЕ ЭТО а то затрахаться можно
	// line == "constructorName#01234567 one:two {three:four} = SomeTypeToInterface;"
	if strings.HasPrefix(line, "//") || line == "" {
		return nil, nil
	}

	constructor := regexp.MustCompilePOSIX("[a-z][a-zA-Z0-9_.]+").FindString(line)
	switch constructor {
	// исключения, не ебу зачем вообще обозначены
	case "int",
		"long",
		"double",
		"string",
		"bytes",
		"true",
		"boolFalse",
		"boolTrue",
		"vector":
		return nil, nil
	}

	if constructor == "" {
		errors.New("constructor wasn't defined")
	}

	// line == "#01234567 one:two {three:four} = SomeTypeToInterface;"
	line = strings.TrimPrefix(line, constructor)

	var crcCode uint32
	crcCodeStr := regexp.MustCompilePOSIX("#[0-9a-f]{1,8}").FindString(line)
	if crcCodeStr != "" {
		line = strings.TrimPrefix(line, crcCodeStr)
		crcCodeStr = strings.TrimPrefix(crcCodeStr, "#")
		crcCodeStr = pad.Left(crcCodeStr, 8, "0") // 8 потому что 8 hex символов на int32

		b, err := hex.DecodeString(crcCodeStr)
		if err != nil {
			return nil, errors.Wrap(err, "decoding CRC")
		}
		crcCode = binary.BigEndian.Uint32(b)
	}
	// line == "one:two {three:four} = SomeTypeToInterface;"
	line = strings.TrimLeft(line, " ")

	params := make([]*Param, 0)
	// mathing "= SomeTypeToInterface;"
	_interface := regexp.MustCompilePOSIX("=[ ]?([a-z]+.)?[A-Z][a-zA-Z0-9_.]+;").FindString(line)
	onlyInterfaceLeft := _interface == line
	if !onlyInterfaceLeft {
		line = strings.TrimSpace(strings.TrimSuffix(line, _interface))

		hasOptionalFields := false

		parameters := strings.Split(line, " ")
		for _, paramStr := range parameters {
			paramSplitted := strings.Split(paramStr, ":")
			if len(paramSplitted) != 2 {
				return nil, errors.New("incorrect parameter declaration: " + paramStr)
			}

			key := paramSplitted[0]
			typ := paramSplitted[1]

			if key == "flags" && typ == "#" {
				hasOptionalFields = true
				typ = "bitflags"
			}

			p := &Param{
				Name: key,
				Type: typ,
			}

			if hasOptionalFields {
				if strings.HasPrefix(typ, "flags.") {
					typ = strings.TrimPrefix(typ, "flags.")
					triggeringBitStr := regexp.MustCompilePOSIX("^[0-9]+").FindString(typ)
					if triggeringBitStr == "" {
						errors.New("expected number of bit for triggering")
					}

					typ = strings.TrimPrefix(typ, triggeringBitStr)
					triggeringBit, err := strconv.Atoi(triggeringBitStr)
					dry.PanicIfErr(err)

					realType := strings.TrimPrefix(typ, "?")

					p.Type = realType
					p.IsOptional = true
					p.BitToTrigger = triggeringBit
				}
			} else if strings.Contains(typ, "?") {
				errors.New("declaration didn't define flags parameter for optional values")
			}

			isVector := strings.HasPrefix(p.Type, "Vector<")
			if isVector {
				p.Type = strings.TrimSuffix(strings.TrimPrefix(p.Type, "Vector<"), ">")
				p.IsList = true
			}

			params = append(params, p)
		}
	}
	_interface = strings.TrimLeftFunc(_interface, func(r rune) bool { return r == '=' || r == ' ' })
	_interface = strings.TrimSuffix(_interface, ";")

	return &DefinitionObject{
		CRC:         crcCode,
		Constructor: constructor,
		Interface:   _interface,
		Parameters:  params,
	}, nil
}

type TLSchema struct {
	Objects []*DefinitionObject
	Methods []*DefinitionMethod
}

func ParseTL(data string) (*TLSchema, error) {
	objects := make([]*DefinitionObject, 0)
	methods := make([]*DefinitionMethod, 0)
	definingFuncs := false
	for lineNumber, line := range strings.Split(data, "\n") {
		lineNumber++ // т.к. начинаем с нуля, а строчки то с 1
		if strings.Contains(line, "---functions---") {
			definingFuncs = true // функции отдельно отрабатываем
			continue
		}
		if strings.Contains(line, "---types---") {
			definingFuncs = false // типы судя по докам могут и после функций идти
			continue
		}
		var err error
		if definingFuncs {
			var method *DefinitionMethod
			method, err = ParseFunction(line)
			if method != nil {
				methods = append(methods, method)
			}
		} else {
			var object *DefinitionObject
			object, err = ParseType(line)
			if object != nil {
				objects = append(objects, object)
			}
		}

		if err != nil {
			return nil, errors.Wrap(err, "line "+strconv.Itoa(lineNumber))
		}
	}
	return &TLSchema{
		Objects: objects,
		Methods: methods,
	}, nil
}
