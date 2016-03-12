package ast

/* Constant Nodes */

var (
	UnparsedValue = &struct{}{}

	// Relaxed types
	InferredType     = BaseTyp("<inferrable>") // could be any type
	InferredText     = BaseTyp("<text>")       // could only be char/text/bytes/etc
	InferredNumber   = BaseTyp("<number>")     // could only be a number
	InferredFloat    = BaseTyp("<float>")      // could only be a float number
	InferredSigned   = BaseTyp("<signed>")     // could only be a signed number
	InferredUnsigned = BaseTyp("<unsigned>")   // could only be an unsigned number

	// Error types
	UninferredType  = BaseTyp("<uninferred>")  // used before it was inferred
	UnresolvedType  = BaseTyp("<unresolved>")  // could not infer type
	UncastableType  = BaseTyp("<uncastable>")  // could not cast type
	PlaceholderType = BaseTyp("<placeholder>") // a placeholder until I implement more complex types

	// Builtin types
	BuiltinEmpty   = BaseTyp("empty") // the 0-byte type
	BuiltinText    = BaseTyp("text")
	BuiltinChar    = BaseTyp("char")
	BuiltinFloat   = BaseTyp("float")
	BuiltinFloat32 = BaseTyp("f32")
	BuiltinFloat64 = BaseTyp("f64")
	BuiltinInt     = BaseTyp("int")
	BuiltinInt8    = BaseTyp("i8")
	BuiltinInt16   = BaseTyp("i16")
	BuiltinInt32   = BaseTyp("i32")
	BuiltinInt64   = BaseTyp("i64")
	BuiltinUint    = BaseTyp("uint")
	BuiltinUint8   = BaseTyp("u8")
	BuiltinUint16  = BaseTyp("u16")
	BuiltinUint32  = BaseTyp("u32")
	BuiltinUint64  = BaseTyp("u64")

	// Logical operators
	BuiltinLogicalOr    = Operator("Logical Or", "or", "_or_", BinaryInfix, LeftAssociative, LogicalOrPrec)
	BuiltinLogicalAnd   = Operator("Logical And", "and", "_and_", BinaryInfix, LeftAssociative, LogicalAndPrec)
	BuiltinElementOf    = Operator("Element Of", "in", "_in_", BinaryInfix, NonAssociative, InclusionPrec)
	BuiltinNotElementOf = Operator("Not Element Of", "not in", "_not_in_", BinaryInfix, NonAssociative, InclusionPrec)

	// Comparison operators
	BuiltinIdentical      = Operator("Identical", "is", "_is_", BinaryInfix, NonAssociative, ComparisonPrec)
	BuiltinEqual          = Operator("Equal", "==", "_eq_", BinaryInfix, NonAssociative, ComparisonPrec)
	BuiltinLess           = Operator("Less", "<", "_lt_", BinaryInfix, NonAssociative, ComparisonPrec)
	BuiltinLessOrEqual    = Operator("Less Or Equal", "<=", "_lte_", BinaryInfix, NonAssociative, ComparisonPrec)
	BuiltinGreater        = Operator("Greater", ">", "_gt_", BinaryInfix, NonAssociative, ComparisonPrec)
	BuiltinGreaterOrEqual = Operator("Greater Or Equal", ">=", "_gte_", BinaryInfix, NonAssociative, ComparisonPrec)

	// Arithmetic operators
	BuiltinCompare   = Operator("Compare", "<=>", "_cmp_", BinaryInfix, LeftAssociative, ArithmeticPrec)
	BuiltinAdd       = Operator("Add", "+", "_add_", BinaryInfix, LeftAssociative, CommutativePrec)
	BuiltinSubtract  = Operator("Subtract", "-", "_sub_", BinaryInfix, LeftAssociative, CommutativePrec)
	BuiltinMultiply  = Operator("Multiply", "*", "_mul_", BinaryInfix, LeftAssociative, DistributivePrec)
	BuiltinDivide    = Operator("Divide", "/", "_div_", BinaryInfix, LeftAssociative, DistributivePrec)
	BuiltinRemainder = Operator("Remainder", "%", "_rem_", BinaryInfix, LeftAssociative, DistributivePrec)
	BuiltinPositive  = Operator("Positive", "+", "pos_", UnaryPrefix, RightAssociative, PrefixPrec)
	BuiltinNegative  = Operator("Negative", "-", "neg_", UnaryPrefix, RightAssociative, PrefixPrec)

	// Pointer operators
	BuiltinReference   = Operator("Reference", "^", "ref_", UnaryPrefix, RightAssociative, PostfixPrec)
	BuiltinDereference = Operator("Dereference", "~", "deref_", UnaryPrefix, RightAssociative, PostfixPrec)
)

/* Abstract Nodes */

type Node interface {
	ImplementsNode()
	GetParent() Node
	SetParent(p Node)
}

type NodeBase struct {
	Parent Node
}

func (n *NodeBase) ImplementsNode()  {}
func (n *NodeBase) GetParent() Node  { return n.Parent }
func (n *NodeBase) SetParent(p Node) { n.Parent = p }

type Scope interface {
	Node
	ImplementsScope()
}

func (t *TopScope) ImplementsScope() {}
func (b *Block) ImplementsScope()    {}

type Evaluable interface {
	Node
	ImplementsEvaluable()
}

func (b *Block) ImplementsEvaluable()         {}
func (s *AsmBlock) ImplementsEvaluable()      {}
func (d *ImmutableDecl) ImplementsEvaluable() {}
func (d *MutableDecl) ImplementsEvaluable()   {}
func (s *IfStmt) ImplementsEvaluable()        {}
func (s *WhileStmt) ImplementsEvaluable()     {}
func (s *ForStmt) ImplementsEvaluable()       {}
func (s *EvalStmt) ImplementsEvaluable()      {}
func (s *AssignStmt) ImplementsEvaluable()    {}
func (s *ReturnStmt) ImplementsEvaluable()    {}
func (s *DoneStmt) ImplementsEvaluable()      {}

type Decl interface {
	Evaluable
	ImplementsDecl()
	GetName() *Identifier
}

func (d *ImmutableDecl) ImplementsDecl() {}
func (d *MutableDecl) ImplementsDecl()   {}

func (d *ImmutableDecl) GetName() *Identifier { return d.Name }
func (d *MutableDecl) GetName() *Identifier   { return d.Name }

type Defn interface {
	Node
	ImplementsDefn()
}

func (d *EnumDefn) ImplementsDefn()     {}
func (d *ConstantDefn) ImplementsDefn() {}
func (d *OperatorDefn) ImplementsDefn() {}
func (d *StructDefn) ImplementsDefn()   {}

type Stmt interface {
	Evaluable
	ImplementsStmt()
}

func (s *IfStmt) ImplementsStmt()     {}
func (s *WhileStmt) ImplementsStmt()  {}
func (s *ForStmt) ImplementsStmt()    {}
func (s *EvalStmt) ImplementsStmt()   {}
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
func (e *ProcedureExpr) ImplementsExpr() {}
func (e *MemberExpr) ImplementsExpr()    {}
func (l *NumberLiteral) ImplementsExpr() {}
func (l *TextLiteral) ImplementsExpr()   {}
func (i *Identifier) ImplementsExpr()    {}

func (e *PostfixExpr) GetType() Type   { return e.Type }
func (e *InfixExpr) GetType() Type     { return e.Type }
func (e *PrefixExpr) GetType() Type    { return e.Type }
func (e *CallExpr) GetType() Type      { return e.Type }
func (e *GroupExpr) GetType() Type     { return e.Type }
func (e *ProcedureExpr) GetType() Type { return e.Type }
func (e *MemberExpr) GetType() Type    { return e.Type }
func (l *NumberLiteral) GetType() Type { return l.Type }
func (l *TextLiteral) GetType() Type   { return l.Type }
func (i *Identifier) GetType() Type    { return i.Type }

type Literal interface {
	Expr
	ImplementsLiteral()
	GetValue() Value
}

func (l *NumberLiteral) ImplementsLiteral() {}
func (l *TextLiteral) ImplementsLiteral()   {}

func (l *NumberLiteral) GetValue() Value { return l.Value }
func (l *TextLiteral) GetValue() Value   { return l.Value }

type Type interface {
	Node
	ImplementsType()
}

func (t *ArrayType) ImplementsType()     {}
func (t *ProcedureType) ImplementsType() {}
func (t *NamedType) ImplementsType()     {}
func (t *PointerType) ImplementsType()   {}
func (t *BaseType) ImplementsType()      {}

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

type TopScope struct {
	NodeBase
	Decls []Decl
}

func Top(decls []Decl) *TopScope {
	return &TopScope{Decls: decls}
}

type Block struct {
	NodeBase
	Nodes []Evaluable // Block, Decl, or Stmt
}

func Blok(nodes []Evaluable) *Block {
	return &Block{Nodes: nodes}
}

type AsmBlock struct {
	NodeBase

	// syntax
	Source string

	// semantics
	Inputs  []AsmBinding
	Outputs []AsmBinding
}

type AsmBinding struct {
	Name   *Identifier
	Offset int
}

func Asm(source string) *AsmBlock {
	return &AsmBlock{Source: source}
}

// A declaration is represented by one of the following
type (
	ImmutableDecl struct {
		NodeBase

		// syntax
		Name *Identifier
		Defn Defn

		// semantic
		Type Type
	}

	MutableDecl struct {
		NodeBase

		// syntax
		Name *Identifier
		Type Type // <-- also semantic right now
		Expr Expr
	}
)

func Immutable(name string, defn Defn) *ImmutableDecl {
	return &ImmutableDecl{
		Name: Ident(name),
		Defn: defn,
		Type: InferredType,
	}
}

func Mutable(name string, typ Type, expr Expr) *MutableDecl {
	if typ == nil {
		typ = InferredType
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
		NodeBase

		// syntax
		Items []EnumItem
	}

	EnumValue struct {
		NodeBase

		// syntax
		Name  *Identifier
		Value NumberLiteral
		Text  TextLiteral
	}

	EnumSeparator struct {
		NodeBase

		// syntax
		Name *Identifier
	}

	ConstantDefn struct {
		NodeBase

		// syntax
		Expr Expr
	}

	OperatorDefn struct {
		NodeBase

		// syntax
		Literal     string
		Overload    string
		Type        OpType
		Associative OpAssociation
		Precedence  OpPrecedence

		// semantics
		Name string
	}

	StructDefn struct {
		NodeBase

		// syntax
		Fields []StructField
	}

	StructField struct {
		NodeBase

		// syntax
		Name *Identifier
		Type Type
	}
)

func Constant(expr Expr) *ConstantDefn {
	return &ConstantDefn{Expr: expr}
}

func Operator(name string, lit string, ident string, typ OpType, asc OpAssociation, prec OpPrecedence) *OperatorDefn {
	return &OperatorDefn{
		Name:        name,
		Literal:     lit,
		Overload:    ident,
		Type:        typ,
		Associative: asc,
		Precedence:  prec,
	}
}

// A statement is represented by a tree of one or more of the following
type (
	IfStmt struct {
		NodeBase

		// syntax
		Cond Expr
		Then Block
		Else Block
	}

	WhileStmt struct {
		NodeBase

		// syntax
		Cond Expr
		Do   Block
	}

	ForStmt struct {
		NodeBase

		// syntax
		Range LoopRange
		Do    Block
	}

	ForRange struct {
		NodeBase

		// syntax
		Decl   MutableDecl
		Cond   Expr
		Update AssignStmt
	}

	EachRange struct {
		NodeBase

		// syntax
		Names []*Identifier
		Expr  Expr
		Range ExprRange
	}

	ExprRange struct {
		NodeBase

		// syntax
		Min Expr
		Max Expr
	}

	AssignStmt struct {
		NodeBase

		// syntax
		Left     []Expr
		Operator *OperatorDefn
		Right    []Expr
	}

	EvalStmt struct {
		NodeBase

		// syntax
		Expr Expr
	}

	ReturnStmt struct {
		NodeBase

		// syntax
		Value Expr
	}

	DoneStmt struct {
		NodeBase
	}
)

func Assign(left []Expr, op *OperatorDefn, right []Expr) *AssignStmt {
	return &AssignStmt{
		Left:     left,
		Operator: op,
		Right:    right,
	}
}

func Eval(expr Expr) *EvalStmt {
	return &EvalStmt{Expr: expr}
}

// An expression is represented by a tree of one or more of the following
type (
	PostfixExpr struct {
		NodeBase

		// syntax
		Subexpr  Expr
		Operator *OperatorDefn

		// semantics
		Type Type
	}

	InfixExpr struct {
		NodeBase

		// syntax
		Left     Expr
		Operator *OperatorDefn
		Right    Expr

		// semantics
		Type Type
	}

	PrefixExpr struct {
		NodeBase

		// syntax
		Operator *OperatorDefn
		Subexpr  Expr

		// semantics
		Type Type
	}

	CallExpr struct {
		NodeBase

		// syntax
		Procedure Expr
		Arguments []Expr

		// semantics
		Type Type
	}

	ProcedureExpr struct {
		NodeBase

		// syntax
		Params []ProcedureParam
		Return Type
		Block  *Block

		// semantics
		Type Type
	}

	ProcedureParam struct {
		NodeBase

		// syntax
		Name *Identifier
		Type Type
	}

	GroupExpr struct {
		NodeBase

		// syntax
		Subexpr Expr

		// semantics
		Type Type
	}

	MemberExpr struct {
		NodeBase

		// syntax
		Left   Expr
		Member *Identifier

		// semantics
		Type Type
	}

	NumberLiteral struct {
		NodeBase

		// syntax
		Literal string

		// semantics
		Type  Type
		Value Value
	}

	TextLiteral struct {
		NodeBase

		// syntax
		Literal string

		// semantics
		Type  Type
		Value Value
	}

	Identifier struct {
		NodeBase

		// syntax
		Literal string

		// semantics
		Type Type
		Decl Decl
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

func CallExp(proc Expr, args []Expr) *CallExpr {
	return &CallExpr{
		Procedure: proc,
		Arguments: args,
		Type:      UninferredType,
	}
}

func ProcExp(params []ProcedureParam, ret Type, block *Block) *ProcedureExpr {
	if ret == nil {
		ret = InferredType
	}

	return &ProcedureExpr{
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
		NodeBase

		// syntax
		Length  Expr
		Element Type
	}

	ProcedureType struct {
		NodeBase

		// syntax
		Params []Type
		Return Type
	}

	NamedType struct {
		NodeBase

		// syntax
		Name Expr
	}

	PointerType struct {
		NodeBase

		// syntax
		PointerTo Type
	}

	BaseType struct {
		NodeBase
		Name string
	}
)

func NamTyp(name Expr) *NamedType {
	return &NamedType{Name: name}
}

func PtrTyp(pointerTo Type) *PointerType {
	return &PointerType{PointerTo: pointerTo}
}

func BaseTyp(name string) *BaseType {
	return &BaseType{Name: name}
}
