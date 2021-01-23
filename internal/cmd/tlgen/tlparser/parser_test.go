package tlparser

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoadTestFile(file string) string {
	_, filename, _, _ := runtime.Caller(0) // nolint:dogsled cause we don't need another stuff
	f, err := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "testdata", file))
	if err != nil {
		panic(err)
	}

	return string(f)
}

func TestSimplestFixture(t *testing.T) {
	file := LoadTestFile("simplest.tl")

	schema, err := ParseSchema(file)

	assert.NoError(t, err)
	assert.Equal(t, &Schema{
		Objects: []Object{
			{
				Name:       "someEnum",
				CRC:        0x5508ec75,
				Parameters: nil,
				Interface:  "CoolEnumerate",
			},
		},
		Methods: []Method{
			{
				Name:       "someFunc",
				CRC:        0x7da07ec9,
				Parameters: nil,
				Response: MethodResponse{
					Type:   "CoolEnumerate",
					IsList: false,
				},
			},
		},
		TypeComments: make(map[string]string),
	}, schema)
}

func TestWithCommentsFixture(t *testing.T) {
	file := LoadTestFile("with_comments.tl")

	schema, err := ParseSchema(file)
	assert.NoError(t, err)
	assert.Equal(t, &Schema{
		Objects: []Object{
			{
				Name:      "someEnum",
				Comment:   "this is really cool enum!",
				CRC:       0x5508ec75,
				Interface: "CoolEnumerate",
			}, {
				Name:    "constructor",
				Comment: "make it struct!",
				CRC:     0x12345678,
				Parameters: []Parameter{
					{
						Name:         "paramA",
						Type:         "string",
						Comment:      "this is the A paramateter",
						IsVector:     false,
						IsOptional:   false,
						BitToTrigger: 0,
					},
					{
						Name:         "paramB",
						Type:         "bool",
						Comment:      "and this is the B one.",
						IsVector:     false,
						IsOptional:   false,
						BitToTrigger: 0,
					},
				},
				Interface: "CoolStruct",
			},
		},
		TypeComments: map[string]string{
			"CoolEnumerate": "some cool enums!",
		},
	}, schema)
}
