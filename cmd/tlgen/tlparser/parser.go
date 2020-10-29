package tlparser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type definition struct {
	Name       string
	CRC        uint32
	Params     []Parameter
	EqType     string
	IsEqVector bool
}

func ParseSchema(source string) (Schema, error) {
	cur := NewCursor(source)

	var (
		objects     []Object
		methods     []Method
		isFunctions = false
	)

	for {
		cur.SkipSpaces()
		if cur.IsNext("---functions---") {
			isFunctions = true
			continue
		}

		if cur.IsNext("---types---") {
			isFunctions = false
			continue
		}

		if cur.IsNext("//") {
			cur.ReadAt('\n')
			cur.Skip(1)
			continue
		}

		def, err := parseDefinition(cur)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.As(err, &errExcluded{}) {
				continue
			}

			return Schema{}, err
		}

		if isFunctions {
			methods = append(methods, Method{
				Name:       def.Name,
				CRC:        def.CRC,
				Parameters: def.Params,
				Response: MethodResponse{
					Type:   def.EqType,
					IsList: def.IsEqVector,
				},
			})
			continue
		}

		if def.IsEqVector {
			panic("wut")
		}

		objects = append(objects, Object{
			Name:       def.Name,
			CRC:        def.CRC,
			Parameters: def.Params,
			Interface:  def.EqType,
		})
	}

	return Schema{
		Objects: objects,
		Methods: methods,
	}, nil
}

func parseDefinition(cur *Cursor) (def definition, err error) {
	cur.SkipSpaces()

	{
		typSpace, err := cur.ReadAt(' ')
		if err != nil {
			return def, fmt.Errorf("parse def row: %w", err)
		}

		if _, found := excludedTypes[typSpace]; found {
			if _, err := cur.ReadAt(';'); err != nil {
				return def, err
			}

			cur.Skip(1)
			return def, errExcluded{typSpace}
		}

		cur.Unread(len(typSpace))
	}

	def.Name, err = cur.ReadAt('#')
	if err != nil {
		return def, fmt.Errorf("parse object name: %w", err)
	}

	if _, found := excludedDefinitions[def.Name]; found {
		if _, err := cur.ReadAt(';'); err != nil {
			return def, err
		}

		cur.Skip(1)
		return def, errExcluded{def.Name}
	}

	cur.Skip(1) // skip #
	crcString, err := cur.ReadAt(' ')
	if err != nil {
		return def, fmt.Errorf("parse object crc: %w", err)
	}

	cur.SkipSpaces()
	for !cur.IsNext("=") {
		param, err := parseParam(cur)
		if err != nil {
			return def, fmt.Errorf("parse param: %w", err)
		}

		cur.SkipSpaces()
		if param.Name == "flags" && param.Type == "#" {
			param.Type = "bitflags"
		}

		def.Params = append(def.Params, param)
	}

	cur.SkipSpaces()
	if cur.IsNext("Vector") {
		cur.Skip(1) // skip <
		def.EqType, err = cur.ReadAt('>')
		if err != nil {
			return def, fmt.Errorf("parse def eq type: %w", err)
		}

		def.IsEqVector = true
		cur.Skip(2) // skip >;
	} else {
		def.EqType, err = cur.ReadAt(';')
		if err != nil {
			return def, fmt.Errorf("parse obj interface: %w", err)
		}

		cur.Skip(1) // skip ;
	}

	crc, err := strconv.ParseUint(crcString, 16, 32)
	if err != nil {
		return def, err
	}
	def.CRC = uint32(crc)
	return def, nil
}

func parseParam(cur *Cursor) (param Parameter, err error) {
	cur.SkipSpaces()

	param.Name, err = cur.ReadAt(':')
	if err != nil {
		return param, fmt.Errorf("read param name: %w", err)
	}
	cur.Skip(1)

	if cur.IsNext("flags.") {
		//fmt.Println("read digit:", string(cur.source[cur.pos-1:cur.pos+10]))
		r, err := cur.ReadDigits() //read bit index, must be digit
		if err != nil {
			return param, fmt.Errorf("read param bitflag: %w", err)
		}

		param.BitToTrigger, err = strconv.Atoi(r)
		if err != nil {
			return param, fmt.Errorf("invalid bitflag index: %s", string(r))
		}

		if !cur.IsNext("?") {
			return param, fmt.Errorf("expected '?'")
		}
		param.IsOptional = true
	}

	if cur.IsNext("Vector") {
		cur.Skip(1) // skip <
		param.IsVector = true
		param.Type, err = cur.ReadAt('>')
		if err != nil {
			return param, fmt.Errorf("read param type: %w", err)
		}

		cur.Skip(1) // skip >
	} else {
		param.Type, err = cur.ReadAt(' ')
		if err != nil {
			return param, fmt.Errorf("read param type: %w", err)
		}
	}

	return param, nil
}
