package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app             = kingpin.New("geode", "Compiler for the Geode Programming Language").Author(AUTHOR)
	dumpResult      = app.Flag("dump", "Print either llvm or ASM code after compiled (llvm by default, asm if --asm is passed)").Short('S').Bool()
	buildOutput     = app.Flag("output", "Output binary name.").Short('o').Default("a.out").String()
	optimize        = app.Flag("optimize", "Enable full optimization").Short('O').Bool()
	printVerbose    = app.Flag("verbose", "Enable verbose printing").Short('v').Bool()
	disableEmission = app.Flag("disable-emission", "Disable emission and only run through the syntax checking process").Bool()

	versionCMD = app.Command("version", "Display the version")

	buildCMD      = app.Command("build", "Build an executable.")
	buildInput    = buildCMD.Arg("input", "Geode source file or package").Default(".").String()
	emitASM       = buildCMD.Flag("asm", "Set the target to .s asm files with intel syntax instead of a single binary.").Bool()
	dumpScopeTree = buildCMD.Flag("dump-scope-tree", "Dump a tree representation of the scope to stdout").Bool()

	runCMD   = app.Command("run", "Build and run an executable, clean up afterwards").Default()
	runInput = runCMD.Arg("input", "Geode source file or package").String()
	runArgs  = runCMD.Arg("args", "Arguments to be passed into the program after building").Strings()

	testCMD = app.Command("test", "Run tests in the ./tests/ directory")

	newTestCMD  = app.Command("new-test", "Create a new test")
	newTestName = newTestCMD.Arg("name", "the name of the test").Required().String()

	cleanCMD = app.Command("clean", "Remove the hidden build directory")

	infoCMD   = app.Command("info", "Get information about a program (does not compile, just lexes and parses)")
	infoInput = infoCMD.Arg("input", "Geode source file or package").String()
)
