package codegen

import (
	"aloisdeniel/uipack"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type FlutterCodeGen struct {
	Builder strings.Builder
}

func (g *FlutterCodeGen) GenerateDefinitions(metadata *uipack.BundleMetadata, bundles *[]uipack.Bundle) string {
	g.Builder.WriteString("import 'package:flutter/widgets.dart';\n\n")

	for _, bundle := range *bundles {
		identifier := fmt.Sprintf("%x", bundle.Variant)
		g.Builder.WriteString(fmt.Sprintf("import 'bundle_%s.g.dart';\n", identifier))
	}

	g.Builder.WriteString("\n")

	g.generateModeTypeDefinitions(metadata)
	g.generateVariantTypeDefinition(metadata, bundles)
	g.generateBundleDataTypeDefinition(metadata)
	g.generateBundleExtension(metadata)
	return g.Builder.String()
}

func (g *FlutterCodeGen) generateModeTypeDefinitions(metadata *uipack.BundleMetadata) {
	for _, mode := range metadata.Modes {
		g.generateModeEnumDefinition(&mode)
	}
}
func (g *FlutterCodeGen) generateModeEnumDefinition(mode *uipack.ModeMetadata) {
	g.Builder.WriteString("/// Index : " + fmt.Sprintf("%x", mode.Identifier) + "\n")
	g.Builder.WriteString("enum ")
	g.Builder.WriteString(dartType(mode.Name))
	g.Builder.WriteString("Mode {\n")
	for _, v := range mode.Variants {
		g.Builder.WriteString(dartField(v.Name))
		g.Builder.WriteString(",\n")
	}
	g.Builder.WriteString("}\n")
}

func (g *FlutterCodeGen) generateCustomTypeDefinitions(metadata *uipack.BundleMetadata) {
	for _, typedef := range metadata.Types {
		g.Builder.WriteString(fmt.Sprintf("class %s {\n", dartType(typedef.Name)))

		g.Builder.WriteString(fmt.Sprintf("const %s({\n", dartType(typedef.Name)))
		for _, v := range typedef.Properties {
			g.Builder.WriteString(fmt.Sprintf("this.%s,\n", dartField(v.Name)))
		}

		g.Builder.WriteString("});\n")

		for _, v := range typedef.Properties {
			t := generateDartType(metadata, v.Type)
			g.Builder.WriteString(fmt.Sprintf("final %s %s;\n", t, dartField(v.Name)))

		}
		g.Builder.WriteString("}\n\n")
	}
}

func (g *FlutterCodeGen) generateVariantTypeDefinition(metadata *uipack.BundleMetadata, bundles *[]uipack.Bundle) {
	g.Builder.WriteString("typedef Variant = ({")
	for _, mode := range metadata.Modes {
		g.Builder.WriteString(dartType(mode.Name))
		g.Builder.WriteString("Mode ")
		g.Builder.WriteString(dartField(mode.Name))
		g.Builder.WriteString(",")
	}
	g.Builder.WriteString("});")

	g.Builder.WriteString("extension VariantExtension on Variant {")

	// Identifier
	g.Builder.WriteString("int get identifier {")
	g.Builder.WriteString("var result = 0;")
	for i, mode := range metadata.Modes {
		g.Builder.WriteString("result |= ")
		g.Builder.WriteString(dartField(mode.Name))
		g.Builder.WriteString(fmt.Sprintf(".index << %d;", i*4))
	}
	g.Builder.WriteString("return result;")
	g.Builder.WriteString("}")

	// Bundle
	g.Builder.WriteString("Bundle? get bundle {")
	g.Builder.WriteString("switch (identifier) {")
	for _, bundle := range *bundles {
		g.Builder.WriteString("case 0x")
		g.Builder.WriteString(fmt.Sprintf("%x", bundle.Variant))
		g.Builder.WriteString(": return bundle")
		g.Builder.WriteString(fmt.Sprintf("%x", bundle.Variant))
		g.Builder.WriteString("();")
	}
	g.Builder.WriteString("default: return null;")

	g.Builder.WriteString("}")
	g.Builder.WriteString("}")

	g.Builder.WriteString("}")

	// Custom types
	g.generateCustomTypeDefinitions(metadata)
}

func (g *FlutterCodeGen) generateBundleExtension(metadata *uipack.BundleMetadata) {

	g.Builder.WriteString("extension BundleExtension on Bundle {")

	// Metadata
	g.Builder.WriteString("List<(String, dynamic)> get metadata {")
	g.Builder.WriteString("return [")
	for _, v := range metadata.Variables {
		splits := strings.Split(v.Name, "/")
		for i, split := range splits {
			splits[i] = dartField(split)
		}
		path := strings.Join(splits, ".")
		g.Builder.WriteString(fmt.Sprintf("('%s', %s),", strings.ReplaceAll(path, "$", "\\$"), path))
	}
	g.Builder.WriteString("];")
	g.Builder.WriteString("}")

	g.Builder.WriteString("}")
}

func (g *FlutterCodeGen) generateBundleDataTypeDefinition(metadata *uipack.BundleMetadata) {
	collections := metadata.BuildTree()
	g.Builder.WriteString("typedef Bundle = ({\n")
	g.Builder.WriteString("\tint identifier,\n")
	g.generateBundleVariableCollectionTypeDefinition(metadata, collections)
	g.Builder.WriteString("});")
}

func (g *FlutterCodeGen) GenerateBundleLoader(metadata *uipack.BundleMetadata) string {
	g.Builder.WriteString("import 'dart:typed_data';\n\n")
	g.Builder.WriteString("import 'package:flutter/widgets.dart';\n\n")
	g.Builder.WriteString("import 'data.g.dart';\n\n")
	g.Builder.WriteString(`class BundleLoader {
  const BundleLoader(this.d);
  final ByteData d;
  Bundle load() {
    var o = 0;
    T read<T>(T Function(int offset) f, int size) {
      final result = f(o);
      o += size;
      return result;
    }

    int uint8() => read(d.getUint8, 1);
    int uint32() => read(d.getUint32, 4);
    int uint64() => read(d.getUint64, 8);
    double float64() => read(d.getFloat64, 8);`)
	g.Builder.WriteString(`String string() {
      final l = uint32();
      final offset = o;
      o += l;
      return String.fromCharCodes(
        d.buffer.asUint8List(offset, l),
      );
    }`)
	g.Builder.WriteString(`Color color() => Color(uint32());`)
	g.Builder.WriteString(`TextStyle textStyle() => TextStyle(
          fontFamily: string(),
          fontSize: float64(),
          fontWeight: switch (uint8()) {
            0 => FontWeight.w100,
            1 => FontWeight.w200,
            2 => FontWeight.w300,
            4 => FontWeight.w500,
            5 => FontWeight.w600,
            6 => FontWeight.w700,
            7 => FontWeight.w800,
            8 => FontWeight.w900,
            _ => FontWeight.w400,
          },
          letterSpacing: float64(),
          wordSpacing: float64(),
          height: float64(),
        );`)
	g.Builder.WriteString(`LinearGradient linearGradient() => LinearGradient(
          begin: Alignment(float64(), float64()),
          end: Alignment(float64(), float64()),
          colors: List.generate(
            uint32(),
            (_) => color(),
          ),
          stops: List.generate(
            uint32(),
            (_) => float64(),
          ),
        );`)
	g.Builder.WriteString(`RadialGradient radialGradient() => RadialGradient(
		  center: Alignment(float64(), float64()),
		  radius: float64(),
		  colors: List.generate(
			uint32(),
			(_) => color(),
		  ),
		  stops: List.generate(
			uint32(),
			(_) => float64(),
		  ),	
		);`)
	g.Builder.WriteString(`Offset offset() => Offset(float64(), float64());`)
	g.Builder.WriteString(`Radius radius() => Radius.circular(float64(), float64());`)
	g.Builder.WriteString(`BorderRadius borderRadius() => BorderRadius.only(
		  topLeft: radius(),
		  topRight: radius(),
		  bottomLeft: radius(),
		  bottomRight: radius(),
		);`)

	getMethod := func(v uipack.ValueType) string {
		switch v.Type {
		case uipack.ColorType:
			return "color()"
		case uipack.TextStyleType:
			return "textStyle()"
		case uipack.LinearGradientType:
			return "linearGradient()"
		case uipack.RadialGradientType:
			return "radialGradient()"
		case uipack.StringType:
			return "string()"
		case uipack.BooleanType:
			return "uint8() == 1"
		case uipack.FloatType:
			return "float64()"
		case uipack.IntegerType:
			return "uint64()"
		case uipack.OffsetType:
			return "offset()"
		case uipack.RadiusType:
			return "radius`()"
		case uipack.BorderRadiusType:
			return "borderRadius()"
		case uipack.LabelType:
			return "label()"
		case uipack.CustomType:
			return "instance()"
		}
		return ""
	}

	g.Builder.WriteString(`Object instance() {
			switch(uint64()) {`)

	for _, typedef := range metadata.Types {
		g.Builder.WriteString(fmt.Sprintf("case 0x%x: return %s(", typedef.Identifier, dartType(typedef.Name)))
		for _, v := range typedef.Properties {
			g.Builder.WriteString(fmt.Sprintf("%s: %s,", dartField(v.Name), getMethod(v.Type)))
		}
		g.Builder.WriteString(");")
	}
	g.Builder.WriteString(`
			_ => throw Exception('Unknown instance type');
		  }
		}`)
	g.Builder.WriteString("final values = <dynamic>[\n")
	g.Builder.WriteString("uint64(), // Identifier\n")

	for _, v := range metadata.Variables {
		g.Builder.WriteString(getMethod(v.Type))
		g.Builder.WriteString(fmt.Sprintf("//%s\n", v.Name))

	}

	g.Builder.WriteString("];\n")

	g.Builder.WriteString("return (\n")
	g.Builder.WriteString("identifier: values[0] as int,\n")

	g.generateBundleLoaderCollectionInstance(metadata.BuildTree())
	g.Builder.WriteString(");")
	g.Builder.WriteString("}")
	g.Builder.WriteString("}")

	return g.Builder.String()
}

func (g *FlutterCodeGen) generateBundleLoaderCollectionInstance(collection uipack.VariableCollection) {
	for _, v := range collection.Variables {
		g.Builder.WriteString(dartField(v.Name))

		g.Builder.WriteString(fmt.Sprintf(": values[%d],\n", 1+v.Variable.Identifier))
	}

	for _, c := range collection.Collections {
		g.Builder.WriteString(dartField(c.Name))
		g.Builder.WriteString(": (\n")
		g.generateBundleLoaderCollectionInstance(c)
		g.Builder.WriteString("),\n")
	}
}

func (g *FlutterCodeGen) generateBundleVariableCollectionTypeDefinition(metadata *uipack.BundleMetadata, collection uipack.VariableCollection) {

	for _, v := range collection.Variables {
		g.generateBundleVariantVariableDefinition(metadata, &v)
		g.Builder.WriteString(",\n")
	}

	for _, c := range collection.Collections {
		g.Builder.WriteString("({\n")
		g.generateBundleVariableCollectionTypeDefinition(metadata, c)
		g.Builder.WriteString("})")
		g.Builder.WriteString(" ")
		g.Builder.WriteString(dartField(c.Name))
		g.Builder.WriteString(",\n")
	}

}

func generateDartType(metadata *uipack.BundleMetadata, t uipack.ValueType) string {
	switch t.Type {
	case uipack.ColorType:
		return "Color"
	case uipack.TextStyleType:
		return "TextStyle"
	case uipack.LinearGradientType:
		return "LinearGradient"
	case uipack.RadialGradientType:
		return "RadialGradient"
	case uipack.StringType:
		return "String"
	case uipack.BooleanType:
		return "bool"
	case uipack.FloatType:
		return "double"
	case uipack.IntegerType:
		return "int"
	case uipack.OffsetType:
		return "Offset"
	case uipack.RadiusType:
		return "Radius"
	case uipack.BorderRadiusType:
		return "BorderRadius"
	case uipack.CustomType:
		typedef := metadata.FindTypeDefintion(t.CustomType)
		return dartType(typedef.Name)
	}
	return ""
}

func (g *FlutterCodeGen) generateBundleVariantVariableDefinition(metadata *uipack.BundleMetadata, v *uipack.VariableCollectionVariable) {
	g.Builder.WriteString(generateDartType(metadata, v.Variable.Type))
	g.Builder.WriteString(" ")
	g.Builder.WriteString(dartField(v.Name))
}

func (g *FlutterCodeGen) GenerateBundle(metadata *uipack.BundleMetadata, bundle *uipack.Bundle) string {
	g.Builder.WriteString("// ignore_for_file: prefer_const_constructors\n\n")
	g.Builder.WriteString("import 'package:flutter/widgets.dart';\n\n")
	g.Builder.WriteString("import 'data.g.dart';\n\n")

	g.generateBundleInstance(metadata, bundle)

	return g.Builder.String()
}

func (g *FlutterCodeGen) generateBundleInstance(metadata *uipack.BundleMetadata, bundle *uipack.Bundle) {

	g.Builder.WriteString("// Variant :")
	for _, mode := range metadata.Modes {
		g.Builder.WriteString(" ")
		g.Builder.WriteString(mode.Name)
		g.Builder.WriteString(":")
		value := mode.Variants[bundle.Variant.GetMode(mode.Identifier)]
		g.Builder.WriteString(value.Name)
	}
	g.Builder.WriteString("\n")

	identifier := fmt.Sprintf("%x", bundle.Variant)
	collections := metadata.BuildTree()
	g.Builder.WriteString(fmt.Sprintf("Bundle bundle%s() => (\n", identifier))
	g.Builder.WriteString(fmt.Sprintf("identifier: 0x%s,", identifier))
	g.generateBundleVariableCollectionInstance(collections, bundle)
	g.Builder.WriteString(");")
}

func (g *FlutterCodeGen) generateBundleVariableCollectionInstance(collection uipack.VariableCollection, bundle *uipack.Bundle) {
	for _, v := range collection.Variables {

		g.Builder.WriteString(dartField(v.Name))
		g.Builder.WriteString(": \n")
		value := bundle.Values[v.Variable.Identifier]
		g.generateBundleVariableInstance(value)
		g.Builder.WriteString(",\n")
	}

	for _, c := range collection.Collections {
		g.Builder.WriteString(dartField(c.Name))
		g.Builder.WriteString(": (\n")
		g.generateBundleVariableCollectionInstance(c, bundle)
		g.Builder.WriteString("),\n")
	}
}

func (g *FlutterCodeGen) generateBundleVariableInstance(v interface{}) {

	generateGradientStops := func(stops []uipack.GradientStop) {
		g.Builder.WriteString("colors: [")
		for _, stop := range stops {
			g.Builder.WriteString(fmt.Sprintf("Color(0x%s),", stop.Color.ToHexString()))
		}
		g.Builder.WriteString("],")
		g.Builder.WriteString("stops: [")
		for _, stop := range stops {
			g.Builder.WriteString(fmt.Sprintf("%.2f,", stop.Offset))
		}
		g.Builder.WriteString("],")
	}

	switch v := v.(type) {
	case uipack.Color:
		g.Builder.WriteString(fmt.Sprintf("Color(0x%s)", v.ToHexString()))
	case uipack.TextStyle:
		fontFamily := v.FontFamily
		if fontFamily == "SF Pro Display" {
			fontFamily = ".SF UI Display"
		}
		if fontFamily == "SF Pro" {
			fontFamily = ".SF UI Text"
		}

		g.Builder.WriteString("TextStyle(")
		g.Builder.WriteString(fmt.Sprintf("fontFamily: %s,", dartStringLiteral(fontFamily)))
		g.Builder.WriteString(fmt.Sprintf("fontSize: %.2f,", v.FontSize))
		g.Builder.WriteString(fmt.Sprintf("letterSpacing: %.2f,", v.LetterSpacing))
		g.Builder.WriteString(fmt.Sprintf("fontWeight: %s,", generateFlutterFontWeight(v.FontWeight)))
		g.Builder.WriteString(fmt.Sprintf("wordSpacing: %.2f,", v.WordSpacing))
		g.Builder.WriteString(fmt.Sprintf("height: %.2f,", generateFlutterLineHeight(v)))
		g.Builder.WriteString("fontVariations: const [")
		for _, variation := range v.FontVariations {
			g.Builder.WriteString(fmt.Sprintf("FontVariation('%s', %.2f),", variation.Axis, variation.Value))
		}
		g.Builder.WriteString("],")

		g.Builder.WriteString(")")
	case uipack.LinearGradient:
		g.Builder.WriteString("LinearGradient(")
		g.Builder.WriteString("begin: Alignment(")
		g.Builder.WriteString(fmt.Sprintf("%.2f, %.2f", v.Begin.X, v.Begin.Y))
		g.Builder.WriteString("),")
		g.Builder.WriteString("end: Alignment(")
		g.Builder.WriteString(fmt.Sprintf("%.2f, %.2f", v.End.X, v.End.Y))
		g.Builder.WriteString("),")
		generateGradientStops(v.Stops)
		g.Builder.WriteString(")")
	case uipack.RadialGradient:
		g.Builder.WriteString("RadialGradient(")
		g.Builder.WriteString("center: Alignment(")
		g.Builder.WriteString(fmt.Sprintf("%.2f, %.2f", v.Center.X, v.Center.Y))
		g.Builder.WriteString("),")
		g.Builder.WriteString("radius: ")
		g.Builder.WriteString(fmt.Sprintf("%.2f", v.Radius))
		generateGradientStops(v.Stops)
		g.Builder.WriteString(")")
	case uipack.Offset:
		if v.X == 0 && v.Y == 0 {
			g.Builder.WriteString("Offset.zero")
		} else {
			g.Builder.WriteString(fmt.Sprintf("Offset(%f, %f)", v.X, v.Y))
		}
	case uipack.Radius:
		if v.X == v.Y {
			if v.X == 0 {
				g.Builder.WriteString("Radius.zero")
			} else {
				g.Builder.WriteString(fmt.Sprintf("Radius.circular(%f)", v.X))
			}
		} else {
			g.Builder.WriteString(fmt.Sprintf("Radius.elliptical(%f, %f)", v.X, v.Y))
		}
	case uipack.BorderRadius:
		g.Builder.WriteString("BorderRadius.only(")

		g.Builder.WriteString("topLeft:")
		g.generateBundleVariableInstance(v.TopLeft)
		g.Builder.WriteString(",")

		g.Builder.WriteString("topRight:")
		g.generateBundleVariableInstance(v.TopRight)
		g.Builder.WriteString(",")

		g.Builder.WriteString("bottomLeft:")
		g.generateBundleVariableInstance(v.BottomLeft)
		g.Builder.WriteString(",")

		g.Builder.WriteString("bottomRight:")
		g.generateBundleVariableInstance(v.BottomRight)
		g.Builder.WriteString(")")

	case string:
		g.Builder.WriteString(dartStringLiteral(v))
	case int64:
		g.Builder.WriteString(fmt.Sprintf("%d", v))
	case float64:
		g.Builder.WriteString(fmt.Sprintf("%f", v))
	case bool:
		switch v {
		case true:
			g.Builder.WriteString("true")
		case false:
			g.Builder.WriteString("false")
		}
	default:
		panic(fmt.Sprint("Unknown variable type ", v))
	}
}

func generateFlutterLineHeight(t uipack.TextStyle) float64 {
	return t.LineHeight / t.FontSize
}

func generateFlutterFontWeight(v uint8) string {
	switch v {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return fmt.Sprintf("FontWeight.w%d00", v)
	default:
		return "FontWeight.w400"
	}
}

func dartField(name string) string {
	return escapeDartKeywords(strcase.ToLowerCamel(cleanName(name)))
}

func dartType(name string) string {
	return escapeDartKeywords(strcase.ToCamel(cleanName(name)))
}

func dartStringLiteral(name string) string {
	return fmt.Sprintf("'%s'", name)
}

func cleanName(name string) string {
	cleanRegexp := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return cleanRegexp.ReplaceAllString(name, "")
}

func escapeDartKeywords(name string) string {
	switch name {
	case "default", "class", "enum", "switch", "while":
		return name + "$"
	}

	match, _ := regexp.MatchString("^[0-9]", name)
	if match {
		return "v" + name
	}

	return name
}
