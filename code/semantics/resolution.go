package semantics

import (
	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/utils"
)

type ScopedName struct {
	Scope ast.Scope
	Name  string
}

func FindParentScope(node ast.Node) ast.Scope {
	node = node.GetParent()
	for node != nil {
		if scope, ok := node.(ast.Scope); ok {
			return scope
		}
	}
	return nil
}

func ResolveNames(cs *code.Section) {
	_, ok := cs.Root.(ast.Scope)
	utils.Assert(ok, "ResolveNames currently expects Section.Root to be a lexical scope")

	var current ast.Scope
	var lookup = make(map[ScopedName]ast.Decl)
	for _, node := range cs.Nodes {
		// Easy way to track the current scope; should always be correct but not sure
		// Suppose an example tree...
		//
		//   BlockA
		//     Decl     -> parent is BlockA
		//     BlockB
		//       Decl   -> parent is BlockB
		//       Stmt   -> parent is BlockB
		//         Expr
		//     Decl     -> parent is BlockA
		//
		if scope, ok := node.GetParent().(ast.Scope); ok {
			current = scope
		}

		switch n := node.(type) {
		case *ast.ImmutableDecl:
			lookup[ScopedName{current, n.Name.Literal}] = n
		case *ast.MutableDecl:
			lookup[ScopedName{current, n.Name.Literal}] = n
		case *ast.Identifier:
			search := current
			for {
				if decl, ok := lookup[ScopedName{search, n.Literal}]; ok {
					n.Decl = decl
					break
				} else if search == nil {
					panic("TODO: out-of-order declaration lookup")
				} else {
					search = FindParentScope(search)
				}
			}
		}
	}
}
