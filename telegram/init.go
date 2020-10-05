// init.go нужен для того что бы иметь один единственный инит в пакете, и при этом не заставляет его генерироваться через generate-tl-files
package telegram

import (
	"github.com/xelaj/mtproto/serialize"
)

func init() {
	serialize.AddObjectConstructor(GenerateStructByConstructor)
}
