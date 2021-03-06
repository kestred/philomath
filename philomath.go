package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/bytecode"
	"github.com/kestred/philomath/code/interpreter"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
)

var ArgTrace = flag.Bool("trace", false, "")

func init() {
	log.SetFlags(0)
	log.SetPrefix("phi: ")
	flag.Parse()
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, `
Phi is the compiler for Philomath.

Usage:
  phi COMMAND [OPTIONS] [ARGS]

Commands:
  run     interpret a .phi source file
`[1:])
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	// handle command
	command := args[0]
	switch command {
	case "run":
		doRun(args[1:])
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

	psr := parser.Make(args[0], *ArgTrace, []byte(source))
	tree := psr.ParseTop()
	errcount := len(psr.Errors)
	if errcount > 0 {
		for _, err := range psr.Errors {
			fmt.Printf("%v\n", err)
		}

		if errcount >= parser.MaxErrors {
			log.Fatalf("aborted after the first %v errors...\n", errcount)
		} else {
			log.Fatalf("found %v syntax error(s)\n", errcount)
		}
	}

	for _, decl := range tree.Decls {
		if decl.GetName().Literal == "main" {
			if imm, ok := decl.(*ast.ImmutableDecl); ok {
				if con, ok := imm.Defn.(*ast.ConstantDefn); ok {
					if _, ok := con.Expr.(*ast.ProcedureExpr); ok {
						break
					}
				}

				log.Fatalf(`expected "main" to be a procedure (eg. "main :: () { ... }")`)
			} else {
				log.Fatalf(`your "main" procedure must use "::" instead of ":="`)
			}
		}
	}

	section := semantics.FlattenTree(tree, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	// TODO: maybe add an errors list to Section?
	errs := semantics.CheckTypes(&section)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Printf("%v\n", err)
		}
		log.Fatalf("found %v semantic error(s)\n", len(errs))
	}
	program := bytecode.NewProgram()
	program.Extend(tree)

	if _, ok := program.Text["main"]; !ok {
		log.Fatalf(`unable to find a procedure named "main"`)
	}

	interpreter.Run(program)
}
