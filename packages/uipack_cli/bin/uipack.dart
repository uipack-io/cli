import 'dart:io';
import 'package:args/args.dart';
import 'package:path/path.dart' as path;
import 'package:uipack/uipack.dart';

void main(List<String> arguments) {
  final parser = ArgParser()
    ..addCommand('codegen')
    ..addCommand('describe');

  try {
    final results = parser.parse(arguments);

    if (results.command == null) {
      print('Welcome to uipack');
      print('Usage: uipack_dart <command> [arguments]');
      print('\nAvailable commands:');
      print('  codegen    Generate code for target platform');
      exit(0);
    }

    final command = results.command!;

    switch (command.name) {
      case 'codegen':
        _handleCodegen(command);
        break;
      default:
        print('Unknown command: ${command.name}');
        exit(1);
    }
  } catch (e) {
    print('Error: $e');
    exit(1);
  }
}

void _handleCodegen(ArgResults command) {
  final codegenParser = ArgParser()
    ..addOption('target', defaultsTo: 'flutter', help: 'The target platform')
    ..addOption(
      'output',
      defaultsTo: './lib/src/theme',
      help: 'The output directory',
    );

  final codegenResults = codegenParser.parse(command.arguments);

  if (codegenResults.rest.isEmpty) {
    print('Error: Please provide a package path');
    exit(1);
  }

  final packagePath = codegenResults.rest.last;
  final target = codegenResults['target'] as String;
  final outputPath = codegenResults['output'] as String;

  final package = PackageLoader.loadPackage(packagePath);
  switch (target) {
    case 'flutter':
      final outputDir = Directory(outputPath);
      outputDir.createSync(recursive: true);

      final gen = FlutterCodeGen();

      final file = File(path.join(outputDir.path, 'data.g.dart'));
      file.createSync(recursive: true);
      final code = gen.generateDefinitions(package.metadata, package.bundles);
      file.writeAsStringSync(code);
      print('Generated ${file.path}');

      for (final bundle in package.bundles) {
        final file = File(
            path.join(outputDir.path, 'bundle_${bundle.variant.value}.g.dart'));
        file.createSync(recursive: true);
        final code = gen.generateBundle(package.metadata, bundle);
        file.writeAsStringSync(code);
        print('Generated ${file.path}');
      }

    default:
      print('Error: Unsupported target: $target');
      exit(1);
  }
}

