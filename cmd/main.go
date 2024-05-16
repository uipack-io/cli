package main

import (
	"aloisdeniel/uipack"
	"aloisdeniel/uipack/cmd/base"
	"aloisdeniel/uipack/cmd/commands"
	"aloisdeniel/uipack/importers"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println(base.TitleStyle.Render("Welcome to uipack"))

	subcommand := os.Args[1]
	switch subcommand {
	case "codegen":
		codeGenCmd()
	case "describe":
		describe()
	default:
		fmt.Println("Unknown command")
	}
}

func loadPackage(cmd *flag.FlagSet) *uipack.Package {
	args := cmd.Args()
	input := args[len(args)-1]

	path, err := filepath.Abs(input)
	if err != nil {
		panic(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	fmt.Println(base.ListStarted("Loading package from " + path))
	p := uipack.Package{}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		e := uipack.PackageEncoder{Value: &p}
		e.Load(input)
	case mode.IsRegular():
		if filepath.Ext(path) == ".json" {
			encoded_json, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			importers.DecodeJson(&p, encoded_json)
		} else {
			panic("zip archive support not implemented")
		}
	}

	fmt.Println(base.ListSubinfo(p.Metadata.Name))
	fmt.Println(base.ListSubinfo(fmt.Sprint("Version ", p.Metadata.Version.Major, ".", p.Metadata.Version.Minor)))
	fmt.Println(base.ListSubinfo(fmt.Sprint(len(p.Metadata.Modes), " mode(s)")))
	fmt.Println(base.ListSubinfo(fmt.Sprint(len(p.Metadata.Variables), " variable(s)")))
	fmt.Println(base.ListSubinfo(fmt.Sprint(len(p.Bundles), " bundle(s)")))
	fmt.Println(base.ListDone("Package loaded"))

	return &p
}

func codeGenCmd() bool {
	codegenCmd := flag.NewFlagSet("codegen", flag.ExitOnError)
	targetPtr := codegenCmd.String("target", "flutter", "the target platform")
	output := codegenCmd.String("output", "./flutter", "the output directory")

	err := codegenCmd.Parse(os.Args[2:])

	if err != nil {
		return false
	}

	p := loadPackage(codegenCmd)
	switch *targetPtr {
	case "flutter":
		commands.FlutterCodeGen(output, p)
	}
	return true

}

func describe() bool {
	describeCmd := flag.NewFlagSet("describe", flag.ExitOnError)
	err := describeCmd.Parse(os.Args[2:])

	if err != nil {
		return false
	}

	p := loadPackage(describeCmd)
	commands.Describe(p)

	return true
}
