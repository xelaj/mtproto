package tlparser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// строчка в tl
type definition struct { //nolint:maligned для удобочитаемости, сохранен порядок филдов как в схеме
	Name       string      // название в самом начале
	CRC        uint32      // crc после #
	Params     []Parameter // параметры после crc
	EqType     string      // тип после параметров и знака равенства
	IsEqVector bool        // тип после знака равенства векторный?
}

func ParseSchema(source string) (*Schema, error) {
	cur := NewCursor(source)

	var (
		objects       []Object
		methods       []Method
		typeComments  = make(map[string]string)
		paramComments = make(map[string]string)

		isFunctions        = false
		nextTypeComment    string
		constructorComment string
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
			cur.SkipSpaces()
			ctype, err := cur.ReadAt(' ')
			if err != nil {
				return nil, fmt.Errorf("read comment type: %w", err)
			}

			cur.SkipSpaces()

			switch ctype {
			case "@type":
				comment, err := cur.ReadAt('\n')
				if err != nil {
					return nil, fmt.Errorf("read comment: %w", err)
				}
				nextTypeComment = strings.TrimSpace(comment)
			case "@enum", "@constructor", "@method":
				comment, err := cur.ReadAt('\n')
				if err != nil {
					return nil, fmt.Errorf("read comment: %w", err)
				}
				constructorComment = strings.TrimSpace(comment)
			case "@param":
				pname, err := cur.ReadAt(' ')
				if err != nil {
					return nil, fmt.Errorf("read comment param name: %w", err)
				}

				cur.SkipSpaces()
				pcomment, err := cur.ReadAt('\n')
				if err != nil {
					return nil, fmt.Errorf("read comment param: %w", err)
				}

				paramComments[pname] = strings.TrimSpace(pcomment)
			default:
				return nil, fmt.Errorf("unknown comment type: %s", ctype)
			}

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

			return nil, err
		}

		for i := range def.Params {
			def.Params[i].Comment = paramComments[def.Params[i].Name]
		}

		if isFunctions {
			methods = append(methods, Method{
				Name:       def.Name,
				Comment:    constructorComment,
				CRC:        def.CRC,
				Parameters: def.Params,
				Response: MethodResponse{
					Type:   def.EqType,
					IsList: def.IsEqVector,
				},
			})
			continue
		} else {
			// если мы спарсили объект
			// тип после знака равенства это интерфейс
			// и он не может быть вектором
			if def.IsEqVector {
				return nil, errors.New("type can't be a vector")
			}

			objects = append(objects, Object{
				Name:       def.Name,
				Comment:    constructorComment,
				CRC:        def.CRC,
				Parameters: def.Params,
				Interface:  def.EqType,
			})
		}

		if nextTypeComment != "" {
			typeComments[def.EqType] = nextTypeComment
			nextTypeComment = ""
		}
		constructorComment = ""
		paramComments = make(map[string]string)
	}

	return &Schema{
		Objects:      objects,
		Methods:      methods,
		TypeComments: typeComments,
	}, nil
}

func parseDefinition(cur *Cursor) (def definition, err error) {
	cur.SkipSpaces()

	{
		// блок тупо для проверки записей типа таких:
		// int ? = Int;
		//    ↑ - читаем до этого пробела и смотрим тип
		var typSpace string
		typSpace, err = cur.ReadAt(' ') // typSpace is int for example
		if err != nil {
			return def, fmt.Errorf("parse def row: %w", err)
		}

		// если он в excludedTypes
		if _, found := excludedTypes[typSpace]; found {
			// дочитываем строку до конца
			if _, err = cur.ReadAt(';'); err != nil {
				return def, err
			}

			cur.Skip(1) // пропускаем ';'
			// говорим что прочитали хуйню
			return def, errExcluded{typSpace}
		}

		cur.Unread(len(typSpace))
	}

	// ipPort#d433ad73 ipv4:int port:int = IpPort;
	//       ↑ - читаем до решеточки, получаем название
	def.Name, err = cur.ReadAt('#') // def.Name is ipPort for this example
	if err != nil {
		return def, fmt.Errorf("parse object name: %w", err)
	}

	// проверяем название на наличие в блоклисте
	if _, found := excludedDefinitions[def.Name]; found {
		// дочитываем строку до конца
		if _, err = cur.ReadAt(';'); err != nil {
			return def, err
		}

		cur.Skip(1) // пропускаем ';'
		// говорим что прочитали хуйню
		return def, errExcluded{def.Name}
	}

	cur.Skip(1) // skip #

	//        ↓ - курсор здесь
	// ipPort#d433ad73 ipv4:int port:int = IpPort;
	//                ↑ - читаем до пробела, получаем crc
	crcString, err := cur.ReadAt(' ')
	if err != nil {
		return def, fmt.Errorf("parse object crc: %w", err)
	}

	cur.SkipSpaces()

	//                 ↓ - курсор здесь
	// ipPort#d433ad73 ipv4:int port:int = IpPort;

	// читаем параметры
	for !cur.IsNext("=") {
		var param Parameter
		param, err = parseParam(cur)
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
	// читаем тип после знака =
	//                                     ↓ - курсор здесь
	// ipPort#d433ad73 ipv4:int port:int = IpPort;
	if cur.IsNext("Vector") {
		cur.Skip(1) // skip <
		def.EqType, err = cur.ReadAt('>')
		if err != nil {
			return def, fmt.Errorf("parse def eq type: %w", err)
		}

		def.IsEqVector = true
		cur.Skip(len(">;")) // skip >;
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
	// ↓ - курсор здесь
	// correct_answers:flags.0?Vector<bytes> foo:bar
	//                ↑ - читаем до двоеточия, получаем название параметра
	param.Name, err = cur.ReadAt(':')
	if err != nil {
		return param, fmt.Errorf("read param name: %w", err)
	}
	cur.Skip(1) // skip :

	//                  ↓ - курсор здесь
	//  correct_answers:flags.0?Vector<bytes> foo:bar

	// если после ':' идет flags.
	if cur.IsNext("flags.") {
		//                       ↓ - курсор здесь
		// correct_answers:flags.0?Vector<bytes> foo:bar

		// читаем цифры
		var digits string
		digits, err = cur.ReadDigits() // read bit index, must be digit
		if err != nil {
			return param, fmt.Errorf("read param bitflag: %w", err)
		}

		param.BitToTrigger, err = strconv.Atoi(digits)
		if err != nil {
			return param, fmt.Errorf("invalid bitflag index: %s", digits)
		}

		//                        ↓ - курсор здесь
		// correct_answers:flags.0?Vector<bytes> foo:bar
		// ожидаем знак '?'
		if !cur.IsNext("?") {
			return param, fmt.Errorf("expected '?'")
		}
		param.IsOptional = true
	}

	// читаем тип параметра
	if cur.IsNext("Vector") {
		//                               ↓ - курсор здесь
		// correct_answers:flags.0?Vector<bytes> foo:bar

		cur.Skip(1) // skip <
		param.IsVector = true
		param.Type, err = cur.ReadAt('>')
		if err != nil {
			return param, fmt.Errorf("read param type: %w", err)
		}

		cur.Skip(1) // skip >
	} else {
		// если не вектор, просто вычитываем тип до пробела
		param.Type, err = cur.ReadAt(' ')
		if err != nil {
			return param, fmt.Errorf("read param type: %w", err)
		}
	}

	return param, nil
}
