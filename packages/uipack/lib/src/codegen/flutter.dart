import '../models/bundle.dart';
import '../models/metadata.dart';
import '../models/color.dart';
import '../models/text_style.dart';

class FlutterCodeGen {
  final StringBuffer _builder = StringBuffer();

  String generateDefinitions(BundleMetadata metadata, List<Bundle> bundles) {
    _builder.clear();
    _builder.writeln("import 'package:flutter/widgets.dart';\n");

    for (final bundle in bundles) {
      final identifier = bundle.variant.value.toRadixString(16);
      _builder.writeln("import 'bundle_$identifier.g.dart';");
    }

    _builder.writeln();

    _generateModeTypeDefinitions(metadata);
    _generateVariantTypeDefinition(metadata, bundles);
    _generateBundleDataTypeDefinition(metadata);
    _generateBundleExtension(metadata);
    return _builder.toString();
  }

  void _generateModeTypeDefinitions(BundleMetadata metadata) {
    for (final mode in metadata.modes) {
      _generateModeEnumDefinition(mode);
    }
  }

  void _generateModeEnumDefinition(ModeMetadata mode) {
    _builder.writeln('/// Index : ${mode.identifier.toRadixString(16)}');
    _builder.writeln('enum ${_dartType(mode.name)}Mode {');
    for (final v in mode.variants) {
      _builder.writeln('${_dartField(v.name)},');
    }
    _builder.writeln('}');
  }

  void _generateCustomTypeDefinitions(BundleMetadata metadata) {
    for (final typedef in metadata.types) {
      _builder.writeln('class ${_dartType(typedef.name)} {');
      _builder.writeln('const ${_dartType(typedef.name)}();');
      _builder.writeln('}\n');
    }
  }

  void _generateVariantTypeDefinition(BundleMetadata metadata, List<Bundle> bundles) {
    _builder.write('typedef Variant = ({');
    for (final mode in metadata.modes) {
      _builder.write('${_dartType(mode.name)}Mode ${_dartField(mode.name)},');
    }
    _builder.writeln('});');

    _builder.writeln('extension VariantExtension on Variant {');

    // Identifier
    _builder.writeln('int get identifier {');
    _builder.writeln('var result = 0;');
    for (int i = 0; i < metadata.modes.length; i++) {
      final mode = metadata.modes[i];
      _builder.writeln('result |= ${_dartField(mode.name)}.index << ${i * 4};');
    }
    _builder.writeln('return result;');
    _builder.writeln('}');

    // Bundle
    _builder.writeln('Bundle? get bundle {');
    _builder.writeln('switch (identifier) {');
    for (final bundle in bundles) {
      _builder.writeln('case 0x${bundle.variant.value.toRadixString(16)}: return bundle${bundle.variant.value.toRadixString(16)}();');
    }
    _builder.writeln('default: return null;');
    _builder.writeln('}');
    _builder.writeln('}');
    _builder.writeln('}');

    // Custom types
    _generateCustomTypeDefinitions(metadata);
  }

  void _generateBundleExtension(BundleMetadata metadata) {
    _builder.writeln('extension BundleExtension on Bundle {');

    // Metadata
    _builder.writeln('List<(String, dynamic)> get metadata {');
    _builder.writeln('return [');
    for (final v in metadata.variables) {
      final splits = v.name.split('/');
      final path = splits.map(_dartField).join('.');
      _builder.writeln("('${path.replaceAll(r'$', r'\$')}', $path),");
    }
    _builder.writeln('];');
    _builder.writeln('}');
    _builder.writeln('}');
  }

  void _generateBundleDataTypeDefinition(BundleMetadata metadata) {
    final collections = _buildTree(metadata);
    _builder.writeln('typedef Bundle = ({');
    _builder.writeln('\tint identifier,');
    _generateBundleVariableCollectionTypeDefinition(metadata, collections);
    _builder.writeln('});');
  }

  String generateBundleLoader(BundleMetadata metadata) {
    _builder.clear();
    _builder.writeln("import 'dart:typed_data';\n");
    _builder.writeln("import 'package:flutter/widgets.dart';\n");
    _builder.writeln("import 'data.g.dart';\n");
    
    _builder.writeln('''class BundleLoader {
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
    double float64() => read(d.getFloat64, 8);''');
    
    _builder.writeln('''String string() {
      final l = uint32();
      final offset = o;
      o += l;
      return String.fromCharCodes(
        d.buffer.asUint8List(offset, l),
      );
    }''');
    
    _builder.writeln('''Color color() => Color.from(
      red: float64(),
      green: float64(),
      blue: float64(),
      alpha: float64(),
      colorSpace: switch (uint8()) {
        0 => ColorSpace.sRGB,
        1 => ColorSpace.extendedSRGB,
        2 => ColorSpace.displayP3,
        _ => ColorSpace.sRGB,
      },
    );''');
    
    _builder.writeln('''TextStyle textStyle() => TextStyle(
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
        );''');
    
    _builder.writeln('''LinearGradient linearGradient() => LinearGradient(
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
        );''');
    
    _builder.writeln('''RadialGradient radialGradient() => RadialGradient(
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
		);''');
    
    _builder.writeln('Offset offset() => Offset(float64(), float64());');
    _builder.writeln('Radius radius() => Radius.circular(float64());');
    _builder.writeln('''BorderRadius borderRadius() => BorderRadius.only(
		  topLeft: radius(),
		  topRight: radius(),
		  bottomLeft: radius(),
		  bottomRight: radius(),
		);''');

    _builder.writeln('Object instance() => switch(uint64()) {');
    
    for (final typedef in metadata.types) {
      _builder.writeln('case 0x${typedef.identifier.toRadixString(16)}: return ${_dartType(typedef.name)}(),');
    }
    _builder.writeln('_ => throw Exception(\'Unknown instance type\'),');
    _builder.writeln('};');
    
    _builder.writeln('final values = <dynamic>[');
    _builder.writeln('uint64(), // Identifier');
    
    for (final v in metadata.variables) {
      _builder.writeln('${_getMethod(v.type)}, //${v.name}');
    }
    
    _builder.writeln('];');
    
    _builder.writeln('return (');
    _builder.writeln('identifier: values[0] as int,');
    
    _generateBundleLoaderCollectionInstance(_buildTree(metadata));
    _builder.writeln(');');
    _builder.writeln('}');
    _builder.writeln('}');

    return _builder.toString();
  }

  void _generateBundleLoaderCollectionInstance(VariableCollection collection) {
    for (final v in collection.variables) {
      _builder.writeln('${_dartField(v.name)}: values[${1 + v.variable.identifier}],');
    }

    for (final c in collection.collections) {
      _builder.writeln('${_dartField(c.name)}: (');
      _generateBundleLoaderCollectionInstance(c);
      _builder.writeln('),');
    }
  }

  void _generateBundleVariableCollectionTypeDefinition(BundleMetadata metadata, VariableCollection collection) {
    for (final v in collection.variables) {
      _generateBundleVariantVariableDefinition(metadata, v);
      _builder.writeln(',');
    }

    for (final c in collection.collections) {
      _builder.writeln('({');
      _generateBundleVariableCollectionTypeDefinition(metadata, c);
      _builder.writeln('}) ${_dartField(c.name)},');
    }
  }

  String _generateDartType(BundleMetadata metadata, ValueType t) {
    switch (t.type) {
      case MainValueType.color:
        return 'Color';
      case MainValueType.textStyle:
        return 'TextStyle';
      case MainValueType.linearGradient:
        return 'LinearGradient';
      case MainValueType.radialGradient:
        return 'RadialGradient';
      case MainValueType.string:
        return 'String';
      case MainValueType.boolean:
        return 'bool';
      case MainValueType.float:
        return 'double';
      case MainValueType.integer:
        return 'int';
      case MainValueType.offset:
        return 'Offset';
      case MainValueType.radius:
        return 'Radius';
      case MainValueType.borderRadius:
        return 'BorderRadius';
      case MainValueType.custom:
        final typedef = metadata.findTypeDefinition(t.customType);
        return _dartType(typedef?.name ?? 'Object');
      default:
        return 'dynamic';
    }
  }

  String _getMethod(ValueType v) {
    switch (v.type) {
      case MainValueType.color:
        return 'color()';
      case MainValueType.textStyle:
        return 'textStyle()';
      case MainValueType.linearGradient:
        return 'linearGradient()';
      case MainValueType.radialGradient:
        return 'radialGradient()';
      case MainValueType.string:
        return 'string()';
      case MainValueType.boolean:
        return 'uint8() == 1';
      case MainValueType.float:
        return 'float64()';
      case MainValueType.integer:
        return 'uint64()';
      case MainValueType.offset:
        return 'offset()';
      case MainValueType.radius:
        return 'radius()';
      case MainValueType.borderRadius:
        return 'borderRadius()';
      case MainValueType.label:
        return 'label()';
      case MainValueType.custom:
        return 'instance()';
      default:
        return '';
    }
  }

  void _generateBundleVariantVariableDefinition(BundleMetadata metadata, VariableCollectionVariable v) {
    _builder.write('${_generateDartType(metadata, v.variable.type)} ${_dartField(v.name)}');
  }

  String generateBundle(BundleMetadata metadata, Bundle bundle) {
    _builder.clear();
    _builder.writeln('// ignore_for_file: prefer_const_constructors\n');
    _builder.writeln("import 'dart:ui';\n");
    _builder.writeln("import 'package:flutter/widgets.dart';\n");
    _builder.writeln("import 'data.g.dart';\n");

    _generateBundleInstance(metadata, bundle);

    return _builder.toString();
  }

  void _generateBundleInstance(BundleMetadata metadata, Bundle bundle) {
    _builder.write('// Variant :');
    for (final mode in metadata.modes) {
      _builder.write(' ${mode.name}:');
      final value = mode.variants[bundle.variant.getMode(mode.identifier)];
      _builder.write(value.name);
    }
    _builder.writeln();

    final identifier = bundle.variant.value.toRadixString(16);
    final collections = _buildTree(metadata);
    _builder.writeln('Bundle bundle$identifier() => (');
    _builder.writeln('identifier: 0x$identifier,');
    _generateBundleVariableCollectionInstance(collections, bundle);
    _builder.writeln(');');
  }

  void _generateBundleVariableCollectionInstance(VariableCollection collection, Bundle bundle) {
    for (final v in collection.variables) {
      _builder.writeln('${_dartField(v.name)}: ');
      final value = bundle.values[v.variable.identifier];
      _generateBundleVariableInstance(value);
      _builder.writeln(',');
    }

    for (final c in collection.collections) {
      _builder.writeln('${_dartField(c.name)}: (');
      _generateBundleVariableCollectionInstance(c, bundle);
      _builder.writeln('),');
    }
  }

  void _generateBundleVariableInstance(dynamic v) {
    void generateColor(Color color) {
      _builder.write('Color.from(');
      _builder.write('red: ${color.red.toStringAsFixed(4)},');
      _builder.write('green: ${color.green.toStringAsFixed(4)},');
      _builder.write('blue: ${color.blue.toStringAsFixed(4)},');
      _builder.write('alpha: ${color.alpha.toStringAsFixed(4)},');
      _builder.write('colorSpace: ');
      switch (color.colorSpace) {
        case ColorSpace.srgb:
          _builder.write('ColorSpace.sRGB');
          break;
        case ColorSpace.extendedSrgb:
          _builder.write('ColorSpace.extendedSRGB');
          break;
        case ColorSpace.displayP3:
          _builder.write('ColorSpace.displayP3');
          break;
      }
      _builder.write(',)');
    }

    void generateGradientStops(List<GradientStop> stops) {
      _builder.write('colors: [');
      for (final stop in stops) {
        generateColor(stop.color);
      }
      _builder.write('],');
      _builder.write('stops: [');
      for (final stop in stops) {
        _builder.write('${stop.offset.toStringAsFixed(2)},');
      }
      _builder.write('],');
    }

    switch (v.runtimeType) {
      case Color:
        generateColor(v as Color);
        break;
      case TextStyle:
        final textStyle = v as TextStyle;
        String fontFamily = textStyle.fontFamily;
        if (fontFamily == 'SF Pro Display') {
          fontFamily = '.SF UI Display';
        }
        if (fontFamily == 'SF Pro') {
          fontFamily = '.SF UI Text';
        }

        _builder.write('TextStyle(');
        _builder.write('fontFamily: ${_dartStringLiteral(fontFamily)},');
        _builder.write('fontSize: ${textStyle.fontSize.toStringAsFixed(2)},');
        _builder.write('letterSpacing: ${textStyle.letterSpacing.toStringAsFixed(2)},');
        _builder.write('fontWeight: ${_generateFlutterFontWeight(textStyle.fontWeight)},');
        _builder.write('wordSpacing: ${textStyle.wordSpacing.toStringAsFixed(2)},');
        _builder.write('height: ${_generateFlutterLineHeight(textStyle).toStringAsFixed(2)},');
        _builder.write('fontVariations: const [');
        for (final variation in textStyle.fontVariations) {
          _builder.write("FontVariation('${variation.axis}', ${variation.value.toStringAsFixed(2)}),");
        }
        _builder.write('],');
        _builder.write(')');
        break;
      case LinearGradient:
        final gradient = v as LinearGradient;
        _builder.write('LinearGradient(');
        _builder.write('begin: Alignment(${gradient.begin.x.toStringAsFixed(2)}, ${gradient.begin.y.toStringAsFixed(2)}),');
        _builder.write('end: Alignment(${gradient.end.x.toStringAsFixed(2)}, ${gradient.end.y.toStringAsFixed(2)}),');
        generateGradientStops(gradient.stops);
        _builder.write(')');
        break;
      case RadialGradient:
        final gradient = v as RadialGradient;
        _builder.write('RadialGradient(');
        _builder.write('center: Alignment(${gradient.center.x.toStringAsFixed(2)}, ${gradient.center.y.toStringAsFixed(2)}),');
        _builder.write('radius: ${gradient.radius.toStringAsFixed(2)},');
        generateGradientStops(gradient.stops);
        _builder.write(')');
        break;
      case Offset:
        final offset = v as Offset;
        if (offset.x == 0 && offset.y == 0) {
          _builder.write('Offset.zero');
        } else {
          _builder.write('Offset(${offset.x}, ${offset.y})');
        }
        break;
      case Radius:
        final radius = v as Radius;
        if (radius.x == radius.y) {
          if (radius.x == 0) {
            _builder.write('Radius.zero');
          } else {
            _builder.write('Radius.circular(${radius.x})');
          }
        } else {
          _builder.write('Radius.elliptical(${radius.x}, ${radius.y})');
        }
        break;
      case BorderRadius:
        final borderRadius = v as BorderRadius;
        _builder.write('BorderRadius.only(');
        _builder.write('topLeft:');
        _generateBundleVariableInstance(borderRadius.topLeft);
        _builder.write(',');
        _builder.write('topRight:');
        _generateBundleVariableInstance(borderRadius.topRight);
        _builder.write(',');
        _builder.write('bottomLeft:');
        _generateBundleVariableInstance(borderRadius.bottomLeft);
        _builder.write(',');
        _builder.write('bottomRight:');
        _generateBundleVariableInstance(borderRadius.bottomRight);
        _builder.write(')');
        break;
      case String:
        _builder.write(_dartStringLiteral(v as String));
        break;
      case int:
        _builder.write('${v as int}');
        break;
      case double:
        _builder.write('${v as double}');
        break;
      case bool:
        _builder.write(v as bool ? 'true' : 'false');
        break;
      default:
        throw Exception('Unknown variable type $v');
    }
  }

  double _generateFlutterLineHeight(TextStyle t) {
    return t.lineHeight / t.fontSize;
  }

  String _generateFlutterFontWeight(int v) {
    switch (v) {
      case 1:
      case 2:
      case 3:
      case 4:
      case 5:
      case 6:
      case 7:
      case 8:
      case 9:
        return 'FontWeight.w${v}00';
      default:
        return 'FontWeight.w400';
    }
  }

  String _dartField(String name) {
    return _escapeDartKeywords(_lowerCamelCase(_cleanName(name)));
  }

  String _dartType(String name) {
    return _escapeDartKeywords(_upperCamelCase(_cleanName(name)));
  }

  String _dartStringLiteral(String name) {
    return "'$name'";
  }

  String _cleanName(String name) {
    return name.replaceAll(RegExp(r'[^a-zA-Z0-9]'), '');
  }

  String _escapeDartKeywords(String name) {
    switch (name) {
      case 'default':
      case 'class':
      case 'enum':
      case 'switch':
      case 'while':
        return '$name\$';
    }

    if (RegExp(r'^[0-9]').hasMatch(name)) {
      return 'v$name';
    }

    return name;
  }

  String _lowerCamelCase(String input) {
    if (input.isEmpty) return input;
    final words = input.split(RegExp(r'[^a-zA-Z0-9]')).where((w) => w.isNotEmpty).toList();
    if (words.isEmpty) return input;
    
    final result = StringBuffer(words.first.toLowerCase());
    for (int i = 1; i < words.length; i++) {
      final word = words[i];
      if (word.isNotEmpty) {
        result.write(word[0].toUpperCase() + word.substring(1).toLowerCase());
      }
    }
    return result.toString();
  }

  String _upperCamelCase(String input) {
    if (input.isEmpty) return input;
    final words = input.split(RegExp(r'[^a-zA-Z0-9]')).where((w) => w.isNotEmpty).toList();
    if (words.isEmpty) return input;
    
    final result = StringBuffer();
    for (final word in words) {
      if (word.isNotEmpty) {
        result.write(word[0].toUpperCase() + word.substring(1).toLowerCase());
      }
    }
    return result.toString();
  }

  // Helper method to build tree structure - this would need to be implemented
  // based on the actual metadata structure in your Dart models
  VariableCollection _buildTree(BundleMetadata metadata) {
    // This is a placeholder implementation
    // You would need to implement the actual tree building logic
    // based on your metadata structure
    return VariableCollection(
      name: 'root',
      variables: metadata.variables.map((v) => 
        VariableCollectionVariable(name: v.name, variable: v)
      ).toList(),
      collections: [],
    );
  }
}

// Helper classes that may need to be defined based on your models
class VariableCollection {
  final String name;
  final List<VariableCollectionVariable> variables;
  final List<VariableCollection> collections;

  VariableCollection({
    required this.name,
    required this.variables,
    required this.collections,
  });
}

class VariableCollectionVariable {
  final String name;
  final VariableMetadata variable;

  VariableCollectionVariable({
    required this.name,
    required this.variable,
  });
}

// These classes would need to be defined or imported from your models
class LinearGradient {
  final Offset begin;
  final Offset end;
  final List<GradientStop> stops;

  LinearGradient({
    required this.begin,
    required this.end,
    required this.stops,
  });
}

class RadialGradient {
  final Offset center;
  final double radius;
  final List<GradientStop> stops;

  RadialGradient({
    required this.center,
    required this.radius,
    required this.stops,
  });
}

class GradientStop {
  final Color color;
  final double offset;

  GradientStop({
    required this.color,
    required this.offset,
  });
}

class Offset {
  final double x;
  final double y;

  Offset(this.x, this.y);
}

class Radius {
  final double x;
  final double y;

  Radius(this.x, this.y);
}

class BorderRadius {
  final Radius topLeft;
  final Radius topRight;
  final Radius bottomLeft;
  final Radius bottomRight;

  BorderRadius({
    required this.topLeft,
    required this.topRight,
    required this.bottomLeft,
    required this.bottomRight,
  });
}

