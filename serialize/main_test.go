package serialize

import (
	"os"
	"testing"
)

func tearup() {
	AddObjectConstructor(generateDummyObjects)
}

func teardown() {

}

func TestMain(m *testing.M) {
	tearup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
