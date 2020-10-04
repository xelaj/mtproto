package main

import (
	"github.com/dave/jennifer/jen"
)

func GenerateEnumDefinitions(file *jen.File, data *FileStructure) error {
	for enumType, values := range data.Enums {
		file.Add(GenerateSpecificEnum(enumType, values)...)
	}
	return nil
}

func GenerateSpecificEnum(enumType string, enumValues []*EnumObject) []jen.Code {
	total := make([]jen.Code, 0)

	typeId := normalizeID(enumType, true)

	enumDef := jen.Type().Id(typeId).Uint32()
	total = append(total, enumDef, jen.Line())

	opc := make([]jen.Code, len(enumValues))
	cases := make([]jen.Code, len(enumValues))
	for i, id := range enumValues {
		name := normalizeID(id.Name, false)

		opc[i] = jen.Id(name).Id(typeId).Op("=").Lit(int(id.CRCCode))
		cases[i] = jen.Case(jen.Id(typeId).Call(jen.Lit(int(id.CRCCode)))).Block(jen.Return(jen.Lit(id.Name)))
	}

	total = append(total, jen.Const().Defs(opc...), jen.Line())

	cases = append(cases, jen.Default().Block(jen.Return(jen.Lit("<UNKNOWN "+enumType+">"))))

	f := jen.Func().Params(jen.Id("e").Id(typeId)).Id("String").Params().String().Block(
		jen.Switch(jen.Id("e")).Block(cases...),
	)

	total = append(total, f, jen.Line())

	//CRC() uint32
	f = jen.Func().Params(jen.Id("e").Id(typeId)).Id("CRC").Params().Uint32().Block(
		jen.Return(jen.Uint32().Call(jen.Id("e"))),
	)

	total = append(total, f, jen.Line())

	// Ecncode() []byte
	f = jen.Func().Params(jen.Id("e").Id(typeId)).Id("Encode").Params().Index().Byte().Block(
		jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto", "NewEncoder").Call(),
		jen.Id("buf.PutCRC").Call(jen.Uint32().Call(jen.Id("e"))),
		jen.Line(),
		jen.Return(jen.Id("buf.Result").Call()),
	)

	total = append(total, f, jen.Line())

	return total

	/*
		!type InputWallPaper uint32
		!
		!const (
		!	InputWallPaperNoFile InputWallPaper = 2217196460
		!)
		!
		!func (e *InputWallPaper) String() string {
		!	switch e {
		!	case 2217196460:
		!		return "inputWallPaperNoFile"
		!	default:
		!		return "<UNKNOWN input.WallPaper>"
		!	}
		!}
		!
		!func (e *InputWallPaper) Encode() []byte {
		!	buf := mtproto.NewEncoder()
		!	buf.PutUint(uint32(e))
		!
		!	return buf.Result()
		!}
	*/
}
