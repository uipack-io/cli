package uipack

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Package struct {
	Metadata BundleMetadata
	Bundles  []Bundle
}
type PackageEncoder struct {
	Value *Package
}

type PackageExtension string

const (
	JsonExtension   PackageExtension = ".json"
	BinaryExtension PackageExtension = ""
)

func (p *Package) GetBundle(variant Variant) *Bundle {
	for _, b := range p.Bundles {
		if b.Variant == variant {
			return &b
		}
	}
	return nil
}

func (p *PackageEncoder) Load(path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():

		// We first load metadata from the dedicated file
		meta_path := filepath.Join(path, "_")

		file, err := os.OpenFile(meta_path, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		r := bufio.NewReader(file)
		p.Value.Metadata.Decode(r)
		file.Close()

		// Then we load all of the bundles
		entries, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		bundles := make([]Bundle, 0)

		for _, e := range entries {
			if e.Type().IsRegular() {
				bundle_path := filepath.Join(path, e.Name())
				switch filepath.Ext(bundle_path) {
				case string(JsonExtension):

					switch e.Name() {
					case "_": // Metadata already loaded
						break
					default: // Bundle
						bundle := Bundle{}

						file, err := os.OpenFile(bundle_path, os.O_RDONLY, 0644)
						if err != nil {
							panic(err)
						}
						r := bufio.NewReader(file)
						bundle.Decode(r, &p.Value.Metadata)
						file.Close()

						bundles = append(bundles, bundle)
					}
				}
			}
		}
		p.Value.Bundles = bundles
	case mode.IsRegular():
		panic("zip archive support not implemented")
	}

}

func (p *PackageEncoder) Save(path string, extension PackageExtension) {

	// We first save metadata to the dedicated file
	file, err := os.OpenFile(filepath.Join(path, "_"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(file)
	p.Value.Metadata.Encode(w)
	w.Flush()
	file.Close()

	// Then we save all of the bundles
	for _, b := range p.Value.Bundles {
		file, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("%x", b.Variant)), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(file)
		b.Encode(w, &p.Value.Metadata)
		w.Flush()
		file.Close()
	}

}
