package gen

import (
	"strings"

	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

// для понимания как преобразовано название типа
type goifiedName = string
type nativeName = string

// предполагаем, что мы пишем в файл данные через goify, поэтому тут все символы нативные
type internalSchema struct {
	InterfaceCommnets    map[nativeName]string
	Types                map[nativeName][]tlparser.Object
	SingleInterfaceTypes []tlparser.Object
	Enums                map[nativeName][]enum
	Methods              []tlparser.Method
}

type enum struct {
	Name nativeName
	CRC  uint32
}

func createInternalSchema(nativeSchema *tlparser.Schema) (*internalSchema, error) {
	internalSchema := &internalSchema{
		InterfaceCommnets:    make(map[string]string),
		Enums:                make(map[string][]enum),
		Types:                make(map[string][]tlparser.Object),
		SingleInterfaceTypes: make([]tlparser.Object, 0),
		Methods:              make([]tlparser.Method, 0),
	}

	// реверсим, т.к. все обозначается по интерфейсам, а на конструкторы насрать видимо.
	reversedObjects := make(map[string][]tlparser.Object)
	for _, obj := range nativeSchema.Objects {
		if reversedObjects[obj.Interface] == nil {
			reversedObjects[obj.Interface] = make([]tlparser.Object, 0)
		}

		reversedObjects[obj.Interface] = append(reversedObjects[obj.Interface], obj)
	}

	for interfaceName, objects := range reversedObjects {
		// ну тогда это просто объект с интерфейсом получается, раз не енум и не одиночный объект
		for _, obj := range objects {
			// некоторые конструкторы абсолютно идентичны типу по названию
			if strings.ToLower(obj.Name) == strings.ToLower(obj.Interface) {
				obj.Name += "Obj"
			}
		}

		// может это енум?
		if interfaceIsEnum(objects) {
			enums := make([]enum, len(objects))
			for i, obj := range objects {
				enums[i] = enum{
					Name: obj.Name,
					CRC:  obj.CRC,
				}
			}

			internalSchema.Enums[interfaceName] = enums
			continue
		}

		// а может конкретная структура?
		// если у нас в интерфейсе только одна структура, то считаем ее уникальной, и пихаем ее как сингл-тип
		// (потому что зачем лишний интерфейс делать)
		//? Зачем нужны типы с единственным конструктором:
		//? есть такое предположение: возможно конструкторы разбросаны по типам, аггрегируя их (в tl схеме
		//? телеги на 2000 с хреном строк всего 300+ интерфейсов) ВОЗМОЖНО (не уверен), сервер проверяет на
		//? типизацию так: он сначала проходится по типам (интерфейсам), которых не так много, и в каждом типе
		//? проверяет, соблюдает ли конструктор этот тип (интерфейс), если нет, то идет дальше. ВОЗМОЖНО это
		//? сделано чисто для оптимизации, хуй его знает. Но другого объяснения, почему в методы отдают вот
		//? прям только интерфейсы и ничего больше, у меня вариантов тупо нет
		if len(objects) == 1 {
			internalSchema.SingleInterfaceTypes = append(internalSchema.SingleInterfaceTypes, objects[0])
			// delete(reversedObjects, interfaceName)
			continue
		}

		internalSchema.Types[interfaceName] = objects
	}

	internalSchema.Methods = nativeSchema.Methods
	internalSchema.InterfaceCommnets = nativeSchema.TypeComments
	return internalSchema, nil
}

func (g *Generator) getAllConstructors() (structs, enums []goifiedName) {
	structs, enums = make([]string, 0), make([]string, 0)

	for _, items := range g.schema.Types {
		for _, _struct := range items {
			t := goify(_struct.Name, true)
			if goify(_struct.Name, true) == goify(_struct.Interface, true) {
				t = goify(_struct.Name+"Obj", true)
			}
			structs = append(structs, t)
		}
	}
	for _, _struct := range g.schema.SingleInterfaceTypes {
		structs = append(structs, goify(_struct.Name, true))
	}
	for _, method := range g.schema.Methods {
		structs = append(structs, goify(method.Name+"Params", true))
	}

	for _, items := range g.schema.Enums {
		for _, enum := range items {
			enums = append(enums, goify(enum.Name, true))
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
