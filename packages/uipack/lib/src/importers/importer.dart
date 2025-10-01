import 'dart:typed_data';

import 'package:uipack/src/models/package.dart';

abstract class Importer {
  Package load(Uint8List bytes);
}
