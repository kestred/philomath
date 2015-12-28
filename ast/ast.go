package ast

/* Abstract Nodes */

type Node interface {
	IsNode()
}

func (r *Root) IsNode() {}

// Declarations
func (d *ConstantDecl) IsNode() {}
func (d *MutableDecl) IsNode()  {}

// Definitions (and related nodes)
func (d *EnumDefn) IsNode()      {}
func (v *EnumValue) IsNode()     {}
func (s *EnumSeparator) IsNode() {}
func (d *StructDefn) IsNode()    {}
func (f *StructField) IsNode()   {}
func (d *ValueDefn) IsNode()     {}

// Expressions (and related nodes)
func (e *InfixExpr) IsNode()     {}
func (e *PrefixExpr) IsNode()    {}
func (e *CallExpr) IsNode()      {}
func (e *GroupExpr) IsNode()     {}
func (e *FunctionExpr) IsNode()  {}
func (p *FunctionParam) IsNode() {}
func (e *ValueExpr) IsNode()     {}

// Statements (and related nodes)
func (b *Block) IsNode()      {}
func (s *IfStmt) IsNode()     {}
func (s *WhileStmt) IsNode()  {}
func (s *ForStmt) IsNode()    {}
func (r *ForRange) IsNode()   {}
func (r *EachRange) IsNode()  {}
func (r *ExprRange) IsNode()  {}
func (s *AssignStmt) IsNode() {}

// Types
func (t *ArrayType) IsNode()    {}
func (t *FunctionType) IsNode() {}
func (t *NamedType) IsNode()    {}
func (t *PointerType) IsNode()  {}

// Literals
func (l *NumberLiteral) IsNode() {}
func (l *TextLiteral) IsNode()   {}
func (o *Operator) IsNode()      {}
func (i *ScopedIdent) IsNode()   {}
func (i *Ident) IsNode()         {}

type Decl interface {
	Node
	IsDecl()
}

func (d *ConstantDecl) IsDecl() {}
func (d *MutableDecl) IsDecl()  {}

type Defn interface {
	Node
	IsDefn()
}

func (d *EnumDefn) IsDefn()   {}
func (d *StructDefn) IsDefn() {}
func (d *ValueDefn) IsDefn()  {}

type Expr interface {
	Node
	IsExpr()
}

func (e *InfixExpr) IsExpr()    {}
func (e *PrefixExpr) IsExpr()   {}
func (e *CallExpr) IsExpr()     {}
func (e *GroupExpr) IsExpr()    {}
func (e *FunctionExpr) IsExpr() {}
func (e *ValueExpr) IsExpr()    {}

type Stmt interface {
	Node
	IsStmt()
}

func (b *Block) IsStmt()      {}
func (s *IfStmt) IsStmt()     {}
func (s *WhileStmt) IsStmt()  {}
func (s *ForStmt) IsStmt()    {}
func (s *AssignStmt) IsStmt() {}

type Type interface {
	Node
	IsType()
}

func (t *ArrayType) IsType()    {}
func (t *FunctionType) IsType() {}
func (t *NamedType) IsType()    {}
func (t *PointerType) IsType()  {}

type Literal interface {
	Node
	IsLiteral()
}

func (l *NumberLiteral) IsLiteral() {}
func (l *TextLiteral) IsLiteral()   {}
func (i *ScopedIdent) IsLiteral()   {}

type EnumItem interface {
	IsEnumItem()
}

func (d *EnumDefn) IsEnumItem()  {}
func (v *EnumValue) IsEnumItem() {}

type LoopRange interface {
	IsLoopRange()
}

func (r *ForRange) IsLoopRange()  {}
func (r *EachRange) IsLoopRange() {}

/* Concrete Nodes */

type Root struct {
	Stmts []Stmt
	Decls []Decl
}

// A declaration is represented by one of the following
type (
	ConstantDecl struct {
		Name Ident
		Defn Defn
	}

	MutableDecl struct {
		Name Ident
		Type Type
		Expr Expr
	}
)

// A definition is represented by a tree of one or more of the following
type (
	EnumDefn struct {
		Items []EnumItem
	}

	EnumValue struct {
		Name  Ident
		Value NumberLiteral
		Text  TextLiteral
	}

	EnumSeparator struct {
		Name Ident
	}

	StructDefn struct {
		Fields []StructField
	}

	StructField struct {
		Name Ident
		Type Type
	}

	ValueDefn struct {
		Expr Expr
	}
)

// An expression is represented by a tree of one or more of the following
type (
	InfixExpr struct {
		Left     Expr
		Operator Operator
		Right    Expr
	}

	PrefixExpr struct {
		Operator Operator
		Subexpr  Expr
	}

	CallExpr struct {
		Function  Expr
		Arguments []Expr
	}

	FunctionExpr struct {
		Params []FunctionParam
		Return Type
		Stmt   Stmt
	}

	FunctionParam struct {
		Name Ident
		Type Type
	}

	GroupExpr struct {
		Subexpr Expr
	}

	ValueExpr struct {
		Literal Literal
	}
)

// A statement is represented by a tree of one or more of the following
type (
	Block struct {
		Stmts []Stmt
	}

	IfStmt struct {
		Cond Expr
		Then Stmt
		Else Stmt
	}

	WhileStmt struct {
		Cond Expr
		Do   Stmt
	}

	ForStmt struct {
		Range LoopRange
		Expr  Expr
	}

	ForRange struct {
		Decl   MutableDecl
		Cond   Expr
		Update AssignStmt
	}

	EachRange struct {
		Names []Ident
		Expr  Expr
		Range ExprRange
	}

	ExprRange struct {
		Min Expr
		Max Expr
	}

	AssignStmt struct {
		Assignee ScopedIdent
		Operator Operator
		Value    Expr
	}
)

// A type is represented by a tree of one or more of the following
type (
	ArrayType struct {
		Length  Expr
		Element Type
	}

	FunctionType struct {
		Params []Type
		Return Type
		// TODO: Multiple return values?
		//Returns []Type
	}

	NamedType struct {
		Name ScopedIdent
	}

	PointerType struct {
		Reference Type
	}
)

// A literal is represented by one of the following
type (
	NumberLiteral struct {
		Literal string
	}

	TextLiteral struct {
		Literal string
	}

	ScopedIdent struct {
		Scope string
		Name  Ident
	}

	Ident struct {
		Literal string
	}

	Operator struct {
		Literal string
	}
)
