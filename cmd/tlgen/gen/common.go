package gen

import (
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

func (g *Generator) generateStruct(str tlparser.Object, data *internalSchema) (*jen.Statement, error) {
	fields := make([]jen.Code, 0, len(str.Parameters))

	for _, field := range str.Parameters {
		name := strcase.ToCamel(field.Name)
		typ := field.Type
		valueInsideFlag := false

		if name == "Flags" && typ == "bitflags" {
			fields = append(fields, jen.Comment("flags position"))
			continue
		}

		f := jen.Id(name)
		if field.IsVector {
			f = f.Index()
		}

		switch typ {
		case "Bool":
			f = f.Bool()
		case "long":
			f = f.Int64()
		case "double":
			f = f.Float64()
		case "int":
			f = f.Int32()
		case "string":
			f = f.String()
		case "bytes":
			f = f.Index().Byte()
		case "bitflags":
			f = f.Struct().Comment("flags param position")
		case "true":
			f = f.Bool()
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			valueInsideFlag = true
		default:
			if _, ok := data.Enums[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := data.Types[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := data.SingleInterfaceCanonical[typ]; ok {
				f = f.Id("*" + g.goify(typ))
				break
			}

			//pp.Fprintln(os.Stderr, data)
			panic("пробовали обработать '" + field.Type + "'")
		}

		tags := map[string]string{}
		if !field.IsOptional {
			tags["validate"] = "required"
		} else {
			tags["flag"] = strconv.Itoa(field.BitToTrigger)
			if valueInsideFlag {
				tags["flag"] += ",encoded_in_bitflags"
			}
		}

		f.Tag(tags)
		fields = append(fields, f)
	}

	structName := g.goify(str.Name)

	return jen.Type().Id(structName).Struct(
		fields...,
	), nil
}

func (g *Generator) generateEncodeFunc(str tlparser.Object, data *internalSchema) (*jen.Statement, error) {
	calls := make([]jen.Code, 0)
	if len(str.Parameters) > 0 {
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

	// sort.Slice(str.Parameters, func(i, j int) bool {
	// 	return str.Parameters[i].Name < str.Parameters[j].Name
	// })

	for _, field := range str.Parameters {
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

		if field.IsVector {
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

	if haveOptionalParams(str.Parameters) {
		// string это fieldname
		sortedOptionalValues := make([][]tlparser.Parameter, maxBitflag(str.Parameters)+1)
		for _, field := range str.Parameters {
			if !field.IsOptional {
				continue
			}
			if sortedOptionalValues[field.BitToTrigger] == nil {
				sortedOptionalValues[field.BitToTrigger] = make([]tlparser.Parameter, 0)
			}

			sortedOptionalValues[field.BitToTrigger] = append(sortedOptionalValues[field.BitToTrigger], tlparser.Parameter{
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

	return jen.Func().Params(jen.Id("e").Id("*" + g.goify(str.Name))).Id("Encode").Params().Index().Byte().Block(
		calls...,
	), nil
}

func (g *Generator) generateEncodeNonreflectFunc(str tlparser.Object, data *internalSchema) (*jen.Statement, error) {
	calls := make([]jen.Code, 0)
	if len(str.Parameters) > 0 {
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

	for _, field := range str.Parameters {
		valueInsideFlag := false
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
			valueInsideFlag = true
		default:
			putFuncId = "buf.PutRawBytes"
		}

		if field.IsVector {
			putFuncId = "buf.PutVector"
		}

		putFunc := jen.Null()
		switch {
		case name == "__flagsPosition":
			putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
		case valueInsideFlag:
			// не делаем ничего, значение уже заложили в флаг
		case putFuncId == "buf.PutRawBytes":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name).Dot("Encode").Call())
		case putFuncId != "":
			putFunc = jen.Id(putFuncId).Call(jen.Id("e." + name))
		default:
			panic("putFincID is empty!")
		}

		if field.IsOptional && !valueInsideFlag {
			checkStmt, err := g.createZeroValCheckStmt(field, data)
			if err != nil {
				kek := jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name)))
				fmt.Println("goodv:", kek)
				panic(err)
			}

			putFunc = jen.If(checkStmt).Block(
				putFunc,
			).Comment("check_bitflag_set")
		}

		putCalls = append(putCalls, putFunc)
	}

	if haveOptionalParams(str.Parameters) {
		// string это fieldname
		sortedOptionalValues := make([][]tlparser.Parameter, maxBitflag(str.Parameters)+1)
		for _, field := range str.Parameters {
			if !field.IsOptional {
				continue
			}
			if sortedOptionalValues[field.BitToTrigger] == nil {
				sortedOptionalValues[field.BitToTrigger] = make([]tlparser.Parameter, 0)
			}

			sortedOptionalValues[field.BitToTrigger] = append(sortedOptionalValues[field.BitToTrigger], tlparser.Parameter{
				Name:         field.Name,
				Type:         field.Type,
				IsVector:     field.IsVector,
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

				checkStmt, err := g.createZeroValCheckStmt(field, data)
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

	return jen.Func().Params(jen.Id("e").Id("*" + g.goify(str.Name))).Id("EncodeNonreflect").Params().Index().Byte().Block(
		calls...,
	), nil
}

func (g *Generator) generateMethodCallerFunc(method tlparser.Method, data *internalSchema) (*jen.Statement, error) {
	resp := createParamsStructFromMethod(method)
	maximumPositionalArguments := 0
	if haveOptionalParams(resp.Parameters) {
		maximumPositionalArguments++
	}

	funcParameters := make([]jen.Code, 0)
	methodName := g.goify(method.Name)
	typeName := methodName + "Params"

	argsAsSingleItem := false
	if len(resp.Parameters) > maximumPositionalArguments {
		argsAsSingleItem = true
		funcParameters = []jen.Code{jen.Id("params").Id("*" + typeName)}
	}

	parameters := jen.Dict{}
	for _, arg := range funcParameters {
		parameters[arg] = arg
	}

	requestStruct := jen.Op("&").Id(typeName).Values(parameters)
	if argsAsSingleItem {
		requestStruct = jen.Id("params")
	}

	assertedType := g.goify(method.Response.Type)

	if _, ok := data.SingleInterfaceCanonical[method.Response.Type]; ok {
		assertedType = "*" + assertedType
	}

	firstErrorReturn := jen.Code(jen.Nil())
	if assertedType == "Bool" {
		assertedType = "*serialize.Bool"
		//firstErrorReturn = jen.False()
	}
	if assertedType == "Long" {
		assertedType = "*serialize.Long"
		//firstErrorReturn = jen.Lit(0)
	}
	if assertedType == "Int" {
		assertedType = "*serialize.Int"
		//firstErrorReturn = jen.Lit(0)
	}

	calls := make([]jen.Code, 0)
	calls = append(calls,
		jen.List(jen.Id("data"), jen.Err()).Op(":=").Id("c.MakeRequest").Call(requestStruct),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			//jen.Return(firstErrorReturn, jen.Qual("github.com/pkg/errors", "Wrap").Call(jen.Err(), jen.Lit("sedning "+methodName))),
			jen.Return(
				firstErrorReturn,
				jen.Qual("fmt", "Errorf").Call(jen.Lit(methodName+": %w"), jen.Id("err")),
			),
		),
		jen.Line(),
		jen.List(jen.Id("resp"), jen.Id("ok")).Op(":=").Id("data").Assert(jen.Id(assertedType)),
		jen.If(jen.Op("!").Id("ok")).Block(
			//jen.Panic(jen.Lit("got invalid response type: ").Op("+").Qual("reflect", "TypeOf").Call(jen.Id("data")).Dot("String").Call()),
			jen.Err().Op(":=").Qual("fmt", "Errorf").Call(jen.Lit(methodName+": got invalid response type: %T"), jen.Id("data")),
			jen.Comment(
				jen.Return(
					firstErrorReturn,
					jen.Qual("fmt", "Errorf").Call(jen.Lit(methodName+": got invalid response type: %T"), jen.Id("data")),
				).GoString(),
			),
			jen.Panic(jen.Err()),
		),
		jen.Line(),
		jen.Return(jen.Id("resp"), jen.Nil()),
	)

	return jen.Func().Params(jen.Id("c").Id("*Client")).Id(methodName).Params(funcParameters...).Params(jen.Id(assertedType), jen.Error()).Block(
		calls...,
	), nil
}

func (g *Generator) generateStructValidatorFunc(str tlparser.Object, data *internalSchema) (*jen.Statement, error) {
	checks := make([]jen.Code, 0)
	for _, field := range str.Parameters {
		if field.IsOptional {
			_, isStruct := data.SingleInterfaceCanonical[field.Type]
			_, isIface := data.Types[field.Type]
			if !isStruct && !isIface {
				continue
			}
		}

		fv, err := g.createFieldValidation(field, data, false)
		if err != nil {
			return nil, err
		}

		if fv == nil {
			continue
		}

		checks = append(checks, fv)
		checks = append(checks, jen.Line())
	}

	checks = append(checks, jen.Return(jen.Id("nil")))
	structName := g.goify(str.Name)
	return jen.Func().Params(jen.Id("e").Id("*" + structName)).Id("Validate").Params().Params(jen.Error()).Block(
		checks...,
	), nil
}

func (g *Generator) createFieldValidation(field tlparser.Parameter, data *internalSchema, insideRange bool) (jen.Code, error) {
	name := g.goify(field.Name)
	direct := "e." + name
	if insideRange {
		direct = "item"
	}
	typ := field.Type

	if field.IsVector {
		name = g.goify(field.Name)
		direct := "e." + name

		checkLen := jen.If(jen.Len(jen.Id(direct)).Op("==").Id("0")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
		).Comment("slice_len_check")

		cp := tlparser.Parameter{
			Name:         field.Name,
			Type:         field.Type,
			IsOptional:   field.IsOptional,
			BitToTrigger: field.BitToTrigger,
			IsVector:     false,
		}

		fv, err := g.createFieldValidation(cp, data, true)
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
		zeroval = "false"
		//panic("owwwooo22")
	default:
		if !insideRange {
			name = g.goify(field.Name)
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

		_, isStruct := data.SingleInterfaceCanonical[field.Type]
		_, isIface := data.Types[field.Type]
		if isStruct || isIface {
			typComment := "struct"
			if isIface {
				typComment = "interface"
			}

			if field.IsOptional {
				check := jen.If(
					jen.Id(direct).Op("!=").Nil(),
				).Block(
					jen.If(
						jen.Err().Op(":=").Id(direct+".Validate").Call(),
						jen.Err().Op("!=").Nil(),
					).Block(
						jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+name+"': %w"), jen.Id("err"))),
					),
				).Comment("optional_" + typComment + "_valid_check")

				return check, nil
			}

			nilcheck := jen.If(
				jen.Id(direct).Op("==").Nil(),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' required"))),
			).Comment("required_" + typComment + "_nil_check")

			validcheck := jen.If(
				jen.Err().Op(":=").Id(direct+".Validate").Call(),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+name+"': %w"), jen.Id("err"))),
			).Comment("required_" + typComment + "_valid_check")

			return nilcheck.Line().Add(validcheck), nil
		}

		// if isIface {
		// 	return jen.If(
		// 		jen.Id(direct).Op("!=").Nil(),
		// 	).Block(
		// 		jen.If(
		// 			jen.Err().Op(":=").Id(direct+".Validate").Call(),
		// 			jen.Err().Op("!=").Nil(),
		// 		).Block(
		// 			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+name+"': %w"), jen.Id("err"))),
		// 		),
		// 	).Comment("interface_valid_check"), nil
		// }
		panic("wat")

	}

	// обычный билтин изи бризи
	return jen.If(jen.Id(direct).Op("==").Id(zeroval)).Block(
		jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + name + "' is not set"))),
	).Comment("builtin_check"), nil
}

func (g *Generator) createZeroValCheckStmt(field tlparser.Parameter, data *internalSchema) (*jen.Statement, error) {
	name := g.goify(field.Name)
	direct := "e." + name
	typ := field.Type

	if field.IsVector || typ == "bytes" {
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

		check := jen.Id(direct).Op("!=").Nil().Op("&&").
			Id(direct + ".Validate").Call().Op("==").Nil()
		return check, nil
	}

	return jen.Id(direct).Op("!=").Id(zeroval), nil
}
