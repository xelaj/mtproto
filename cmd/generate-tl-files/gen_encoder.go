package main

import (
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func generateEncodeFunc(str *StructObject, data *FileStructure) (*jen.Statement, error) {
	calls := make([]jen.Code, 0)
	if len(str.Fields) > 0 {
		// calls = append(calls,
		// 	jen.Id("err").Op(":=").Qual("github.com/go-playground/validator", "New").Call().Dot("Struct").Call(jen.Id("e")),
		// 	jen.Qual("github.com/xelaj/go-dry", "PanicIfErr").Call(jen.Id("err")),
		// 	jen.Line(),
		// )
		calls = append(calls,
			jen.If(
				jen.Err().Op(":=").Qual("github.com/go-playground/validator", "New").Call().Dot("Struct").Call(jen.Id("e")),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Comment(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validator: %w"), jen.Id("err"))).GoString(),
				),
				jen.Panic(jen.Id("err")),
			),
		)
	}

	putCalls := make([]jen.Code, 0)

	sort.Slice(str.Fields, func(i, j int) bool {
		return str.Fields[i].Name < str.Fields[j].Name
	})
	for _, field := range str.Fields {
		ЗНАЧЕНИЕ_В_ФЛАГЕ := false
		typ := field.Type
		name := strcase.ToCamel(field.Name)
		if name == "Flags" && typ == "bitflags" {
			name = "__flagsPosition"
		}

		putFuncId := ""
		switch typ {
		case "Bool":
			putFuncId = "buf.PutBool"
		case "long":
			putFuncId = "buf.PutLong"
		case "double":
			putFuncId = "buf.PutDouble"
		case "int":
			putFuncId = "buf.PutInt"
		case "string":
			putFuncId = "buf.PutString"
		case "bytes":
			putFuncId = "buf.PutMessage"
		case "true":
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			ЗНАЧЕНИЕ_В_ФЛАГЕ = true
		default:
			putFuncId = "buf.PutRawBytes"
		}

		if field.IsList {
			putFuncId = "buf.PutVector"
		}

		putFunc := jen.Null()
		switch {
		case name == "__flagsPosition":
			putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
		case ЗНАЧЕНИЕ_В_ФЛАГЕ:
			// не делаем ничего, значение уже заложили в флаг
		case putFuncId == "buf.PutRawBytes":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name).Dot("Encode").Call())
		case putFuncId != "":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name))
		default:
			panic("putFincID is empty!")
		}

		if field.IsOptional && !ЗНАЧЕНИЕ_В_ФЛАГЕ {
			putFunc = jen.If(jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name)))).Block(
				putFunc,
			)
		}

		putCalls = append(putCalls, putFunc)
	}

	if str.HaveOptionalArgs() {
		// string это fieldname
		sortedOptionalValues := make([][]*Param, str.MaxFlagBit()+1)
		for _, field := range str.Fields {
			if !field.IsOptional {
				continue
			}
			if sortedOptionalValues[field.BitToTrigger] == nil {
				sortedOptionalValues[field.BitToTrigger] = make([]*Param, 0)
			}

			sortedOptionalValues[field.BitToTrigger] = append(sortedOptionalValues[field.BitToTrigger], &Param{
				Name: field.Name,
				Type: field.Type,
			})
		}

		flagchecks := make([]jen.Code, 0, len(sortedOptionalValues))
		for i, fields := range sortedOptionalValues {
			if len(fields) == 0 {
				continue
			}

			statements := jen.Null()
			for j, field := range fields {
				if j != 0 {
					statements.Add(jen.Op("||"))
				}

				//? zero.IsZeroVal(e.Fieldname)
				statements.Add(jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name))))
			}

			flagchecks = append(flagchecks, jen.If(statements).Block(
				//? flag |= 1 << n
				jen.Id("flag").Op("|=").Lit(1).Op("<<").Lit(i),
			))
		}

		calls = append(calls, jen.Var().Id("flag").Uint32())
		calls = append(calls,
			flagchecks...,
		)
	}

	calls = append(calls,
		jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto/serialize", "NewEncoder").Call(),
		jen.Id("buf.PutUint").Call(jen.Id("e.CRC").Call()),
	)

	calls = append(calls, putCalls...)
	calls = append(calls, jen.Return(jen.Id("buf.Result").Call()))

	typeName := normalizeID(str.Name, false)
	return jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("Encode").Params().Index().Byte().Block(
		calls...,
	), nil
}

func generateEncodeNonreflectFunc(str *StructObject, data *FileStructure) (*jen.Statement, error) {
	calls := make([]jen.Code, 0)
	if len(str.Fields) > 0 {
		calls = append(calls,
			jen.If(
				jen.Err().Op(":=").Id("e.Validate").Call(),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Comment(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validator: %w"), jen.Id("err"))).GoString(),
				),
				jen.Panic(jen.Id("err")),
			),
		)
	}

	putCalls := make([]jen.Code, 0)

	sort.Slice(str.Fields, func(i, j int) bool {
		return str.Fields[i].Name < str.Fields[j].Name
	})
	for _, field := range str.Fields {
		ЗНАЧЕНИЕ_В_ФЛАГЕ := false
		typ := field.Type
		name := strcase.ToCamel(field.Name)
		if name == "Flags" && typ == "bitflags" {
			name = "__flagsPosition"
		}

		putFuncId := ""
		switch typ {
		case "Bool":
			putFuncId = "buf.PutBool"
		case "long":
			putFuncId = "buf.PutLong"
		case "double":
			putFuncId = "buf.PutDouble"
		case "int":
			putFuncId = "buf.PutInt"
		case "string":
			putFuncId = "buf.PutString"
		case "bytes":
			putFuncId = "buf.PutMessage"
		case "true":
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			ЗНАЧЕНИЕ_В_ФЛАГЕ = true
		default:
			putFuncId = "buf.PutRawBytes"
		}

		if field.IsList {
			putFuncId = "buf.PutVector"
		}

		putFunc := jen.Null()
		switch {
		case name == "__flagsPosition":
			putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
		case ЗНАЧЕНИЕ_В_ФЛАГЕ:
			// не делаем ничего, значение уже заложили в флаг
		case putFuncId == "buf.PutRawBytes":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name).Dot("Encode").Call())
		case putFuncId != "":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name))
		default:
			panic("putFincID is empty!")
		}

		if field.IsOptional && !ЗНАЧЕНИЕ_В_ФЛАГЕ {
			checkStmt, err := createZeroValCheckStmt(field, data)
			if err != nil {
				kek := jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name)))
				fmt.Println("goodv:", kek)
				panic(err)
			}

			putFunc = jen.If(checkStmt).Block(
				putFunc,
			)
		}

		putCalls = append(putCalls, putFunc)
	}

	if str.HaveOptionalArgs() {
		// string это fieldname
		sortedOptionalValues := make([][]*Param, str.MaxFlagBit()+1)
		for _, field := range str.Fields {
			if !field.IsOptional {
				continue
			}
			if sortedOptionalValues[field.BitToTrigger] == nil {
				sortedOptionalValues[field.BitToTrigger] = make([]*Param, 0)
			}

			sortedOptionalValues[field.BitToTrigger] = append(sortedOptionalValues[field.BitToTrigger], &Param{
				Name:         field.Name,
				Type:         field.Type,
				IsList:       field.IsList,
				IsOptional:   field.IsOptional,
				BitToTrigger: field.BitToTrigger,
			})
		}

		flagchecks := make([]jen.Code, 0, len(sortedOptionalValues))
		for i, fields := range sortedOptionalValues {
			if len(fields) == 0 {
				continue
			}

			statements := jen.Null()
			for j, field := range fields {
				if j != 0 {
					statements.Add(jen.Op("||"))
				}

				checkStmt, err := createZeroValCheckStmt(field, data)
				if err != nil {
					kek := jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name))).GoString()
					fmt.Println("goodv:", kek)
					panic(err)
				}

				// if field.Name == "documents" && field.Type == "Document" {
				// 	pp.Println(str)
				// 	panic("done")
				// }

				statements.Add(checkStmt)
			}

			flagchecks = append(flagchecks, jen.If(statements).Block(
				//? flag |= 1 << n
				jen.Id("flag").Op("|=").Lit(1).Op("<<").Lit(i),
			))
		}

		calls = append(calls, jen.Var().Id("flag").Uint32())
		calls = append(calls,
			flagchecks...,
		)
	}

	calls = append(calls,
		jen.Id("buf").Op(":=").Qual("github.com/xelaj/mtproto/serialize", "NewEncoder").Call(),
		jen.Id("buf.PutUint").Call(jen.Id("e.CRC").Call()),
	)

	calls = append(calls, putCalls...)
	calls = append(calls, jen.Return(jen.Id("buf.Result").Call()))

	typeName := normalizeID(str.Name, false)
	return jen.Func().Params(jen.Id("e").Id("*" + typeName)).Id("EncodeNonreflect").Params().Index().Byte().Block(
		calls...,
	), nil
}
