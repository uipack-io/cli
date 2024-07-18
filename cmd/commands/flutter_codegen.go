package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	uipack "github.com/uipack-io/cli"
	"github.com/uipack-io/cli/cmd/base"
	"github.com/uipack-io/cli/codegen"
)

func FlutterCodeGen(output *string, p *uipack.Package) {

	fmt.Println(base.ListStarted("Flutter code generation to " + *output))

	dartFiles := make([]string, len(p.Bundles)+2)

	// Type definitions
	g := codegen.FlutterCodeGen{}
	code := g.GenerateDefinitions(&p.Metadata, &p.Bundles)
	path := filepath.Join(*output, "data.g.dart")
	dartFiles[0] = path
	saveDartFile(path, code)

	// Loader
	g = codegen.FlutterCodeGen{}
	code = g.GenerateBundleLoader(&p.Metadata)
	path = filepath.Join(*output, "loader.g.dart")
	dartFiles[1] = path
	saveDartFile(path, code)

	// Bundle definitions
	for i, bundle := range p.Bundles {
		identifier := fmt.Sprintf("%x", bundle.Variant)
		g := codegen.FlutterCodeGen{}
		code := g.GenerateBundle(&p.Metadata, &bundle)
		path := filepath.Join(*output, "bundle_"+identifier+".g.dart")
		dartFiles[i+2] = path
		saveDartFile(path, code)
	}

	fmt.Println(base.ListDone("Code generated!"))

	// Formatting the code with the Dart cm
	fmt.Println(base.ListStarted("Formatting Flutter code"))
	for _, file := range dartFiles {
		formatDartFile(file)
	}
	fmt.Println(base.ListDone("Code formatted!"))
}

func saveDartFile(path string, code string) {
	f, err := os.Create(path)
	fmt.Println(base.ListSubinfo(path))
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(code)
	if err != nil {
		panic(err)
	}
}

func formatDartFile(path string) {
	fmt.Println(base.ListSubinfo(path))
	cmd := exec.Command("dart", "format", path)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
