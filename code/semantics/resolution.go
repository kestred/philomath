package semantics

import "github.com/kestred/philomath/code/ast"
import "github.com/kestred/philomath/code/utils"

type ScopedName struct {
	Scope ast.Scope
	Name  string
}

func ResolveNames(cs *Section) {
	utils.Assert(!cs.DidSteps(Step_ResolveNames), "Tried to perform name resolution twice on the same code section")

	current := FindParentScope(cs.Root)
	var lookup = make(map[ScopedName]ast.Decl)
	for _, node := range cs.Nodes {
		// Track the current scope by always updating the current scope if we reach
		// a node and that node's parent provides a scope.
		//
		// It should mostly be correct, but I haven't thought about it that hard.
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
					utils.NotImplemented("Name resolution for out-of-order declarations")
				} else {
					search = FindParentScope(search)
				}
			}
		}
	}

	cs.StepsCompleted |= Step_ResolveNames
}

func FindParentScope(node ast.Node) ast.Scope {
	node = node.GetParent()
	for node != nil {
		if scope, ok := node.(ast.Scope); ok {
			return scope
		}
		node = node.GetParent()
	}

	return nil
}
