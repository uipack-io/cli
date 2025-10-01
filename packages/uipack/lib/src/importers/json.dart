import 'dart:convert';
import 'dart:typed_data';
import 'package:uipack/src/importers/importer.dart';

import '../models/package.dart';
import '../models/bundle.dart';
import '../models/metadata.dart';
import '../models/version.dart';
import '../models/variant.dart';
import '../models/color.dart';
import '../models/text_style.dart';

class JsonImporter extends Importer {
  @override
  Package load(Uint8List path) {
    final jsonString = utf8.decode(path);
    final json = jsonDecode(jsonString) as Map<String, dynamic>;
    if (!json.containsKey('collections') ||
        json['collections'] is! List<dynamic>) {
      throw Exception('Invalid JSON format: missing collections');
    }
    final inputCollections = (json['collections'] as List<dynamic>)
        .map((e) => JsonCollection.fromJson(e as Map<String, dynamic>))
        .toList();

    final bundles = <Bundle>[];
    final metadata = _toBundleMetadata(inputCollections);
    final combinations = metadata.generateModeCombinations();

    for (final combination in combinations) {
      Variant vid = Variant(0);
      Map<String, String> variant = {};

      for (int mi = 0; mi < combination.length; mi++) {
        final mv = combination[mi];
        final mode = metadata.modes[mi];
        variant[mode.name] = mv.name;
        vid = vid.setMode(mode.identifier, mv.identifier);
      }

      bundles.add(_toBundle(inputCollections, vid, variant));
    }
    return Package(metadata: metadata, bundles: bundles);
  }

  BundleMetadata _toBundleMetadata(List<JsonCollection> collections) {
    final result = BundleMetadata(
      version: Version(major: 1, minor: 0),
      name: "Figma",
      modes: [],
      variables: [],
      types: [],
    );

    int vi = 0;
    for (int ci = 0; ci < collections.length; ci++) {
      final collection = collections[ci];
      final mm = ModeMetadata(
        identifier: ci,
        name: collection.name,
        variants: [],
      );

      for (int i = 0; i < collection.modes.length; i++) {
        final mode = collection.modes[i];
        mm.variants.add(ModeVariantMetadata(
          identifier: i,
          name: mode.name,
        ));

        if (i == 0) {
          for (final variable in mode.variables) {
            final vm = VariableMetadata(
              identifier: vi,
              name: variable.name,
              type: figmaToVariableType(variable.type),
            );
            result.variables.add(vm);
            vi++;
          }
        }
      }

      result.modes.add(mm);
    }
    return result;
  }

  Bundle _toBundle(List<JsonCollection> collections, Variant identifier,
      Map<String, String> variant) {
    final result = Bundle(
      version: Version(major: 1, minor: 0),
      variant: identifier,
      values: [],
    );

    for (final collection in collections) {
      final String? mn = variant[collection.name];
      if (mn != null) {
        final mode = collection.findMode(mn);
        if (mode != null) {
          for (final variable in mode.variables) {
            result.values.add(_resolveVariable(
                collections, collection.name, variant, variable));
          }
        }
      }
    }
    return result;
  }

  JsonCollection? _findCollection(
      List<JsonCollection> collections, String name) {
    for (final collection in collections) {
      if (collection.name == name) {
        return collection;
      }
    }
    return null;
  }

  dynamic _resolveVariable(
      List<JsonCollection> collections,
      String currentCollection,
      Map<String, String> variant,
      JsonVariable variable) {
    if (variable.isAlias) {
      final alias = variable.alias(currentCollection);

      final acol = _findCollection(collections, alias.collection);
      if (acol != null) {
        final String? mn = variant[alias.collection];
        if (mn != null) {
          final mode = acol.findMode(mn);
          if (mode != null) {
            final aliasVariable = mode.findVariable(alias.name);
            if (aliasVariable != null) {
              return _resolveVariable(
                  collections, currentCollection, variant, aliasVariable);
            }
            throw Exception('Alias variable not found ${alias.name}');
          }
          throw Exception('Alias mode not found $mn');
        }
        throw Exception(
            'Alias mode not found for collection ${alias.collection}');
      }
      throw Exception('Alias collection not found ${alias.collection}');
    }

    return figmaToVariableValue(variable);
  }
}

dynamic figmaToVariableValue(JsonVariable variable) {
  if (variable.isAlias) {
    throw Exception('Alias not resolved');
  }

  switch (variable.type) {
    case 'color':
      if (variable.value is String) {
        final result = Color(
          red: 0.0,
          green: 0.0,
          blue: 0.0,
          alpha: 1.0,
          colorSpace: ColorSpace.srgb,
        );
        result.parseHexString(variable.value as String);
        return result;
      } else {
        throw Exception('Unknown color type');
      }

    case 'typography':
      if (variable.value is Map<String, dynamic>) {
        final typographyData = variable.value as Map<String, dynamic>;
        final fontSize = (typographyData['fontSize'] as num).toDouble();
        return TextStyle(
          fontFamily: typographyData['fontFamily'] as String,
          fontSize: fontSize,
          fontWeight:
              figmaFontWeightToIndex(typographyData['fontWeight'] as String),
          letterSpacing: figmaToLetterSpacing(
            (typographyData['letterSpacing'] as num).toDouble(),
            fontSize,
            typographyData['letterSpacingUnit'],
          ),
          wordSpacing: 0.0,
          lineHeight: figmaToLineHeight(
            (typographyData['lineHeight'] as num).toDouble(),
            fontSize,
            typographyData['lineHeightUnit'],
          ),
          fontVariations: [],
        );
      } else {
        throw Exception('Typography should be a map');
      }

    case 'number':
      if (variable.value is int) {
        return (variable.value as int).toDouble();
      } else if (variable.value is double) {
        return variable.value as double;
      } else {
        throw Exception('Number should be an int or a double');
      }

    case 'string':
      if (variable.value is String) {
        return variable.value as String;
      } else {
        throw Exception('String should be a string');
      }

    default:
      throw Exception('Unknown type ${variable.type}');
  }
}

double figmaToLineHeight(double value, double fontSize, dynamic unit) {
  if (unit is String) {
    switch (unit) {
      case 'PERCENT':
        return value * 0.01 * fontSize;
      case 'PIXELS':
        return value;
      default:
        return value;
    }
  }
  return value;
}

double figmaToLetterSpacing(double value, double fontSize, dynamic unit) {
  if (unit is String) {
    switch (unit) {
      case 'PERCENT':
        return value * 0.01 * fontSize;
      case 'PIXELS':
        return value;
      default:
        return value;
    }
  }
  return value;
}

ValueType figmaToVariableType(String type) {
  switch (type) {
    case 'color':
      return ValueType(type: MainValueType.color);
    case 'typography':
      return ValueType(type: MainValueType.textStyle);
    case 'number':
      return ValueType(type: MainValueType.float);
    case 'string':
      return ValueType(type: MainValueType.string);
    default:
      return ValueType(type: MainValueType.textStyle);
  }
}

int figmaFontWeightToIndex(String value) {
  switch (value.toLowerCase()) {
    case 'thin':
      return 1;
    case 'extralight':
      return 2;
    case 'light':
      return 3;
    case 'regular':
      return 4;
    case 'medium':
      return 5;
    case 'semibold':
      return 6;
    case 'bold':
      return 7;
    case 'extrabold':
      return 8;
    case 'black':
      return 9;
    default:
      return 4;
  }
}

class JsonCollection {
  String name;
  List<JsonMode> modes;

  JsonCollection({required this.name, required this.modes});

  factory JsonCollection.fromJson(Map<String, dynamic> json) {
    return JsonCollection(
      name: json['name'] as String,
      modes: (json['modes'] as List<dynamic>)
          .map((e) => JsonMode.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  JsonMode? findMode(String name) {
    for (final mode in modes) {
      if (mode.name == name) {
        return mode;
      }
    }
    return null;
  }
}

class JsonMode {
  String name;
  List<JsonVariable> variables;

  JsonMode({required this.name, required this.variables});

  factory JsonMode.fromJson(Map<String, dynamic> json) {
    return JsonMode(
      name: json['name'] as String,
      variables: (json['variables'] as List<dynamic>)
          .map((e) => JsonVariable.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  JsonVariable? findVariable(String name) {
    for (final variable in variables) {
      if (variable.name == name) {
        return variable;
      }
    }
    return null;
  }
}

class JsonVariable {
  String name;
  String type;
  bool isAlias;
  dynamic value;

  JsonVariable({
    required this.name,
    required this.type,
    required this.isAlias,
    required this.value,
  });

  factory JsonVariable.fromJson(Map<String, dynamic> json) {
    return JsonVariable(
      name: json['name'] as String,
      type: json['type'] as String,
      isAlias: json['isAlias'] as bool,
      value: json['value'],
    );
  }

  JsonAlias alias(String currentCollection) {
    if (isAlias && value is Map<String, dynamic>) {
      final aliasData = value as Map<String, dynamic>;
      final name = aliasData['name'] as String;
      final collection =
          aliasData['collection'] as String? ?? currentCollection;
      return JsonAlias(collection: collection, name: name);
    }
    return JsonAlias(collection: '', name: '');
  }
}

class JsonAlias {
  String collection;
  String name;

  JsonAlias({required this.collection, required this.name});
}

