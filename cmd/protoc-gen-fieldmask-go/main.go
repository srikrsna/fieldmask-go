package main

import (
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	generatedFilenameExtension = ".fieldmask.go"
	generatedPackageSuffix     = "fieldmask"

	fieldmaskPkg = protogen.GoImportPath("github.com/srikrsna/fieldmask-go")
)

var (
	wktSet = map[protoreflect.FullName]bool{
		(new(structpb.NullValue)).Descriptor().FullName():                  true,
		(&structpb.Struct{}).ProtoReflect().Descriptor().FullName():        true,
		(&structpb.ListValue{}).ProtoReflect().Descriptor().FullName():     true,
		(&structpb.Value{}).ProtoReflect().Descriptor().FullName():         true,
		(&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName():  true,
		(&timestamppb.Timestamp{}).ProtoReflect().Descriptor().FullName():  true,
		(&durationpb.Duration{}).ProtoReflect().Descriptor().FullName():    true,
		(&anypb.Any{}).ProtoReflect().Descriptor().FullName():              true,
		(&emptypb.Empty{}).ProtoReflect().Descriptor().FullName():          true,
		(&wrapperspb.BoolValue{}).ProtoReflect().Descriptor().FullName():   true,
		(&wrapperspb.StringValue{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.Int32Value{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.Int64Value{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.FloatValue{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor().FullName(): true,
	}
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range plugin.Files {
			if file.Generate {
				gen(plugin, file)
			}
		}
		return nil
	})
}

func gen(plugin *protogen.Plugin, file *protogen.File) {
	if len(file.Messages) == 0 {
		return
	}
	file.GoPackageName += generatedPackageSuffix
	file.GoPackageName = protogen.GoPackageName(strings.TrimPrefix(string(file.GoPackageName), "_"))

	dir := path.Dir(file.GeneratedFilenamePrefix)
	base := path.Base(file.GeneratedFilenamePrefix)
	file.GeneratedFilenamePrefix = path.Join(
		dir,
		string(file.GoPackageName),
		base,
	)
	generatedFile := plugin.NewGeneratedFile(
		file.GeneratedFilenamePrefix+generatedFilenameExtension,
		protogen.GoImportPath(path.Join(
			string(file.GoImportPath),
			string(file.GoPackageName),
		)),
	)
	genFileHeader(generatedFile, file)
	for _, message := range file.Messages {
		genMessage(generatedFile, file, message)
	}
}

func genFileHeader(g *protogen.GeneratedFile, file *protogen.File) {
	g.P("// Code generated by ", path.Base(os.Args[0]), ". DO NOT EDIT.")
	g.P("//")
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("//", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// Source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
}

func genMessage(g *protogen.GeneratedFile, file *protogen.File, message *protogen.Message) {
	g.Annotate(message.GoIdent.GoName, message.Location)
	g.P("var ", message.GoIdent.GoName, "Mask = ", message.GoIdent.GoName, "(\"\")")
	g.P()
	g.P("// ", message.GoIdent.GoName, " is the field mask for ", message.Desc.Name())
	g.P("type ", message.GoIdent.GoName, " string")
	for _, field := range message.Fields {
		returnType := "string"
		if !field.Desc.IsMap() && field.Message != nil && !wktSet[field.Desc.Message().FullName()] {
			if file.GoImportPath == field.Message.GoIdent.GoImportPath {
				returnType = field.Message.GoIdent.GoName
			} else {
				ident := field.Message.GoIdent
				ident.GoImportPath = protogen.GoImportPath(path.Join(
					string(ident.GoImportPath),
					string(path.Base(string(ident.GoImportPath))+generatedPackageSuffix),
				))
				returnType = g.QualifiedGoIdent(ident)
			}
		}
		if field.Desc.IsList() {
			returnType = g.QualifiedGoIdent(fieldmaskPkg.Ident("Slice[" + returnType + "]"))
		}
		g.P("func (f ", message.GoIdent.GoName, ")", field.GoName, "() ", returnType, "{")
		g.P("if f == \"\" {")
		g.P("return \"", field.Desc.Name(), "\"")
		g.P("}")
		g.P("return ", returnType, "(", "string(f) + \".\" + ", "\"", field.Desc.Name(), "\"", ")")
		g.P("}")
		g.P("")
	}
	for _, message := range message.Messages {
		if message.Desc.IsMapEntry() {
			continue
		}
		genMessage(g, file, message)
	}
}
