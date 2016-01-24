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
	return Section{Root: root, Nodes: flattenTree(root), Parent: parent}
}

func flattenTree(node ast.Node) []ast.Node {
	nodes := []ast.Node{node}
	switch n := node.(type) {
	case *ast.Block:
		for _, subnode := range n.Nodes {
			nodes = append(nodes, flattenTree(subnode)...)
		}

	// Declarations
	case *ast.ConstantDecl:
		nodes = append(nodes, n.Name)
		nodes = append(nodes, flattenTree(n.Defn)...)
	case *ast.MutableDecl:
		nodes = append(nodes, n.Name)
		//nodes = append(nodes, flattenTree(n.Type)...)
		nodes = append(nodes, flattenTree(n.Expr)...)

	// Definitions
	case *ast.ExprDefn:
		nodes = append(nodes, flattenTree(n.Expr)...)

	// Statements
	case *ast.ExprStmt:
		nodes = append(nodes, flattenTree(n.Expr)...)
	case *ast.AssignStmt:
		for _, expr := range n.Assignees {
			nodes = append(nodes, flattenTree(expr)...)
		}
		for _, expr := range n.Values {
			nodes = append(nodes, flattenTree(expr)...)
		}

	// Expressions
	case *ast.PostfixExpr:
		nodes = append(nodes, flattenTree(n.Subexpr)...)
		nodes = append(nodes, n.Operator)
	case *ast.InfixExpr:
		nodes = append(nodes, flattenTree(n.Left)...)
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Right)...)
	case *ast.PrefixExpr:
		nodes = append(nodes, n.Operator)
		nodes = append(nodes, flattenTree(n.Subexpr)...)
	case *ast.GroupExpr:
		nodes = append(nodes, flattenTree(n.Subexpr)...)
	case *ast.ValueExpr:
		nodes = append(nodes, n.Literal)

	default:
		panic("TODO: Handle all nodes types")
		utils.InvalidCodePath()
	}

	return nodes
}
