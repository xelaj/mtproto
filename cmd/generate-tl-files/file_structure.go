package main

type FileStructure struct {
	Types                    map[typeName][]*StructObject
	SingleInterfaceTypes     []*StructObject
	SingleInterfaceCanonical map[typeName]typeName
	Enums                    map[typeName][]*EnumObject
	Methods                  []*FuncObject

	_d *FileDeclarations //? кэш Declarations(), т.к. так удобнее, Declarations предоставляет более простое описание файла, но долго считает, поэтому так проще
}

func FileFromTlSchema(schema *TLSchema) (*FileStructure, error) {
	res := &FileStructure{
		Enums:                    make(map[typeName][]*EnumObject),
		Types:                    make(map[typeName][]*StructObject),
		SingleInterfaceTypes:     make([]*StructObject, 0),
		SingleInterfaceCanonical: make(map[typeName]typeName),
	}

	// реверсим, т.к. все обозначается по интерфейсам, а на конструкторы насрать видимо.
	reversedObjects := make(map[typeName][]*DefinitionObject)
	for _, obj := range schema.Objects {
		if reversedObjects[obj.Interface] == nil {
			reversedObjects[obj.Interface] = make([]*DefinitionObject, 0)
		}

		reversedObjects[obj.Interface] = append(reversedObjects[obj.Interface], obj)
	}

	for interfaceName, objects := range reversedObjects {
		// может это енум?
		if InterfaceIsEnum(objects) {
			enums := make([]*EnumObject, len(objects))
			for i, obj := range objects {
				enums[i] = &EnumObject{
					Name:    obj.Constructor,
					CRCCode: obj.CRC,
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
		if InterfaceIsSpecific(objects) {
			if len(objects) != 1 {
				panic("defined as single object, but in real has multiple constructors")
			}
			singleObject := &StructObject{
				Name:    objects[0].Constructor,
				CRCCode: objects[0].CRC,
				Fields:  objects[0].Parameters,
			}

			res.SingleInterfaceTypes = append(res.SingleInterfaceTypes, singleObject)
			res.SingleInterfaceCanonical[interfaceName] = singleObject.Name
			delete(reversedObjects, interfaceName)
			continue
		}
		// ну тогда это просто объект с интерфейсом получается

		resultStructs := make([]*StructObject, len(objects))
		for i, obj := range objects {
			constructor := obj.Constructor
			// некоторые конструкторы абсолютно идентичны типу
			if normalizeID(constructor, false) == normalizeID(interfaceName, true) {
				constructor += "Obj"
			}

			resultStructs[i] = &StructObject{
				Name:    constructor,
				CRCCode: obj.CRC,
				Fields:  obj.Parameters,
			}
		}
		res.Types[interfaceName] = resultStructs
	}

	// погнали по методам
	for _, method := range schema.Methods {
		HasOptional := false
		for _, param := range method.Parameters {
			if param.IsOptional {
				HasOptional = true
				break
			}
		}

		m := &FuncObject{
			CRCCode:     method.CRC,
			Name:        method.MethodName,
			Arguments:   method.Parameters,
			HasOptional: HasOptional,
			Returns:     method.Response,
		}

		res.Methods = append(res.Methods, m)
	}

	return res, nil
}

func (s *FileStructure) GetAllConstructors() (structs, enums map[uint32]string) {
	structs = make(map[uint32]string)
	for _, items := range s.Types {
		for _, _struct := range items {
			structs[_struct.CRCCode] = normalizeID(_struct.Name, false)
		}
	}
	for _, _struct := range s.SingleInterfaceTypes {
		for k, v := range s.SingleInterfaceCanonical {
			if v == _struct.Name {
				structs[_struct.CRCCode] = normalizeID(k, false)
			}
		}

	}

	enums = make(map[uint32]string)
	for _, items := range s.Enums {
		for _, enum := range items {
			enums[enum.CRCCode] = normalizeID(enum.Name, false)
		}
	}

	return structs, enums
}

func InterfaceIsEnum(in []*DefinitionObject) bool {
	for _, obj := range in {
		if obj.Parameters == nil || len(obj.Parameters) > 0 {
			return false
		}
	}
	return true
}

func InterfaceIsSpecific(in []*DefinitionObject) bool {
	return len(in) == 1
}

type typeName = string

type EnumObject struct {
	Name    string
	CRCCode uint32
}

type StructObject struct {
	Name    string
	CRCCode uint32
	Fields  []*Param
}

type FuncObject struct {
	CRCCode     uint32
	Name        string
	Arguments   []*Param
	HasOptional bool
	TooManyArgs bool
	Returns     *MethodResponse
}

type FileDeclarations struct {
	Enums           map[string]string
	SpecificStructs map[string]string
	Interfaces      map[string]string
	Methods         map[string]string
}

type typ uint8

const (
	typUnknown typ = iota
	typEnum
	typStruct
	typInterface
	typMethod
)

func (d *FileDeclarations) Find(canonicalTypeName string) (generatedName string, realType typ) {
	if v, ok := d.Enums[canonicalTypeName]; ok {
		return v, typEnum
	}
	if v, ok := d.SpecificStructs[canonicalTypeName]; ok {
		return v, typStruct
	}
	if v, ok := d.Interfaces[canonicalTypeName]; ok {
		return v, typInterface
	}
	if v, ok := d.Methods[canonicalTypeName]; ok {
		return v, typMethod
	}

	return "", typUnknown
}
