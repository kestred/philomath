package ast

/* Null Nodes */
var Undefined = &struct{}{}
var Inferred = &NamedType{Name: &ValueExpr{Literal: &Ident{"<Inferred>"}}}

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
func (e *MemberExpr) ImplementsNode()    {}
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
func (e *MemberExpr) ImplementsExpr()   {}
func (e *ValueExpr) ImplementsExpr()    {}

type Stmt interface {
	Node
	ImplementsStmt()
}

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
func (i *Ident) ImplementsLiteral()         {}

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
		// syntax
		Subexpr  Expr
		Operator Operator

		// semantics
		Type Type
	}

	InfixExpr struct {
		// syntax
		Left     Expr
		Operator Operator
		Right    Expr

		// semantics
		Type Type
	}

	PrefixExpr struct {
		// syntax
		Operator Operator
		Subexpr  Expr

		// semantics
		Type Type
	}

	CallExpr struct {
		// syntax
		Function  Expr
		Arguments []Expr

		// semantics
		Type Type
	}

	FunctionExpr struct {
		// syntax
		Params []FunctionParam
		Return Type
		Block Block

		// semantics
		Type Type
	}

	FunctionParam struct {
		Name Ident
		Type Type
	}

	GroupExpr struct {
		// syntax
		Subexpr Expr

		// semantics
		Type Type
	}

	MemberExpr struct {
		// syntax
		Left Expr
		Member Ident

		// semantics
		Type Type
	}

	ValueExpr struct {
		// syntax
		Literal Literal

		// semantics
		Type  Type
	}
)

func NewPostfixExpr(subexpr Expr, op Operator) *PostfixExpr {
	return &PostfixExpr{
		Subexpr: subexpr,
		Operator: op,
		Type: Inferred,
	}
}

func NewInfixExpr(left Expr, op Operator, right Expr) *InfixExpr {
	return &InfixExpr{
		Left: left,
		Operator: op,
		Right: right,
		Type: Inferred,
	}
}

func NewPrefixExpr(op Operator, subexpr Expr) *PrefixExpr {
	return &PrefixExpr{
		Operator: op,
		Subexpr: subexpr,
		Type: Inferred,
	}
}

func NewCallExpr(fn Expr, args []Expr) *CallExpr {
	return &CallExpr{
		Function: fn,
		Arguments: args,
		Type: Inferred,
	}
}

func NewFunctionExpr(params []FunctionParam, ret Type, block Block) *FunctionExpr {
	return &FunctionExpr {
		Params: params,
		Return: ret,
		Block: block,
		Type: Inferred,
	}
}

func NewGroupExpr(subexpr Expr) *GroupExpr {
	return &GroupExpr{
		Subexpr: subexpr,
		Type: Inferred,
	}
}

func NewMemberExpr(left Expr, member Ident) *MemberExpr {
	return &MemberExpr{
		Left: left,
		Member: member,
		Type: Inferred,
	}
}

func NewValueExpr(literal Literal) *ValueExpr {
	return &ValueExpr{
		Literal: literal,
		Type: Inferred,
	}
}

// A statement is represented by a tree of one or more of the following
type (
	Block struct {
		Stmts []Stmt
	}

	IfStmt struct {
		Cond Expr
		Then Block
		Else Block
	}

	WhileStmt struct {
		Cond Expr
		Do   Block
	}

	ForStmt struct {
		Range LoopRange
		Do Block
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
		Assignee Expr
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
		Name Expr
	}

	PointerType struct {
		PointerTo Type
	}
)

func NewNamedType(name Expr) *NamedType {
	return &NamedType{Name: name}
}

func NewPointerType(pointerTo Type) *PointerType {
	return &PointerType{PointerTo: pointerTo}
}

// A literal is represented by one of the following
type (
	NumberLiteral struct {
		// syntax
		Literal string

		// semantics
		Value interface{}
	}

	TextLiteral struct {
		// syntax
		Literal string

		// semantics
		Value interface{}
	}

	Ident struct {
		Literal string
	}

	Operator struct {
		Literal string
	}
)

func NewNumberLiteral(literal string) *NumberLiteral {
	return &NumberLiteral{
		Literal: literal,
		Value: Undefined,
	}
}

func NewTextLiteral(literal string) *TextLiteral {
	return &TextLiteral{
		Literal: literal,
		Value: Undefined,
	}
}
