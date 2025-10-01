import 'dart:io';
import 'dart:typed_data';
import 'package:path/path.dart' as path;

import 'models/package.dart';
import 'models/bundle.dart';
import 'models/metadata.dart';
import 'models/version.dart';
import 'models/variant.dart';
import 'encoding/binary_encoding.dart';
import 'helpers/style_helper.dart';

class PackageLoader {
  static Package loadPackage(String inputPath) {
    final absolutePath = path.absolute(inputPath);
    final file = File(absolutePath);
    final directory = Directory(absolutePath);

    print(StyleHelper.listStarted('Loading package from $absolutePath'));

    Package package;

    if (directory.existsSync()) {
      package = _loadFromDirectory(absolutePath);
    } else if (file.existsSync()) {
      if (path.extension(absolutePath) == '.json') {
        final jsonData = file.readAsBytesSync();
        package = _loadFromJson(jsonData);
      } else {
        throw Exception('ZIP archive support not implemented');
      }
    } else {
      throw Exception('Path not found: $absolutePath');
    }

    _printPackageInfo(package);
    return package;
  }

  static Package _loadFromDirectory(String dirPath) {
    final metadataFile = File(path.join(dirPath, '_'));
    if (!metadataFile.existsSync()) {
      throw Exception('Metadata file not found: ${metadataFile.path}');
    }

    final metadataBytes = metadataFile.readAsBytesSync();
    final metadataReader = ByteDataReader(metadataBytes);
    final metadata = BundleMetadata(
      version: Version(major: 0, minor: 0),
      name: '',
      modes: [],
      variables: [],
      types: [],
    );
    metadata.decode(metadataReader);

    final bundles = <Bundle>[];
    final directory = Directory(dirPath);

    for (final entity in directory.listSync()) {
      if (entity is File && entity.path != metadataFile.path) {
        final fileName = path.basename(entity.path);
        if (path.extension(fileName).isEmpty && fileName != '_') {
          final bundleBytes = entity.readAsBytesSync();
          final bundleReader = ByteDataReader(bundleBytes);
          final bundle = Bundle(
            version: Version(major: 0, minor: 0),
            variant: Variant(0),
            values: [],
          );
          bundle.decode(bundleReader, metadata);
          bundles.add(bundle);
        }
      }
    }

    return Package(metadata: metadata, bundles: bundles);
  }

  static Package _loadFromJson(Uint8List jsonData) {
    throw UnimplementedError('JSON import not yet implemented');
  }

  static void _printPackageInfo(Package package) {
    print(StyleHelper.listSubinfo(package.metadata.name));
    print(
      StyleHelper.listSubinfo(
        'Version ${package.metadata.version.major}.${package.metadata.version.minor}',
      ),
    );
    print(StyleHelper.listSubinfo('${package.metadata.modes.length} mode(s)'));
    print(
      StyleHelper.listSubinfo(
        '${package.metadata.variables.length} variable(s)',
      ),
    );
    print(StyleHelper.listSubinfo('${package.bundles.length} bundle(s)'));
    print(StyleHelper.listDone('Package loaded'));
  }
}

