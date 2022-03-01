package ast

import "onix/lex/token"

// build ast tree
// and print all found grammer errors.

type Expr interface {
	expr()
}

type Stmt interface {
	stmt()
}

type Decl interface {
	decl()
}

type Spec interface {
	spec()
}

type AstFile struct {
	Package *Ident
	Decls   []Decl
}

type Field struct {
	Names []*Ident
	Type  Expr
	Tag   *BasicLit
}

type FieldList struct {
	Fields []*Field
}

// Expr
type Ident struct {
	Name string
}

type BasicLit struct {
	Type  token.Token
	Value string
}

type FuncType struct {
	Params  *FieldList
	Returns *FieldList
}

type CallExpr struct {
	Fun  *SelectorExpr
	Args []Expr
}

type SelectorExpr struct {
	X   *Ident
	Sel *Ident
}

func (*Ident) expr()    {}
func (*BasicLit) expr() {}
func (*FuncType) expr() {}

// Stmt
type BlockStmt struct {
	List []Stmt
}

type ExprStmt struct {
	X Expr
}

// Spec
type ImportSpec struct {
	Path *BasicLit
}

func (*ImportSpec) spec() {}

// Decl
type GenDecl struct {
	Type  token.Token
	Specs []Spec
}

type FuncDecl struct {
	Name *Ident
	Type *FuncType
	Body *BlockStmt
}

func (*GenDecl) decl()  {}
func (*FuncDecl) decl() {}
