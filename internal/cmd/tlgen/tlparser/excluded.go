package tlparser

import "fmt"

var excludedDefinitions = map[string]struct{}{
	"true":      {},
	"boolFalse": {},
	"boolTrue":  {},
	"vector":    {},

	"invokeAfterMsg":          {},
	"invokeAfterMsgs":         {},
	"initConnection":          {},
	"invokeWithLayer":         {},
	"invokeWithoutUpdates":    {},
	"invokeWithMessagesRange": {},
	"invokeWithTakeout":       {},
}

var excludedTypes = map[string]struct{}{
	"int":    {},
	"long":   {},
	"double": {},
	"string": {},
	"bytes":  {},
}

type errExcluded struct {
	name string
}

func (e errExcluded) Error() string {
	return fmt.Sprintf("excluded: %s", e.name)
}
