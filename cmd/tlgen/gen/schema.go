package gen

import (
	"strings"

	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

type internalSchema struct {
	Types                    map[string][]tlparser.Object
	SingleInterfaceTypes     []tlparser.Object
	SingleInterfaceCanonical map[string]string
	Enums                    map[string][]enum
	Methods                  []tlparser.Method
}

type enum struct {
	Name string
	CRC  uint32
}

func (g *Generator) createInternalSchema(schema tlparser.Schema) (*internalSchema, error) {
	res := &internalSchema{
		Enums:                    make(map[string][]enum),
		Types:                    make(map[string][]tlparser.Object),
		SingleInterfaceTypes:     make([]tlparser.Object, 0),
		SingleInterfaceCanonical: make(map[string]string),
		Methods:                  make([]tlparser.Method, 0),
	}

	// реверсим, т.к. все обозначается по интерфейсам, а на конструкторы насрать видимо.
	reversedObjects := make(map[string][]tlparser.Object)
	for _, obj := range schema.Objects {
		if reversedObjects[obj.Interface] == nil {
			reversedObjects[obj.Interface] = make([]tlparser.Object, 0)
		}

		reversedObjects[obj.Interface] = append(reversedObjects[obj.Interface], obj)
	}

	for interfaceName, objects := range reversedObjects {
		// может это енум?
		if interfaceIsEnum(objects) {
			enums := make([]enum, len(objects))
			for i, obj := range objects {
				enums[i] = enum{
					Name: obj.Name,
					CRC:  obj.CRC,
				}
			}

			res.Enums[interfaceName] = enums
			continue
		}

		// а может конкретная структура?
		// если у нас в интерфейсе только одна структура, то считаем ее уникальной, и пихаем ее как сингл-тип (потому что зачем лишний интерфейс делать)
		//? Зачем нужны типы с единственным конструктором:
		//? есть такое предположение: возможно конструкторы разбросаны по типам, аггрегируя их (в tl схеме телеги на 2000 с хреном строк всего 300+ интерфейсов)
		//? ВОЗМОЖНО (не уверен), сервер проверяет на типизацию так: он сначала проходится по типам (интерфейсам), которых не так много, и в каждом типе проверяет, соблюдает ли
		//? конструктор этот тип (интерфейс), если нет, то идет дальше. ВОЗМОЖНО это сделано чисто для оптимизации, хуй его знает. Но другого объяснения, почему в методы
		//? отдают вот прям только интерфейсы и ничего больше, у меня вариантов тупо нет
		if interfaceIsSpecific(objects) {
			if len(objects) != 1 {
				panic("defined as single object, but in real has multiple constructors")
			}

			// singleObject := &tlparser.Object{
			// 	Name:       objects[0].Name,
			// 	CRC:        objects[0].CRC,
			// 	Parameters: objects[0].Parameters,
			// 	Interface:  objects[0].Interface,
			// }

			res.SingleInterfaceTypes = append(res.SingleInterfaceTypes, objects[0])
			res.SingleInterfaceCanonical[interfaceName] = objects[0].Name
			delete(reversedObjects, interfaceName)
			continue
		}

		// ну тогда это просто объект с интерфейсом получается
		resultStructs := make([]tlparser.Object, len(objects))
		for i, obj := range objects {
			constructor := obj.Name
			// некоторые конструкторы абсолютно идентичны типу
			if strings.ToLower(constructor) == strings.ToLower(interfaceName) {
				constructor += "Obj"
			}

			resultStructs[i] = tlparser.Object{
				Name:       constructor,
				CRC:        obj.CRC,
				Parameters: obj.Parameters,
				Interface:  obj.Interface,
			}
		}
		res.Types[interfaceName] = resultStructs
	}

	// погнали по методам
	for _, method := range schema.Methods {
		res.Methods = append(res.Methods, method)
	}

	return res, nil
}

func (g *Generator) getAllConstructors(s *internalSchema) (structs, enums map[uint32]string) {
	structs = make(map[uint32]string)
	for _, items := range s.Types {
		for _, _struct := range items {
			structs[_struct.CRC] = g.goify(_struct.Name)
		}
	}
	for _, _struct := range s.SingleInterfaceTypes {
		for k, v := range s.SingleInterfaceCanonical {
			if v == _struct.Name {
				structs[_struct.CRC] = g.goify(k)
			}
		}

	}

	enums = make(map[uint32]string)
	for _, items := range s.Enums {
		for _, enum := range items {
			enums[enum.CRC] = g.goify(enum.Name)
		}
	}

	return structs, enums
}

func interfaceIsEnum(in []tlparser.Object) bool {
	for _, obj := range in {
		if len(obj.Parameters) > 0 {
			return false
		}
	}

	return true
}

func interfaceIsSpecific(in []tlparser.Object) bool {
	return len(in) == 1
}
