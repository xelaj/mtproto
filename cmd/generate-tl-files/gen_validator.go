package main

import (
	"github.com/dave/jennifer/jen"
)

func generateStructValidatorFunc(str *StructObject, data *FileStructure) (*jen.Statement, error) {
	checks := make([]jen.Code, 0)
	i := 0
	for _, field := range str.Fields {
		if field.IsOptional {
			continue
		}
		i++
		if i == 10 {
			break
		}
		fv, err := createFieldValidation(field, data, false)
		if err != nil {
			return nil, err
		}

		checks = append(checks, fv)
		checks = append(checks, jen.Line())
	}

	checks = append(checks, jen.Return(jen.Id("nil")))
	structName := normalizeID(str.Name, false)
	return jen.Func().Params(jen.Id("e").Id("*" + structName)).Id("Validate").Params().Params(jen.Error()).Block(
		checks...,
	), nil
}

func createFieldValidation(field *Param, data *FileStructure, insideRange bool) (jen.Code, error) {
	name := normalizeID(field.Name, true)
	name = goify(field.Name)
	direct := "e." + name
	if insideRange {
		direct = "item"
	}
	typ := field.Type

	if field.IsList {
		name := normalizeID(field.Name, false)
		name = goify(field.Name)
		direct := "e." + name

		checkLen := jen.If(jen.Len(jen.Id(direct)).Op("==").Id("0")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
		).Comment("slice_len_check")

		cp := new(Param)
		cp.Name = field.Name
		cp.Type = field.Type
		cp.IsOptional = field.IsOptional
		cp.BitToTrigger = field.BitToTrigger
		cp.IsList = false

		fv, err := createFieldValidation(cp, data, true)
		if err != nil {
			return nil, err
		}

		if fv != nil {
			iterCheck := jen.For(jen.Id("_").Op(",").Id("item").Op(":=").Range().Id(direct).Block(
				fv,
			)).Comment("subitem_check")

			return jen.Add(checkLen, jen.Line(), iterCheck, jen.Line()), nil
		}

		return jen.Add(checkLen, jen.Line()), nil
	}

	zeroval := ""
	switch typ {
	case "Bool":
		zeroval = "false"
	case "long", "double", "int":
		zeroval = "0"
	case "string":
		zeroval = "\"\""
	case "bytes":
		return jen.If(jen.Len(jen.Id(direct)).Op("==").Id("0")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
		).Comment("byte_slice_check"), nil
	case "bitflags":
		return nil, nil
	case "true":
		panic("owwwooo22")
	default:
		if !insideRange {
			name = normalizeID(field.Name, false)
			name = goify(field.Name)
			direct = "e." + name
		}

		if _, ok := data.Enums[typ]; ok {
			// видимо енумы всегда uint32?
			return jen.If(
				jen.Id(direct).Op("==").Id("0"),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
			).Comment("enum_check"), nil
		}

		// если это не енум
		// дай бог чтобы у него был метод Validate()
		return jen.If(
			jen.Err().Op(":=").Id(direct+".Validate").Call(),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+name+"': %w"), jen.Id("err"))),
		).Comment("type_iter_check"), nil
	}

	// обычный билтин изи бризи
	return jen.If(jen.Id(direct).Op("==").Id(zeroval)).Block(
		jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
	).Comment("builtin_check"), nil
}

func createZeroValCheckStmt(field *Param, data *FileStructure) (*jen.Statement, error) {
	name := normalizeID(field.Name, true)
	name = goify(field.Name)
	direct := "e." + name
	typ := field.Type

	if field.IsList || typ == "bytes" {
		check := jen.Len(jen.Id(direct)).Op(">").Id("0")
		return check, nil
	}

	zeroval := ""
	switch typ {
	case "Bool":
		zeroval = "false"
	case "long", "double", "int":
		zeroval = "0"
	case "string":
		zeroval = "\"\""
	case "bitflags":
		return nil, nil
	case "true":
		// че?
		zeroval = "false"
	default:
		if _, ok := data.Enums[typ]; ok {
			return jen.Id(direct).Op("!=").Id("0"), nil
		}

		return jen.Id(direct + ".Validate").Call().Op("==").Nil(), nil
	}

	return jen.Id(direct).Op("!=").Id(zeroval), nil
}
