package ast

/* Abstract Nodes */

type Node interface {
	ImplementsNode()
}

func (r *Root) ImplementsNode() {}

// Declarations
func (d *ConstantDecl) ImplementsNode() {}
func (d *MutableDecl) ImplementsNode()  {}

// Definitions (and related nodes)
func (d *EnumDefn) ImplementsNode()      {}
func (v *EnumValue) ImplementsNode()     {}
func (s *EnumSeparator) ImplementsNode() {}
func (d *StructDefn) ImplementsNode()    {}
func (f *StructField) ImplementsNode()   {}
func (d *ValueDefn) ImplementsNode()     {}

// Expressions (and related nodes)
func (e *PostfixExpr) ImplementsNode()   {}
func (e *InfixExpr) ImplementsNode()     {}
func (e *PrefixExpr) ImplementsNode()    {}
func (e *CallExpr) ImplementsNode()      {}
func (e *GroupExpr) ImplementsNode()     {}
func (e *FunctionExpr) ImplementsNode()  {}
func (p *FunctionParam) ImplementsNode() {}
func (e *ValueExpr) ImplementsNode()     {}

// Statements (and related nodes)
func (b *Block) ImplementsNode()      {}
func (s *IfStmt) ImplementsNode()     {}
func (s *WhileStmt) ImplementsNode()  {}
func (s *ForStmt) ImplementsNode()    {}
func (r *ForRange) ImplementsNode()   {}
func (r *EachRange) ImplementsNode()  {}
func (r *ExprRange) ImplementsNode()  {}
func (s *AssignStmt) ImplementsNode() {}
func (s *ReturnStmt) ImplementsNode() {}
func (s *DoneStmt) ImplementsNode()   {}

// Types
func (t *ArrayType) ImplementsNode()    {}
func (t *FunctionType) ImplementsNode() {}
func (t *NamedType) ImplementsNode()    {}
func (t *PointerType) ImplementsNode()  {}

// Literals
func (l *NumberLiteral) ImplementsNode() {}
func (l *TextLiteral) ImplementsNode()   {}
func (o *Operator) ImplementsNode()      {}
func (i *ScopedIdent) ImplementsNode()   {}
func (i *Ident) ImplementsNode()         {}

type Decl interface {
	Node
	ImplementsDecl()
}

func (d *ConstantDecl) ImplementsDecl() {}
func (d *MutableDecl) ImplementsDecl()  {}

type Defn interface {
	Node
	ImplementsDefn()
}

func (d *EnumDefn) ImplementsDefn()   {}
func (d *StructDefn) ImplementsDefn() {}
func (d *ValueDefn) ImplementsDefn()  {}

type Expr interface {
	Node
	ImplementsExpr()
}

func (e *PostfixExpr) ImplementsExpr()  {}
func (e *InfixExpr) ImplementsExpr()    {}
func (e *PrefixExpr) ImplementsExpr()   {}
func (e *CallExpr) ImplementsExpr()     {}
func (e *GroupExpr) ImplementsExpr()    {}
func (e *FunctionExpr) ImplementsExpr() {}
func (e *ValueExpr) ImplementsExpr()    {}

type Stmt interface {
	Node
	ImplementsStmt()
}

func (b *Block) ImplementsStmt()      {}
func (s *IfStmt) ImplementsStmt()     {}
func (s *WhileStmt) ImplementsStmt()  {}
func (s *ForStmt) ImplementsStmt()    {}
func (s *AssignStmt) ImplementsStmt() {}
func (s *ReturnStmt) ImplementsStmt() {}
func (s *DoneStmt) ImplementsStmt()   {}

type Type interface {
	Node
	ImplementsType()
}

func (t *ArrayType) ImplementsType()    {}
func (t *FunctionType) ImplementsType() {}
func (t *NamedType) ImplementsType()    {}
func (t *PointerType) ImplementsType()  {}

type Literal interface {
	Node
	ImplementsLiteral()
}

func (l *NumberLiteral) ImplementsLiteral() {}
func (l *TextLiteral) ImplementsLiteral()   {}
func (i *ScopedIdent) ImplementsLiteral()   {}

type EnumItem interface {
	Node
	ImplementsEnumItem()
}

func (d *EnumDefn) ImplementsEnumItem()  {}
func (v *EnumValue) ImplementsEnumItem() {}

type LoopRange interface {
	Node
	ImplementsLoopRange()
}

func (r *ForRange) ImplementsLoopRange()  {}
func (r *EachRange) ImplementsLoopRange() {}

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
	PostfixExpr struct {
		Subexpr  Expr
		Operator Operator
	}

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

	ReturnStmt struct {
		Value Expr
	}

	DoneStmt struct{}
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
