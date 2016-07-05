package semantics

import "github.com/kestred/philomath/code/ast"

// TODO: implement most types of type checking (so far I only check # of arguments in procedure calls)

func CheckTypes(cs *ast.Section) []error {
	// TODO: assert that name resolution and inference have already been run
	// TODO: move error printing out from parser and into its own module; then add source ranges to AST nodes
	return nil
}
