package gen

import (
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/xelaj/mtproto/cmd/tlgen/tlparser"
)

func (g *Generator) generateStruct(str tlparser.Object) (*jen.Statement, error) {
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
			if _, ok := g.schema.Enums[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := g.schema.Types[typ]; ok {
				f = f.Id(g.goify(typ))
				break
			}
			if _, ok := g.schema.SingleInterfaceCanonical[typ]; ok {
				f = f.Id("*" + g.goify(typ))
				break
			}

			//pp.Fprintln(os.Stderr, g.schema)
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

func (g *Generator) generateEncodeFunc(str tlparser.Object) (*jen.Statement, error) {
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

	for _, field := range str.Parameters {
		valueInsideFlag := false
		typ := field.Type
		name := strcase.ToCamel(field.Name)
		if name == "Flags" && typ == "bitflags" {
			name = "__flagsPosition"
		}

		putFuncID := ""
		switch typ {
		case "Bool":
			putFuncID = "buf.PutBool"
		case "long":
			putFuncID = "buf.PutLong"
		case "double":
			putFuncID = "buf.PutDouble"
		case "int":
			putFuncID = "buf.PutInt"
		case "string":
			putFuncID = "buf.PutString"
		case "bytes":
			putFuncID = "buf.PutMessage"
		case "true":
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			valueInsideFlag = true
		default:
			putFuncID = "buf.PutRawBytes"
		}

		if field.IsVector {
			putFuncID = "buf.PutVector"
		}

		putFunc := jen.Null()
		switch {
		case name == "__flagsPosition":
			putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
		case valueInsideFlag:
			// не делаем ничего, значение уже заложили в флаг
		case putFuncID == "buf.PutRawBytes":
			putFunc = jen.Id(putFuncID).Call(jen.Id("e." + name).Dot("Encode").Call())
		case putFuncID != "":
			putFunc = jen.Id(putFuncID).Call(jen.Id("e." + name))
		default:
			panic("putFincID is empty!")
		}

		if field.IsOptional && !valueInsideFlag {
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

func (g *Generator) generateEncodeNonreflectFunc(str tlparser.Object) (*jen.Statement, error) {
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

		putFuncID := ""
		switch typ {
		case "Bool":
			putFuncID = "buf.PutBool"
		case "long":
			putFuncID = "buf.PutLong"
		case "double":
			putFuncID = "buf.PutDouble"
		case "int":
			putFuncID = "buf.PutInt"
		case "string":
			putFuncID = "buf.PutString"
		case "bytes":
			putFuncID = "buf.PutMessage"
		case "true":
			//! ИСКЛЮЧЕНИЕ БЛЯТЬ! ИСКЛЮЧЕНИЕ!!!
			//! если в опциональном флаге указан true, то это значение true уходит в битфлаги и его типа десериализовать не надо!!! ебать!!! ЕБАТЬ!!!
			valueInsideFlag = true
		default:
			putFuncID = "buf.PutRawBytes"
		}

		if field.IsVector {
			putFuncID = "buf.PutVector"
		}

		putFunc := jen.Null()
		switch {
		case name == "__flagsPosition":
			putFunc = jen.Id("buf.PutUint").Call(jen.Id("flag"))
		case valueInsideFlag:
			// не делаем ничего, значение уже заложили в флаг
		case putFuncID == "buf.PutRawBytes":
			putFunc = jen.Id(putFuncID).Call(jen.Id("e." + name).Dot("Encode").Call())
		case putFuncID != "":
			putFunc = jen.Id(putFuncID).Call(jen.Id("e." + name))
		default:
			panic("putFincID is empty!")
		}

		if field.IsOptional && !valueInsideFlag {
			checkStmt, err := g.createZeroValueCheckStmt(field)
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

				checkStmt, err := g.createZeroValueCheckStmt(field)
				if err != nil {
					kek := jen.Op("!").Qual("github.com/vikyd/zero", "IsZeroVal").Call(jen.Id("e." + strcase.ToCamel(field.Name))).GoString()
					return nil, fmt.Errorf("bad zero value check stmt: %w (good example: %s)", err, kek)
				}

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

func (g *Generator) generateMethodCallerFunc(method tlparser.Method) (*jen.Statement, error) {
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

	if _, ok := g.schema.SingleInterfaceCanonical[method.Response.Type]; ok {
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

func (g *Generator) generateStructValidatorFunc(str tlparser.Object) (*jen.Statement, error) {
	checks := make([]jen.Code, 0)
	for _, field := range str.Parameters {
		if field.IsOptional {
			_, isStruct := g.schema.SingleInterfaceCanonical[field.Type]
			_, isIface := g.schema.Types[field.Type]
			isBuiltin := !isStruct && !isIface

			// если поле - опциональный билтин
			// не проверяем его
			if isBuiltin {
				continue
			}

			// если поле - опциональный НЕ билтин (т.е. структура или интерфейс)
			// нужно сначала проверить его на nil
			//
			// если оно nil - не валидируем его т.к. оно опционально
			// если оно НЕ nil - выполняем валидацию
		}

		validateStmt, err := g.createFieldValidation(field, false)
		if err != nil {
			return nil, err
		}

		if validateStmt == nil {
			continue
		}

		checks = append(checks, validateStmt)
		checks = append(checks, jen.Line())
	}

	// если все валидации прошли, отдаем nil
	checks = append(checks, jen.Return(jen.Nil()))

	return jen.Func().Params(jen.Id("e").Id("*" + g.goify(str.Name))).Id("Validate").Params().Params(jen.Error()).Block(
		checks...,
	), nil
}

func (g *Generator) createFieldValidation(field tlparser.Parameter, insideRange bool) (jen.Code, error) {
	// название проверяемого филда без префикса структуры
	// (нужно для ошибок)
	goname := g.goify(field.Name)

	// филд с префиксом структуры
	// юзается непосредственно для обращения к полю
	direct := "e." + goname

	// если проверка идет внтури итерации
	if insideRange {
		direct = "item"
	}

	if field.IsVector {
		// делаем проверку на длину слайса
		checkLen := jen.If(jen.Len(jen.Id(direct)).Op("==").Id("0")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + goname + "' is not set"))),
		).Comment("slice_len_check")

		// копируем текущий параметр, но без вектора
		cp := tlparser.Parameter{
			Name:         field.Name,
			Type:         field.Type,
			IsOptional:   field.IsOptional,
			BitToTrigger: field.BitToTrigger,
			IsVector:     false,
		}

		// мутим для него валидатор, чтобы юзать внутри итерации
		validateStmt, err := g.createFieldValidation(cp, true)
		if err != nil {
			return nil, err
		}

		if validateStmt != nil {
			iterCheck := jen.For(jen.Id("_").Op(",").Id("item").Op(":=").Range().Id(direct).Block(
				validateStmt,
			)).Comment("subitem_check")

			return jen.Add(checkLen, jen.Line(), iterCheck, jen.Line()), nil
		}

		return jen.Add(checkLen, jen.Line()), nil
	}

	zeroval := ""
	switch field.Type {
	case "Bool":
		zeroval = "false"
	case "long", "double", "int":
		zeroval = "0"
	case "string":
		zeroval = `""`
	case "bytes":
		return jen.If(jen.Len(jen.Id(direct)).Op("==").Id("0")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + goname + "' is not set"))),
		).Comment("byte_slice_check"), nil
	case "bitflags":
		return nil, nil
	case "true":
		zeroval = "false"
	default:
		if _, ok := g.schema.Enums[field.Type]; ok {
			// видимо енумы всегда uint32?
			return jen.If(
				jen.Id(direct).Op("==").Id("0"),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + goname + "' is not set"))),
			).Comment("enum_check"), nil
		}

		_, isStruct := g.schema.SingleInterfaceCanonical[field.Type]
		_, isIface := g.schema.Types[field.Type]
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
						jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+goname+"': %w"), jen.Id("err"))),
					),
				).Comment("optional_" + typComment + "_valid_check")

				return check, nil
			}

			nilcheck := jen.If(
				jen.Id(direct).Op("==").Nil(),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + goname + "' required"))),
			).Comment("required_" + typComment + "_nil_check")

			validcheck := jen.If(
				jen.Err().Op(":=").Id(direct+".Validate").Call(),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '"+goname+"': %w"), jen.Id("err"))),
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

	// обычный билтин
	return jen.If(jen.Id(direct).Op("==").Id(zeroval)).Block(
		jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("field '" + goname + "' is not set"))),
	).Comment("builtin_check"), nil
}

// отдает проверку на zero value
// abc != 0
// foo != ""
// len(slice) > 0
func (g *Generator) createZeroValueCheckStmt(field tlparser.Parameter) (*jen.Statement, error) {
	name := g.goify(field.Name)
	direct := "e." + name

	// если вектор или байты, просто проверяем длину слайса
	if field.IsVector || field.Type == "bytes" {
		check := jen.Len(jen.Id(direct)).Op(">").Id("0")
		return check, nil
	}

	zeroval := ""
	switch field.Type {
	case "Bool":
		zeroval = "false"
	case "long", "double", "int":
		zeroval = "0"
	case "string":
		zeroval = "\"\""
	case "bitflags":
		return nil, nil
	case "true":
		zeroval = "false"
	default:
		// енум
		if _, ok := g.schema.Enums[field.Type]; ok {
			return jen.Id(direct).Op("!=").Id("0"), nil
		}

		// структура или интерфейс
		// (структуры всегда с указателем, хз почему)
		check := jen.Id(direct).Op("!=").Nil().Op("&&").
			Id(direct + ".Validate").Call().Op("==").Nil()
		return check, nil
	}

	// билтин
	return jen.Id(direct).Op("!=").Id(zeroval), nil
}
