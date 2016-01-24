package ast

import "strconv"

type OpType uint8
type OpAssociation uint8
type OpPrecedence int8

const (
	Nullary OpType = iota
	UnaryPrefix
	UnaryPostfix
	BinaryInfix
)

var operatorTypes = [...]string{
	Nullary:      "Nullary",
	UnaryPrefix:  "Prefix",
	UnaryPostfix: "Postfix",
	BinaryInfix:  "Infix",
}

func (typ OpType) String() string {
	s := ""
	if 0 <= typ && typ < OpType(len(operatorTypes)) {
		s = operatorTypes[typ]
	}
	if s == "" {
		s = "OpType(" + strconv.Itoa(int(typ)) + ")"
	}
	return s
}

const (
	NonAssociative OpAssociation = iota
	LeftAssociative
	RightAssociative
)

var associatives = [...]string{
	NonAssociative:   "NonAssociative",
	LeftAssociative:  "LeftAssociative",
	RightAssociative: "RightAssociative",
}

func (asc OpAssociation) String() string {
	s := ""
	if 0 <= asc && asc < OpAssociation(len(associatives)) {
		s = associatives[asc]
	}
	if s == "" {
		s = "Associative(" + strconv.Itoa(int(asc)) + ")"
	}
	return s
}

const (
	AssignmentPrec   OpPrecedence = 0
	LogicalOrPrec    OpPrecedence = 10
	LogicalAndPrec   OpPrecedence = 15
	InclusionPrec    OpPrecedence = 20
	ComparisonPrec   OpPrecedence = 40
	ArithmeticPrec   OpPrecedence = 60
	CommutativePrec  OpPrecedence = 70
	DistributivePrec OpPrecedence = 80
	PrefixPrec       OpPrecedence = 100
	PostfixPrec      OpPrecedence = 120
	MaxPrecedence    OpPrecedence = 127
)
