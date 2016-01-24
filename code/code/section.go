package code

import (
	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/utils"
)

type Section struct {
	Root     ast.Node
	Nodes    []ast.Node
	Parent   *Section
	Progress int
}

func PrepareTree(root ast.Node, parent *Section) Section {
	nodes := flattenTree(root, nil) // TODO: Get parentNode from parent section
	return Section{Root: root, Nodes: nodes, Parent: parent}
}

func flattenTree(node ast.Node, parent ast.Node) []ast.Node {
	node.SetParent(parent)
	nodes := []ast.Node{node}
	switch n := node.(type) {
	case *ast.Block:
		for _, subnode := range n.Nodes {
			nodes = append(nodes, flattenTree(subnode, n)...)
		}

	// Declarations
	case *ast.ConstantDecl:
		nodes = append(nodes, n.Name)
		nodes = append(nodes, flattenTree(n.Defn, n)...)
	case *ast.MutableDecl:
		nodes = append(nodes, n.Name)
		//nodes = append(nodes, flattenTree(n.Type, n)...)
		nodes = append(nodes, flattenTree(n.Expr, n)...)

	// Definitions
	case *ast.ConstantDefn:
		nodes = append(nodes, flattenTree(n.Expr, n)...)

	// Statements
	case *ast.EvalStmt:
		nodes = append(nodes, flattenTree(n.Expr, n)...)
	case *ast.AssignStmt:
		for _, expr := range n.Left {
			nodes = append(nodes, flattenTree(expr, n)...)
		}
		for _, expr := range n.Right {
			nodes = append(nodes, flattenTree(expr, n)...)
		}

	// Expressions
	case *ast.PostfixExpr:
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)
		nodes = append(nodes, n.Operator)
	case *ast.InfixExpr:
		nodes = append(nodes, flattenTree(n.Left, n)...)
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Right, n)...)
	case *ast.PrefixExpr:
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)
	case *ast.GroupExpr:
		nodes = append(nodes, flattenTree(n.Subexpr, n)...)

	// Literals
	case *ast.Identifier,
		*ast.NumberLiteral,
		*ast.TextLiteral:
		break // nothing to add

	default:
		panic("TODO: Handle all nodes types")
		utils.InvalidCodePath()
	}

	return nodes
}