package main

import (
	"io"
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

var ArgRun = argparse.BoolP("run", "r", false, "")
var ArgRepl = argparse.BoolP("interactive", "i", false, "")
var ArgTrace = argparse.Bool("trace", false, "")

func init() {
	log.SetFlags(0)
	log.SetPrefix("phi: ")
	argparse.Parse()
}

func main() {
	if *ArgRepl {
		log.Fatalf("you've discovered my dark secret... the lack of a REPL")
	}

	// TODO: Move file reading and handling to a proper home
	// TODO: Update parser to preform streaming parsing
	file := os.Stdin
	args := argparse.Args()
	if len(args) == 1 {
		var err error
		file, err = os.Open(args[0])
		if err != nil {
			log.Fatalf("unable to open file: %v", err)
		}
	} else if len(args) > 2 {
		log.Fatalln("too many arguments") // TODO: Usage text
	}

	source := make([]byte, 8192)
	_, err := file.Read(source)
	if err == io.EOF {
		panic("TODO: handle long inputs")
	} else if err != nil {
		log.Fatalln(err)
	}

	// if *ArgRun {
	p := parser.Make("example", *ArgTrace, []byte(``))
	top := p.ParseTop()
	errcount := len(p.Errors)
	if errcount > 0 {
		for _, err := range p.Errors {
			log.Println(err)
		}

		if errcount >= parser.MaxErrors {
			log.Fatalf("aborted after the first %v errors...\n", errcount)
		} else {
			log.Fatalf("found %v syntax error(s)\n", errcount)
		}
	}

	var inits []ast.Decl
	var main *ast.ProcedureExpr
	for _, decl := range top.Decls {
		if decl.GetName().Literal == "main" {
			if imm, ok := decl.(*ast.ImmutableDecl); ok {
				if con, ok := imm.Defn.(*ast.ConstantDefn); ok {
					if proc, ok := con.Expr.(*ast.ProcedureExpr); ok {
						main = proc
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

	// Initialize constants/globals
	programData := ast.Top(inits)
	initSection := code.PrepareTree(programData, nil)
	semantics.ResolveNames(&initSection)
	semantics.InferTypes(&initSection)
	mainSection := code.PrepareTree(main.Block, &initSection)
	semantics.ResolveNames(&mainSection)
	semantics.InferTypes(&mainSection)
	scope := bytecode.NewScope()
	insts := bytecode.Generate(top, scope)
	temp := interpreter.Evaluate(insts, scope.Constants, scope.NextRegister)
	log.Printf("Result: %v\n", temp)
	// } else {
	//    TODO: COMPILE AND LINK!
	// }
}
