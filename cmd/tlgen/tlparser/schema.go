package tlparser

type Schema struct {
	Objects      []Object
	Methods      []Method
	TypeComments map[string]string
}

type Object struct {
	Name       string
	Comment    string
	CRC        uint32
	Parameters []Parameter
	Interface  string
}

type Parameter struct {
	Name         string
	Type         string
	Comment      string
	IsVector     bool
	IsOptional   bool
	BitToTrigger int
}

type Method struct {
	Name       string
	CRC        uint32
	Comment    string
	Parameters []Parameter
	Response   MethodResponse
}

type MethodResponse struct {
	Type   string
	IsList bool
}
