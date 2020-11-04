package gen

import (
	"strings"

	"github.com/xelaj/mtproto/cmd/tlgen/typelang"
)

type internalSchema struct {
	Types                    map[string][]typelang.Object
	SingleInterfaceTypes     []typelang.Object
	SingleInterfaceCanonical map[string]string // uses for searching from canonical id to go-ified
	Enums                    map[string][]enum
	Methods                  []typelang.Method
}

type enum struct {
	Name string
	CRC  uint32
}

func createInternalSchema(schema *typelang.Schema) *internalSchema {
	ischem := &internalSchema{
		Enums:                    make(map[string][]enum),
		Types:                    make(map[string][]typelang.Object),
		SingleInterfaceTypes:     make([]typelang.Object, 0),
		SingleInterfaceCanonical: make(map[string]string),
		Methods:                  make([]typelang.Method, 0),
	}

	// objects are reversing cause it's generated code is more readable with structs grouped by interface
	reversedObjects := make(map[string][]typelang.Object)
	for _, obj := range schema.Objects {
		if reversedObjects[obj.Interface] == nil {
			reversedObjects[obj.Interface] = make([]typelang.Object, 0)
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

			ischem.Enums[interfaceName] = enums
			continue
		}

		// а может конкретная структура?
		// если у нас в интерфейсе только одна структура, то считаем ее уникальной, и пихаем ее как сингл-тип
		// (потому что зачем лишний интерфейс делать)
		//? Зачем нужны типы с единственным конструктором:
		//? есть такое предположение: возможно конструкторы разбросаны по типам, аггрегируя их (в tl схеме
		//? телеги на 2000 с хреном строк всего 300+ интерфейсов)  ВОЗМОЖНО (не уверен), сервер проверяет на
		//? типизацию так: он сначала проходится по типам (интерфейсам), которых не так много, и в каждом типе
		//? проверяет, соблюдает ли конструктор этот тип (интерфейс), если нет, то идет дальше. ВОЗМОЖНО это
		//? сделано чисто для оптимизации, хуй его знает. Но другого объяснения, почему в методы отдают вот
		//? прям только интерфейсы и ничего больше, у меня вариантов тупо нет
		if interfaceIsSpecific(objects) {
			ischem.SingleInterfaceTypes = append(ischem.SingleInterfaceTypes, objects[0])
			ischem.SingleInterfaceCanonical[interfaceName] = objects[0].Name
			delete(reversedObjects, interfaceName)
			continue
		}

		// ну тогда это просто объект с интерфейсом получается
		resultStructs := make([]typelang.Object, len(objects))
		for i, obj := range objects {
			constructor := obj.Name
			// некоторые конструкторы абсолютно идентичны типу
			if strings.EqualFold(constructor, interfaceName) {
				constructor += "Obj"
			}

			resultStructs[i] = typelang.Object{
				Name:       constructor,
				CRC:        obj.CRC,
				Parameters: obj.Parameters,
				Interface:  obj.Interface,
			}
		}
		ischem.Types[interfaceName] = resultStructs
	}

	// погнали по методам
	ischem.Methods = append(ischem.Methods, schema.Methods...)

	return ischem
}

func (g *Generator) getAllConstructors() (structs, enums map[uint32]string) {
	structs = make(map[uint32]string)

	for _, items := range g.schema.Types {
		for _, _struct := range items {
			structs[_struct.CRC] = noramlizeIdentificator(_struct.Name)
		}
	}
	for _, _struct := range g.schema.SingleInterfaceTypes {
		for k, v := range g.schema.SingleInterfaceCanonical {
			if v == _struct.Name {
				structs[_struct.CRC] = noramlizeIdentificator(k)
			}
		}
	}

	enums = make(map[uint32]string)
	for _, items := range g.schema.Enums {
		for _, enum := range items {
			enums[enum.CRC] = noramlizeIdentificator(enum.Name)
		}
	}

	return structs, enums
}

func interfaceIsEnum(in []typelang.Object) bool {
	for _, obj := range in {
		if len(obj.Parameters) > 0 {
			return false
		}
	}

	return true
}

func interfaceIsSpecific(in []typelang.Object) bool {
	return len(in) == 1
}
