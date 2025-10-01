import 'metadata.dart';
import 'bundle.dart';
import 'variant.dart';

class Package {
  BundleMetadata metadata;
  List<Bundle> bundles;

  Package({required this.metadata, required this.bundles});

  Bundle? getBundle(Variant variant) {
    for (final bundle in bundles) {
      if (bundle.variant == variant) {
        return bundle;
      }
    }
    return null;
  }
}