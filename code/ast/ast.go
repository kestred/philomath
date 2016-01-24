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

	// Relaxed types
	UninferredType   = &SimpleType{"<uninferred>"} // could be any type
	UnresolvedType   = &SimpleType{"<unresolved>"} // could not infer type
	UnknownType      = &SimpleType{"<unknown>"}    // could not infer type
	InferredNumber   = &SimpleType{"<number>"}     // could only be a number
	InferredFloat    = &SimpleType{"<float>"}      // could only be a float number
	InferredSigned   = &SimpleType{"<signed>"}     // could only be a signed number
	InferredUnsigned = &SimpleType{"<unsigned>"}   // could only be an unsigned number

	// Builtin types
	BuiltinEmpty   = &SimpleType{"empty"} // the 0-byte type
	BuiltinFloat   = &SimpleType{"float"}
	BuiltinFloat32 = &SimpleType{"f32"}
	BuiltinFloat64 = &SimpleType{"f64"}
	BuiltinInt     = &SimpleType{"int"}
	BuiltinInt8    = &SimpleType{"i8"}
	BuiltinInt16   = &SimpleType{"i16"}
	BuiltinInt32   = &SimpleType{"i32"}
	BuiltinInt64   = &SimpleType{"i64"}
	BuiltinUint    = &SimpleType{"uint"}
	BuiltinUint8   = &SimpleType{"u8"}
	BuiltinUint16  = &SimpleType{"u16"}
	BuiltinUint32  = &SimpleType{"u32"}
	BuiltinUint64  = &SimpleType{"u64"}

	// Logical operators
	BuiltinLogicalOr    = &OperatorDefn{"Logical Or", "or", "_or_", BinaryInfix, LeftAssociative, LogicalOrPrec}
	BuiltinLogicalAnd   = &OperatorDefn{"Logical And", "and", "_and_", BinaryInfix, LeftAssociative, LogicalAndPrec}
	BuiltinElementOf    = &OperatorDefn{"Element Of", "in", "_in_", BinaryInfix, NonAssociative, InclusionPrec}
	BuiltinNotElementOf = &OperatorDefn{"Not Element Of", "not in", "_not_in_", BinaryInfix, NonAssociative, InclusionPrec}

	// Comparison operators
	BuiltinIdentical      = &OperatorDefn{"Identical", "is", "_is_", BinaryInfix, NonAssociative, ComparisonPrec}
	BuiltinEqual          = &OperatorDefn{"Equal", "==", "_eq_", BinaryInfix, NonAssociative, ComparisonPrec}
	BuiltinLess           = &OperatorDefn{"Less", "<", "_lt_", BinaryInfix, NonAssociative, ComparisonPrec}
	BuiltinLessOrEqual    = &OperatorDefn{"Less Or Equal", "<=", "_lte_", BinaryInfix, NonAssociative, ComparisonPrec}
	BuiltinGreater        = &OperatorDefn{"Greater", ">", "_gt_", BinaryInfix, NonAssociative, ComparisonPrec}
	BuiltinGreaterOrEqual = &OperatorDefn{"Greater Or Equal", ">=", "_gte_", BinaryInfix, NonAssociative, ComparisonPrec}

	// Arithmetic operators
	BuiltinCompare   = &OperatorDefn{"Compare", "<=>", "_cmp_", BinaryInfix, LeftAssociative, ArithmeticPrec}
	BuiltinAdd       = &OperatorDefn{"Add", "+", "_add_", BinaryInfix, LeftAssociative, CommutativePrec}
	BuiltinSubtract  = &OperatorDefn{"Subtract", "-", "_sub_", BinaryInfix, LeftAssociative, CommutativePrec}
	BuiltinMultiply  = &OperatorDefn{"Multiply", "*", "_mul_", BinaryInfix, LeftAssociative, DistributivePrec}
	BuiltinDivide    = &OperatorDefn{"Divide", "/", "_div_", BinaryInfix, LeftAssociative, DistributivePrec}
	BuiltinRemainder = &OperatorDefn{"Remainder", "%", "_rem_", BinaryInfix, LeftAssociative, DistributivePrec}
	BuiltinPositive  = &OperatorDefn{"Positive", "+", "pos_", UnaryPrefix, RightAssociative, PrefixPrec}
	BuiltinNegative  = &OperatorDefn{"Negative", "-", "neg_", UnaryPrefix, RightAssociative, PrefixPrec}

	// Pointer operators
	BuiltinReference   = &OperatorDefn{"Reference", "^", "ref_", UnaryPrefix, RightAssociative, PostfixPrec}
	BuiltinDereference = &OperatorDefn{"Dereference", "~", "deref_", UnaryPrefix, RightAssociative, PostfixPrec}
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
func (d *ConstantDefn) ImplementsNode()  {}
func (d *OperatorDefn) ImplementsNode()  {}
func (d *StructDefn) ImplementsNode()    {}
func (f *StructField) ImplementsNode()   {}

// Statements (and related nodes)
func (s *IfStmt) ImplementsNode()     {}
func (s *WhileStmt) ImplementsNode()  {}
func (s *ForStmt) ImplementsNode()    {}
func (r *ForRange) ImplementsNode()   {}
func (r *EachRange) ImplementsNode()  {}
func (r *ExprRange) ImplementsNode()  {}
func (s *ExprStmt) ImplementsNode()   {}
func (s *AssignStmt) ImplementsNode() {}
func (s *ReturnStmt) ImplementsNode() {}
func (s *DoneStmt) ImplementsNode()   {}

// Expressions (and related nodes)
func (e *PostfixExpr) ImplementsNode()   {}
func (e *InfixExpr) ImplementsNode()     {}
func (e *PrefixExpr) ImplementsNode()    {}
func (e *CallExpr) ImplementsNode()      {}
func (e *GroupExpr) ImplementsNode()     {}
func (e *FunctionExpr) ImplementsNode()  {}
func (p *FunctionParam) ImplementsNode() {}
func (e *MemberExpr) ImplementsNode()    {}

// Literals
func (l *NumberLiteral) ImplementsNode() {}
func (l *TextLiteral) ImplementsNode()   {}
func (i *Identifier) ImplementsNode()    {}

// Types
func (t *ArrayType) ImplementsNode()    {}
func (t *FunctionType) ImplementsNode() {}
func (t *NamedType) ImplementsNode()    {}
func (t *PointerType) ImplementsNode()  {}
func (t *SimpleType) ImplementsNode()   {}

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
func (s *ExprStmt) ImplementsBlockable()     {}
func (s *AssignStmt) ImplementsBlockable()   {}
func (s *ReturnStmt) ImplementsBlockable()   {}
func (s *DoneStmt) ImplementsBlockable()     {}

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

func (d *EnumDefn) ImplementsDefn()     {}
func (d *ConstantDefn) ImplementsDefn() {}
func (d *OperatorDefn) ImplementsDefn() {}
func (d *StructDefn) ImplementsDefn()   {}

type Stmt interface {
	Blockable
	ImplementsStmt()
}

func (s *IfStmt) ImplementsStmt()     {}
func (s *WhileStmt) ImplementsStmt()  {}
func (s *ForStmt) ImplementsStmt()    {}
func (s *ExprStmt) ImplementsStmt()   {}
func (s *AssignStmt) ImplementsStmt() {}
func (s *ReturnStmt) ImplementsStmt() {}
func (s *DoneStmt) ImplementsStmt()   {}

type Expr interface {
	Node
	ImplementsExpr()
	GetType() Type
}

func (e *PostfixExpr) ImplementsExpr()   {}
func (e *InfixExpr) ImplementsExpr()     {}
func (e *PrefixExpr) ImplementsExpr()    {}
func (e *CallExpr) ImplementsExpr()      {}
func (e *GroupExpr) ImplementsExpr()     {}
func (e *FunctionExpr) ImplementsExpr()  {}
func (e *MemberExpr) ImplementsExpr()    {}
func (l *NumberLiteral) ImplementsExpr() {}
func (l *TextLiteral) ImplementsExpr()   {}
func (i *Identifier) ImplementsExpr()    {}

func (e *PostfixExpr) GetType() Type   { return e.Type }
func (e *InfixExpr) GetType() Type     { return e.Type }
func (e *PrefixExpr) GetType() Type    { return e.Type }
func (e *CallExpr) GetType() Type      { return e.Type }
func (e *GroupExpr) GetType() Type     { return e.Type }
func (e *FunctionExpr) GetType() Type  { return e.Type }
func (e *MemberExpr) GetType() Type    { return e.Type }
func (l *NumberLiteral) GetType() Type { return l.Type }
func (l *TextLiteral) GetType() Type   { return l.Type }
func (i *Identifier) GetType() Type    { return i.Type }

type Literal interface {
	Expr
	ImplementsLiteral()
}

func (l *NumberLiteral) ImplementsLiteral() {}
func (l *TextLiteral) ImplementsLiteral()   {}
func (i *Identifier) ImplementsLiteral()    {}

type Type interface {
	Node
	ImplementsType()
}

func (t *ArrayType) ImplementsType()    {}
func (t *FunctionType) ImplementsType() {}
func (t *NamedType) ImplementsType()    {}
func (t *PointerType) ImplementsType()  {}
func (t *SimpleType) ImplementsType()   {}

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
		Name *Identifier
		Defn Defn

		// semantic
		Type Type
	}

	MutableDecl struct {
		// syntactic
		Name *Identifier
		Type Type // <-- also semantic right now
		Expr Expr
	}
)

func Constant(name string, defn Defn) *ConstantDecl {
	return &ConstantDecl{
		Name: Ident(name),
		Defn: defn,
		Type: UninferredType,
	}
}

func Mutable(name string, typ Type, expr Expr) *MutableDecl {
	if typ == nil {
		typ = UninferredType
	}

	return &MutableDecl{
		Name: Ident(name),
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
		Name  *Identifier
		Value NumberLiteral
		Text  TextLiteral
	}

	EnumSeparator struct {
		Name *Identifier
	}

	ConstantDefn struct {
		Expr Expr
	}

	OperatorDefn struct {
		Name        string
		Literal     string
		Overload    string
		Type        OpType
		Associative OpAssociation
		Precedence  OpPrecedence
	}

	StructDefn struct {
		Fields []StructField
	}

	StructField struct {
		Name *Identifier
		Type Type
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
		Update AssignStmt
	}

	EachRange struct {
		Names []*Identifier
		Expr  Expr
		Range ExprRange
	}

	ExprRange struct {
		Min Expr
		Max Expr
	}

	AssignStmt struct {
		Assignees []Expr
		Operator  *OperatorDefn
		Values    []Expr
	}

	ExprStmt struct {
		Expr Expr
	}

	ReturnStmt struct {
		Value Expr
	}

	DoneStmt struct{}
)

// An expression is represented by a tree of one or more of the following
type (
	PostfixExpr struct {
		// syntax
		Subexpr  Expr
		Operator *OperatorDefn

		// semantics
		Type Type
	}

	InfixExpr struct {
		// syntax
		Left     Expr
		Operator *OperatorDefn
		Right    Expr

		// semantics
		Type Type
	}

	PrefixExpr struct {
		// syntax
		Operator *OperatorDefn
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
		Name *Identifier
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
		Left   Expr
		Member *Identifier

		// semantics
		Type Type
	}

	NumberLiteral struct {
		// syntax
		Literal string

		// semantics
		Type  Type
		Value Value
	}

	TextLiteral struct {
		// syntax
		Literal string

		// semantics
		Type  Type
		Value Value
	}

	Identifier struct {
		// syntax
		Literal string

		// semantics
		Type     Type
		Resolved Node
	}

	Value interface{}
)

func PostExp(subexpr Expr, op *OperatorDefn) *PostfixExpr {
	return &PostfixExpr{
		Subexpr:  subexpr,
		Operator: op,
		Type:     UninferredType,
	}
}

func InExp(left Expr, op *OperatorDefn, right Expr) *InfixExpr {
	return &InfixExpr{
		Left:     left,
		Operator: op,
		Right:    right,
		Type:     UninferredType,
	}
}

func PreExp(op *OperatorDefn, subexpr Expr) *PrefixExpr {
	return &PrefixExpr{
		Operator: op,
		Subexpr:  subexpr,
		Type:     UninferredType,
	}
}

func CallExp(fn Expr, args []Expr) *CallExpr {
	return &CallExpr{
		Function:  fn,
		Arguments: args,
		Type:      UninferredType,
	}
}

func FnExp(params []FunctionParam, ret Type, block Block) *FunctionExpr {
	return &FunctionExpr{
		Params: params,
		Return: ret,
		Block:  block,
		Type:   UninferredType,
	}
}

func GrpExp(subexpr Expr) *GroupExpr {
	return &GroupExpr{
		Subexpr: subexpr,
		Type:    UninferredType,
	}
}

func GetExp(left Expr, member string) *MemberExpr {
	return &MemberExpr{
		Left:   left,
		Member: Ident(member),
		Type:   UninferredType,
	}
}

func NumLit(literal string) *NumberLiteral {
	return &NumberLiteral{
		Literal: literal,
		Value:   UnparsedValue,
		Type:    UninferredType,
	}
}

func TxtLit(literal string) *TextLiteral {
	return &TextLiteral{
		Literal: literal,
		Value:   UnparsedValue,
		Type:    UninferredType,
	}
}

func Ident(literal string) *Identifier {
	return &Identifier{
		Literal: literal,
		Type:    UnresolvedType,
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

	SimpleType struct {
		Name string
	}
)

func NamTyp(name Expr) *NamedType {
	return &NamedType{Name: name}
}

func PtrTyp(pointerTo Type) *PointerType {
	return &PointerType{PointerTo: pointerTo}
}
