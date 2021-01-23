package gen

import (
	"sort"

	"github.com/dave/jennifer/jen"
	
	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

var maximumPositionalArguments = 5

func (g *Generator) generateMethods(f *jen.File) {
	sort.Slice(g.schema.Methods, func(i, j int) bool {
		return g.schema.Methods[i].Name < g.schema.Methods[j].Name
	})

	for _, method := range g.schema.Methods {
		f.Add(g.generateStructTypeAndMethods(tlparser.Object{
			Name:       method.Name + "Params",
			Comment:    method.Comment,
			CRC:        method.CRC,
			Parameters: method.Parameters,
		}, nil))
		f.Line()
		if method.Comment != "" {
			f.Comment(method.Comment)
		}
		f.Add(g.generateMethodFunction(&method))
		f.Line()
	}

	//	sort.Strings(keys)
	//
	//	for _, i := range keys {
	//		structs := g.schema.Types[i]
	//
	//		sort.Slice(structs, func(i, j int) bool {
	//			return structs[i].Name < structs[j].Name
	//		})
	//
	//		for _, _type := range structs {
	//			f.Add(g.generateStructTypeAndMethods(_type, []string{goify(i, true)}))
	//			f.Line()
	//		}
	//	}
}

func (g *Generator) generateMethodFunction(obj *tlparser.Method) jen.Code {
	resp := g.typeIdFromSchemaType(obj.Response.Type)
	if obj.Response.IsList {
		resp = jen.Index().Add(resp)
	}

	// еще одно злоебучее исключение. проблема в том, что bool это вот как бы и объект, да вот как бы и нет
	// трабла только в том, что нельзя просто так взять, и получить bool из MakeRequest. так что
	// возвращаем tl.Bool
	if obj.Response.Type == "Bool" {
		resp = jen.Op("*").Qual(tlPackagePath, "PseudoBool")
	}

	responses := []jen.Code{resp, jen.Error()}

	//*	data, err := c.MakeRequest(params)
	//*	if err != nil {
	//*		return nil, errors.Wrap(err, "sedning AuthSendCode")
	//*	}
	//*
	//*	resp, ok := data.(*AuthSentCode)
	//*	if !ok {
	//*		panic("got invalid response type: " + reflect.TypeOf(data).String())
	//*	}
	//*
	//*	return resp, nil
	method := jen.Func().Params(jen.Id("c").Op("*").Id("Client")).Id(goify(obj.Name, true)).Params(g.generateArgumentsForMethod(obj)...).Params(responses...).Block(
		jen.List(jen.Id("responseData"), jen.Id("err")).Op(":=").Id("c").Dot("MakeRequest").Call(g.generateMethodArgumentForMakingRequest(obj)),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual(errorsPackagePath, "Wrap").Call(jen.Err(), jen.Lit("sending "+goify(obj.Name, true)))),
		),
		jen.Line(),
		jen.List(jen.Id("resp"), jen.Id("ok")).Op(":=").Id("responseData").Assert(resp),
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Panic(jen.Lit("got invalid response type: ").Op("+").Qual("reflect", "TypeOf").Call(jen.Id("responseData")).Dot("String").Call()),
		),
		jen.Return(jen.Id("resp"), jen.Nil()),
	)

	return method
}

func (g *Generator) generateArgumentsForMethod(obj *tlparser.Method) []jen.Code {
	if len(obj.Parameters) == 0 {
		return []jen.Code{}
	}
	if len(obj.Parameters) > maximumPositionalArguments {
		return []jen.Code{jen.Id("params").Op("*").Id(goify(obj.Name, true) + "Params")}
	}

	items := make([]jen.Code, 0)

	for i, p := range obj.Parameters {
		item := jen.Id(goify(p.Name, false))
		if i == len(obj.Parameters)-1 || p.Type != obj.Parameters[i+1].Type || p.IsVector != obj.Parameters[i+1].IsVector {
			if p.Type == "bitflags" {
				continue // ну а зачем?
			}

			if p.IsVector {
				item = item.Add(jen.Index(), g.typeIdFromSchemaType(p.Type))
			} else {
				item = item.Add(g.typeIdFromSchemaType(p.Type))
			}
		}

		items = append(items, item)
	}
	return items
}

func (g *Generator) generateMethodArgumentForMakingRequest(obj *tlparser.Method) *jen.Statement {
	if len(obj.Parameters) > maximumPositionalArguments {
		return jen.Id("params")
	}

	dict := jen.Dict{}
	for _, p := range obj.Parameters {
		if p.Type == "bitflags" {
			continue // ну а зачем?
		}

		dict[jen.Id(goify(p.Name, true))] = jen.Id(goify(p.Name, false))
	}

	return jen.Op("&").Id(goify(obj.Name, true) + "Params").Values(dict)
}
