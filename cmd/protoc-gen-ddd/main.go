package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"flag"
	"os"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/dave/jennifer/jen"
	// "github.com/xoe-labs/ddd-gen/pkg/gen_domain/generator"
)

// fieldGoType returns the Go type used for a field.
//
// If it returns pointer=true, the struct field is a pointer to the type.
func fieldGoType(g *protogen.GeneratedFile, field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}

	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := fieldGoType(g, field.Message.Fields[0])
		valType, _ := fieldGoType(g, field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}
	return goType, pointer
}

func main() {
	// Protoc passes pluginpb.CodeGeneratorRequest in via stdin
	// marshalled with Protobuf
	input, _ := ioutil.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	proto.Unmarshal(input, &req)

	var flags flag.FlagSet
	Entity := flags.String("entity", "", "")
	opts := &protogen.Options{
	  ParamFunc: flags.Set,
	}
	if Entity == nil {
		panic("necessary to define 'entity' option (eg. --ddd_out=entity=Account:.)")
	}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}

	// Protoc passes a slice of File structs for us to process
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		// Specify the output filename
		filename := file.GeneratedFilenamePrefix + ".model.go"
		genFile := plugin.NewGeneratedFile(filename, file.GoImportPath)

		// 1. Initialise a buffer to hold the generated code
		buf := &bytes.Buffer{}

		// 2. Write the package name
		f := jen.NewFilePathName(string(file.GoImportPath), string(file.GoPackageName))

		var visited = make(map[string]bool)

		f.Commentf("%s is a domain model", *Entity)
		f.Type().Id(*Entity).StructFunc(func(g *jen.Group) {
			for _, msg := range file.Messages {
				for _, fld := range msg.Fields {
					if _, ok := visited[fld.GoName]; !ok {
						goType, pointer := fieldGoType(genFile, fld)
						if pointer {
							g.Id(lowerFirst(fld.GoName)).Op("*").Id(goType)
						} else {
							g.Id(lowerFirst(fld.GoName)).Id(goType)
						}
						visited[fld.GoName] = true
					}
				}
			}
		})

		var New = "New"

		f.Commentf("%s%s constructs an empty %s", New, *Entity, *Entity)
		f.Func().Id(New + *Entity).Params().Op("*").Id(*Entity).Block(
			jen.Return().Op("&").Id(*Entity).Values(),
		)

		// var String = "String"

		// f.Commentf("%s implements fmt.Stringer for %s", String, Entity)
		// f.Func().Id(String).Params().String().Block(
		// 	jen.Return().Lit(""),
		// )

		var Apply = "Apply"

		f.Commentf("%s implements app.??? for %s", Apply, *Entity)
		f.Func().Params(
			jen.Id(firstLower(*Entity)).Op("*").Id(*Entity),
		).Id(Apply).Params(
			jen.Id("fact").Interface(),
		).Params(
			jen.Id("success").Bool(),
		).Block(
			jen.Switch(
				jen.Id("f").Op(":=").Id("fact").Assert(jen.Type()),
			).BlockFunc(func(g *jen.Group) {
				for _, msg := range file.Messages {
					g.Case(
						jen.Op("*").Id(msg.GoIdent.GoName),
					).BlockFunc(func(g *jen.Group) {
						for _, fld := range msg.Fields {
							g.Id(
								firstLower(*Entity),
							).Dot(
								lowerFirst(fld.GoName),
							).Op("=").Id("f").Dot(
								fld.GoName,
							)
						}
						g.Return().True()
					})
				}
				g.Default().Block(
					jen.Return().False(),
				)
			}),
		)

		err := f.Render(buf)
		if err != nil {
			panic(err)
		}

		// 5. Pass the data from our buffer to the plugin file struct
		_, err = genFile.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// Write the response to stdout, to be picked up by protoc
	fmt.Fprintf(os.Stdout, string(out))
}

func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func firstLower(str string) string {
	return "r"
}
