package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateInit(file *jen.File) {
	initFunc := jen.Func().Id("init").Params().Block(
		g.createInitStructs(),
		jen.Line(),
		g.createInitEnums(),
	)

	file.Add(initFunc)
}

func (g *Generator) createInitStructs() jen.Code {
	structs, _ := g.getAllConstructors()
	crcs := make([]uint32, 0)

	for crc := range structs {
		crcs = append(crcs, crc)
	}

	sort.Slice(crcs, func(i, j int) bool {
		return crcs[i] < crcs[j]
	})

	var stmts []jen.Code
	for _, crc := range crcs {
		name := structs[crc]
		stmts = append(stmts, jen.Line().Id("&"+name).Values())
	}

	stmts = append(stmts, jen.Line())
	return jen.Qual("github.com/xelaj/mtproto/encoding/tl", "RegisterObjects").Call(stmts...)
}

func (g *Generator) createInitEnums() jen.Code {
	_, enums := g.getAllConstructors()
	crcs := make([]uint32, 0)

	for crc := range enums {
		crcs = append(crcs, crc)
	}

	sort.Slice(crcs, func(i, j int) bool {
		return crcs[i] < crcs[j]
	})

	var stmts []jen.Code
	for _, crc := range crcs {
		name := enums[crc]
		stmts = append(stmts, jen.Line().Id(name))
	}

	stmts = append(stmts, jen.Line())
	return jen.Qual("github.com/xelaj/mtproto/encoding/tl", "RegisterEnums").Call(stmts...)
}
