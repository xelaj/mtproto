package tlparser

import "fmt"

var excludedDefinitions = []string{
	"true",
	"boolFalse",
	"boolTrue",
	"vector",

	"invokeAfterMsg",
	"invokeAfterMsgs",
	"initConnection",
	"invokeWithLayer",
	"invokeWithoutUpdates",
	"invokeWithMessagesRange",
	"invokeWithTakeout",
}

var excludedTypes = []string{
	"int",
	"long",
	"double",
	"string",
	"bytes",
}

type errExcluded struct {
	name string
}

func (e errExcluded) Error() string {
	return fmt.Sprintf("excluded: %s", e.name)
}
