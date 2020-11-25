// init.go нужен для того что бы иметь один единственный инит в пакете, и при этом не заставляет его генерироваться через generate-tl-files
package telegram

// "github.com/xelaj/mtproto/serialize"

const (
	ApiVersion = 121
)

func init() {
	//serialize.AddObjectConstructor(GenerateStructByConstructor)
}
