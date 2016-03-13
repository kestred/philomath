package ast

import "github.com/kestred/philomath/code/utils"

type Section struct {
	Root     Node
	Nodes    []Node
	Parent   *Section
	Progress int
}

func FlattenTree(root Node, parent *Section) Section {
	// TODO: get parentNode properly
	var top Node = nil
	if parent != nil {
		top = parent.Root
	}

	nodes := flattenTree(root, top)
	return Section{Root: root, Nodes: nodes, Parent: parent}
}

func flattenTree(node Node, parent Node) []Node {
	node.SetParent(parent)
	nodes := []Node{node}
	switch n := node.(type) {
	case *TopScope:
		for _, decl := range n.Decls {
			nodes = append(nodes, flattenTree(decl, n)...)
		}
	case *Block:
		for _, subnode := range n.Nodes {
			nodes = append(nodes, flattenTree(subnode, n)...)
		}
	case *AsmBlock:
		for _, binding := range n.Inputs {
			nodes = append(nodes, binding.Name)
		}
		for _, binding := range n.Outputs {
			nodes = append(nodes, binding.Name)
		}

	// declarations
	case *ImmutableDecl:
		nodes = append(nodes, n.Name)
		nodes = append(nodes, flattenTree(n.Defn, n)...)
	case *MutableDecl:
		nodes = append(nodes, n.Name)
		nodes = append(nodes, flattenTree(n.Type, n)...)
		nodes = append(nodes, flattenTree(n.Expr, n)...)

	// definitions
	case *ConstantDefn:
		nodes = append(nodes, flattenTree(n.Expr, n)...)

	// statements
	case *EvalStmt:
		nodes = append(nodes, flattenTree(n.Expr, n)...)
	case *AssignStmt:
		for _, expr := range n.Left {
			nodes = append(nodes, flattenTree(expr, n)...)
		}
		for _, expr := range n.Right {
			nodes = append(nodes, flattenTree(expr, n)...)
		}

	// expressions
	case *PostfixExpr:
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)
		nodes = append(nodes, n.Operator)
	case *InfixExpr:
		nodes = append(nodes, flattenTree(n.Left, n)...)
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Right, n)...)
	case *PrefixExpr:
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)
	case *GroupExpr:
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)
	case *ProcedureExpr:
		nodes = append(nodes, n.Return)
		nodes = append(nodes, flattenTree(n.Block, n)...)

	// literals
	case *Identifier,
		*NumberLiteral,
		*TextLiteral:
		break // nothing to add

	// types
	case *BaseType:
		break // nothing to add

	default:
		utils.InvalidCodePath()
	}

	return nodes
}
