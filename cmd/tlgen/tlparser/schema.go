package tlparser

type Schema struct {
	Objects []Object
	Methods []Method
}

type Object struct {
	Name       string
	CRC        uint32
	Parameters []Parameter
	Interface  string
}

type Parameter struct {
	Name         string
	Type         string
	IsVector     bool
	IsOptional   bool
	BitToTrigger int
}

type Method struct {
	Name string
	CRC        uint32
	Parameters []Parameter
	Response   MethodResponse
}

type MethodResponse struct {
	Type   string
	IsList bool
}


