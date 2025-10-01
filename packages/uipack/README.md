# UIpack Dart CLI

A Dart port of the UIpack CLI tool for processing UI package formats and generating code for various platforms.

## Features

- **Package Loading**: Load UIpack packages from directories or JSON files
- **Code Generation**: Generate Flutter/Dart code from UIpack definitions
- **Package Description**: Display detailed information about package contents
- **Binary Format Support**: Read and write UIpack binary format
- **Styled Terminal Output**: Beautiful terminal output with ANSI colors

## Installation

```bash
dart pub global activate --source path .
```

Or run directly:

```bash
dart run bin/main.dart <command> [arguments]
```

## Usage

### Basic Commands

```bash
# Display help
dart run bin/main.dart

# Describe a package
dart run bin/main.dart describe <package-path>

# Generate Flutter code
dart run bin/main.dart codegen --target flutter --output ./generated <package-path>
```

### Code Generation

The `codegen` command supports the following options:

- `--target`: Target platform (currently supports: `flutter`)
- `--output`: Output directory for generated code (default: `./flutter`)

Example:
```bash
dart run bin/main.dart codegen --target flutter --output ./lib/generated my_package.json
```

### Package Description

The `describe` command shows detailed information about a package:

```bash
dart run bin/main.dart describe my_package_directory/
```

This will display:
- Package metadata (name, version)
- Available modes and variants
- Variable definitions
- Bundle information in a formatted table

## Supported Formats

### Input Formats
- **Directory Structure**: UIpack package directories with metadata and bundle files
- **JSON Files**: UIpack JSON format (basic support)

### Output Formats
- **Flutter/Dart**: Generate type-safe Dart code for Flutter applications

## Architecture

The Dart CLI is structured as follows:

```
lib/src/
├── models/          # Data models (Package, Bundle, Color, etc.)
├── commands/        # Command implementations
├── helpers/         # Utility functions (styling, etc.)
├── encoding/        # Binary encoding/decoding
└── package_loader.dart  # Package loading logic
```

### Key Components

- **Models**: Type-safe representations of UIpack data structures
- **Binary Encoding**: Custom binary format support with ByteDataReader/Writer
- **Code Generation**: Extensible code generation system
- **Package Loader**: Handles loading from various sources

## Comparison with Go CLI

| Feature | Go CLI | Dart CLI | Status |
|---------|--------|----------|--------|
| Package Loading | ✅ | ✅ | Complete |
| Binary Format | ✅ | ✅ | Complete |
| Flutter Codegen | ✅ | ✅ | Complete |
| JSON Import | ✅ | 🚧 | Basic |
| Describe Command | ✅ | ✅ | Complete |
| Terminal Styling | ✅ (lipgloss) | ✅ (ANSI) | Complete |
| ZIP Archives | ❌ | ❌ | Not implemented |

## Development

### Running Tests

```bash
dart test
```

### Formatting

```bash
dart format .
```

### Analysis

```bash
dart analyze
```

## Contributing

This Dart CLI aims to maintain feature parity with the original Go implementation while leveraging Dart's type system and ecosystem.

## License

Same as the original UIpack CLI.