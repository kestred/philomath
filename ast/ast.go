package ast

/* At the moment this package isn't so far off from a concrete syntax tree

	TODO: Maybe it would be worth making the two structures separate?

  For example, the CST would only contain the syntactic information,
	and the AST would then only contain the semantic information.
	The AST could eliminate some nodes:
		For example "GroupExpr" can be removed entirely
		The "ValueExpr", "NumberLiteral", and "TextLiteral" nodes could be collapsed?

*/

/* Const Nodes */
var (
	UnparsedValue = &struct{}{}

	UnknownType      = &BuiltinType{"<unknown>"}  // could not infer type
	InferredType     = &BuiltinType{"<inferred>"} // could be any type
	InferredNumber   = &BuiltinType{"<number>"}   // could be any number
	InferredFloat    = &BuiltinType{"<float>"}    // could only be a float number
	InferredSigned   = &BuiltinType{"<signed>"}   // could only be a signed number
	InferredUnsigned = &BuiltinType{"<unsigned>"} // could only be an unsigned number
	BuiltinEmpty     = &BuiltinType{"empty"}      // the 0-byte type
	BuiltinFloat     = &BuiltinType{"float"}
	BuiltinFloat32   = &BuiltinType{"f32"}
	BuiltinFloat64   = &BuiltinType{"f64"}
	BuiltinInt       = &BuiltinType{"int"}
	BuiltinInt8      = &BuiltinType{"i8"}
	BuiltinInt16     = &BuiltinType{"i16"}
	BuiltinInt32     = &BuiltinType{"i32"}
	BuiltinInt64     = &BuiltinType{"i64"}
	BuiltinUint      = &BuiltinType{"uint"}
	BuiltinUint8     = &BuiltinType{"u8"}
	BuiltinUint16    = &BuiltinType{"u16"}
	BuiltinUint32    = &BuiltinType{"u32"}
	BuiltinUint64    = &BuiltinType{"u64"}
)

/* Abstract Nodes */

type Node interface {
	ImplementsNode()
}

func (b *Block) ImplementsNode() {}

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
func (e *AssignExpr) ImplementsNode()    {}
func (e *MemberExpr) ImplementsNode()    {}
func (e *ValueExpr) ImplementsNode()     {}

// Statements (and related nodes)
func (s *IfStmt) ImplementsNode()     {}
func (s *WhileStmt) ImplementsNode()  {}
func (s *ForStmt) ImplementsNode()    {}
func (r *ForRange) ImplementsNode()   {}
func (r *EachRange) ImplementsNode()  {}
func (r *ExprRange) ImplementsNode()  {}
func (s *ReturnStmt) ImplementsNode() {}
func (s *DoneStmt) ImplementsNode()   {}
func (s *ExprStmt) ImplementsNode()   {}

// Types
func (t *ArrayType) ImplementsNode()    {}
func (t *FunctionType) ImplementsNode() {}
func (t *NamedType) ImplementsNode()    {}
func (t *PointerType) ImplementsNode()  {}
func (t *BuiltinType) ImplementsNode()  {}

// Literals
func (l *NumberLiteral) ImplementsNode() {}
func (l *TextLiteral) ImplementsNode()   {}
func (o *Operator) ImplementsNode()      {}
func (i *Ident) ImplementsNode()         {}

type Blockable interface {
	Node
	ImplementsBlockable()
}

func (b *Block) ImplementsBlockable()        {}
func (d *ConstantDecl) ImplementsBlockable() {}
func (d *MutableDecl) ImplementsBlockable()  {}
func (s *IfStmt) ImplementsBlockable()       {}
func (s *WhileStmt) ImplementsBlockable()    {}
func (s *ForStmt) ImplementsBlockable()      {}
func (s *ReturnStmt) ImplementsBlockable()   {}
func (s *DoneStmt) ImplementsBlockable()     {}
func (s *ExprStmt) ImplementsBlockable()     {}

type Decl interface {
	Blockable
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

type Stmt interface {
	Blockable
	ImplementsStmt()
}

func (s *IfStmt) ImplementsStmt()     {}
func (s *WhileStmt) ImplementsStmt()  {}
func (s *ForStmt) ImplementsStmt()    {}
func (s *ReturnStmt) ImplementsStmt() {}
func (s *DoneStmt) ImplementsStmt()   {}
func (s *ExprStmt) ImplementsStmt()   {}

type Expr interface {
	Node
	ImplementsExpr()
	GetType() Type
}

func (e *PostfixExpr) ImplementsExpr()  {}
func (e *InfixExpr) ImplementsExpr()    {}
func (e *PrefixExpr) ImplementsExpr()   {}
func (e *CallExpr) ImplementsExpr()     {}
func (e *GroupExpr) ImplementsExpr()    {}
func (e *FunctionExpr) ImplementsExpr() {}
func (e *AssignExpr) ImplementsExpr()   {}
func (e *MemberExpr) ImplementsExpr()   {}
func (e *ValueExpr) ImplementsExpr()    {}

func (e *PostfixExpr) GetType() Type  { return e.Type }
func (e *InfixExpr) GetType() Type    { return e.Type }
func (e *PrefixExpr) GetType() Type   { return e.Type }
func (e *CallExpr) GetType() Type     { return e.Type }
func (e *GroupExpr) GetType() Type    { return e.Type }
func (e *FunctionExpr) GetType() Type { return e.Type }
func (e *AssignExpr) GetType() Type   { return e.Type }
func (e *MemberExpr) GetType() Type   { return e.Type }
func (e *ValueExpr) GetType() Type    { return e.Type }

type Type interface {
	Node
	ImplementsType()
}

func (t *ArrayType) ImplementsType()    {}
func (t *FunctionType) ImplementsType() {}
func (t *NamedType) ImplementsType()    {}
func (t *PointerType) ImplementsType()  {}
func (t *BuiltinType) ImplementsType()  {}

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

// A declaration is represented by one of the following
type (
	ConstantDecl struct {
		// syntactic
		Name Ident
		Defn Defn

		// semantic
		Type Type
	}

	MutableDecl struct {
		// syntactic
		Name Ident
		Type Type // <-- also semantic right now
		Expr Expr
	}
)

func Constant(name string, defn Defn) *ConstantDecl {
	return &ConstantDecl{
		Name: Ident{name},
		Defn: defn,
		Type: InferredType,
	}
}

func Mutable(name string, typ Type, expr Expr) *MutableDecl {
	if typ == nil {
		typ = InferredType
	}

	return &MutableDecl{
		Name: Ident{name},
		Type: typ,
		Expr: expr,
	}
}

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

// A statement is represented by a tree of one or more of the following
type (
	Block struct {
		Nodes []Blockable // Block, Decl, or Stmt
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
		Do    Block
	}

	ForRange struct {
		Decl   MutableDecl
		Cond   Expr
		Update AssignExpr
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

	ReturnStmt struct {
		Value Expr
	}

	DoneStmt struct{}

	ExprStmt struct {
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
		Block  Block

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

	AssignExpr struct {
		// syntax
		Assignee Expr
		Operator Operator
		Value    Expr

		// semantics
		Type Type
	}

	MemberExpr struct {
		// syntax
		Left   Expr
		Member Ident

		// semantics
		Type Type
	}

	ValueExpr struct {
		// syntax
		Literal Literal

		// semantics
		Type Type
	}
)

func PostExp(subexpr Expr, op Operator) *PostfixExpr {
	return &PostfixExpr{
		Subexpr:  subexpr,
		Operator: op,
		Type:     InferredType,
	}
}

func InExp(left Expr, op Operator, right Expr) *InfixExpr {
	return &InfixExpr{
		Left:     left,
		Operator: op,
		Right:    right,
		Type:     InferredType,
	}
}

func PreExp(op Operator, subexpr Expr) *PrefixExpr {
	return &PrefixExpr{
		Operator: op,
		Subexpr:  subexpr,
		Type:     InferredType,
	}
}

func CallExp(fn Expr, args []Expr) *CallExpr {
	return &CallExpr{
		Function:  fn,
		Arguments: args,
		Type:      InferredType,
	}
}

func FnExp(params []FunctionParam, ret Type, block Block) *FunctionExpr {
	return &FunctionExpr{
		Params: params,
		Return: ret,
		Block:  block,
		Type:   InferredType,
	}
}

func GrpExp(subexpr Expr) *GroupExpr {
	return &GroupExpr{
		Subexpr: subexpr,
		Type:    InferredType,
	}
}

func SetExp(assignee Expr, op Operator, value Expr) *AssignExpr {
	return &AssignExpr{
		Assignee: assignee,
		Operator: op,
		Value:    value,
		Type:     InferredType,
	}
}

func DotExp(left Expr, member Ident) *MemberExpr {
	return &MemberExpr{
		Left:   left,
		Member: member,
		Type:   InferredType,
	}
}

func ValExp(literal Literal) *ValueExpr {
	return &ValueExpr{
		Literal: literal,
		Type:    InferredType,
	}
}

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

	BuiltinType struct {
		Name string
	}
)

func NamTyp(name Expr) *NamedType {
	return &NamedType{Name: name}
}

func PtrTyp(pointerTo Type) *PointerType {
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

func NumLit(literal string) *NumberLiteral {
	return &NumberLiteral{
		Literal: literal,
		Value:   UnparsedValue,
	}
}

func TxtLit(literal string) *TextLiteral {
	return &TextLiteral{
		Literal: literal,
		Value:   UnparsedValue,
	}
}
