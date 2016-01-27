package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/bytecode"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/interpreter"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
)

// TODO: Replace with "optparse/argparse"-like CLI parser; the Go developers
// should not have chosen X style command-line arguments as the builtin library.
// The "pflag" library supports more conventional UNIX behavior but doesn't
// support all the familar command lines that users would expect.
import argparse "github.com/ogier/pflag"

var ArgTrace = argparse.Bool("trace", false, "")

func init() {
	log.SetFlags(0)
	log.SetPrefix("phi: ")
	argparse.Parse()
	argparse.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, `
Phi is an experimental compiler for AI research.

Usage:
  phi COMMAND [OPTIONS] [ARGS]

Commands:
  build   compile one or more files
  run     interpret the file or input stream
  shell   open an interactive philomath REPL
`[1:])
}

func main() {
	args := argparse.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	// handle command
	command := args[0]
	switch command {
	case "build":
		log.Fatalln("TODO: Implement compilation")
	case "run":
		doRun(args[1:])
	case "shell":
		log.Fatalln("TODO: Implement REPL")
	default:
		if len(command) > 14 {
			command = command[:10] + " ..."
		}
		log.Printf(`error: unknown command "%v"`, command)
		usage()
		os.Exit(1)
	}
}

func doRun(args []string) {
	if len(args) == 0 {
		log.Fatalln(`error: no input files`)
	}

	file, err := os.Open(args[0])
	if err != nil {
		log.Fatalln("error:", err)
	}

	source, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln("error:", err)
	}

	p := parser.Make(args[0], *ArgTrace, []byte(source))
	top := p.ParseTop()
	errcount := len(p.Errors)
	if errcount > 0 {
		for _, err := range p.Errors {
			fmt.Errorf("%v\n", err)
		}

		if errcount >= parser.MaxErrors {
			log.Fatalf("aborted after the first %v errors...\n", errcount)
		} else {
			log.Fatalf("found %v syntax error(s)\n", errcount)
		}
	}

	var inits []ast.Decl
	var mainProc *ast.ProcedureExpr
	for _, decl := range top.Decls {
		if decl.GetName().Literal == "main" {
			if imm, ok := decl.(*ast.ImmutableDecl); ok {
				if con, ok := imm.Defn.(*ast.ConstantDefn); ok {
					if proc, ok := con.Expr.(*ast.ProcedureExpr); ok {
						mainProc = proc
						continue
					}
				}

				log.Fatalf(`expected "main" to be a procedure (eg. "main :: () { ... }")`)
			} else {
				log.Fatalf(`your "main" procedure must use "::" instead of ":="`)
			}
		}

		inits = append(inits, decl)
	}

	if mainProc == nil {
		log.Fatalf(`unable to find a procedure named "main"`)
	}

	// initialize constants/globals
	programData := ast.Top(inits)
	initSection := code.PrepareTree(programData, nil)
	semantics.ResolveNames(&initSection)
	semantics.InferTypes(&initSection)
	program, scope := bytecode.Generate(programData)

	// bytecode for main procedure
	mainSection := code.PrepareTree(mainProc.Block, &initSection)
	semantics.ResolveNames(&mainSection)
	semantics.InferTypes(&mainSection)
	program.Extend(mainProc.Block, scope)

	// interpret bytecode
	temp := interpreter.Evaluate(program, scope.NextRegister)
	log.Println("result:", temp)
}
