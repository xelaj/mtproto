package tlparser

import (
	"fmt"

	"github.com/xelaj/go-dry"
)

var excludedDefinitions = dry.SliceUnique([]string{
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
}).(map[string]struct{})

var excludedTypes = dry.SliceUnique([]string{
	"int",
	"long",
	"double",
	"string",
	"bytes",
}).(map[string]struct{})

type errExcluded struct {
	name string
}

func (e errExcluded) Error() string {
	return fmt.Sprintf("excluded: %s", e.name)
}
