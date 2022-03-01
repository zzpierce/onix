package ast

import (
	"fmt"
	"onix/lex/scanner"
	"onix/lex/token"
)

type AstTreeBuilder struct {
	literals []string
	lp       int
	clit     string // current operating literal
	berr     error
}

func (at *AstTreeBuilder) Init(path string) error {
	sc, err := scanner.NewScanner(path)
	if err != nil {
		return err
	}
	for {
		lit, err := sc.Literal()
		if err != nil {
			return err
		}
		if lit == "" {
			break
		}
		at.literals = append(at.literals, lit)
	}
	return nil
}

func (at *AstTreeBuilder) next() token.Token {
	if at.lp >= len(at.literals) {
		return token.FILE_END
	}
	lit := at.literals[at.lp]
	fmt.Println("next --> " + lit)
	at.lp++
	at.clit = lit
	return token.Lookup(lit)
}

func (at *AstTreeBuilder) prop() token.Token {
	if at.lp >= len(at.literals) {
		return token.FILE_END
	}
	k := at.literals[at.lp]
	fmt.Println("prop --> " + k)
	return token.Lookup(k)
}

func (at *AstTreeBuilder) propv() (token.Token, string) {
	if at.lp >= len(at.literals) {
		return token.FILE_END, ""
	}
	k := at.literals[at.lp]
	return token.Lookup(k), k
}

func (at *AstTreeBuilder) BuildFile() (*AstFile, error) {
	// read first line
	t := at.next()
	if t != token.PACKAGE {
		return nil, expectError("package", at.clit)
	}
	pkg, err := at.bIdent()
	if err != nil {
		return nil, err
	}
	decls := make([]Decl, 0)
	for at.prop() != token.FILE_END {
		di, err := at.bDecl()
		if err != nil {
			return nil, err
		}
		decls = append(decls, di)
	}
	return &AstFile{
		Package: pkg,
		Decls:   decls,
	}, nil
}

func (at *AstTreeBuilder) bExpr() (Expr, error) {
	return nil, nil
}

func (at *AstTreeBuilder) bCallExpr() (*CallExpr, error) {
	return nil, nil
}

func (at *AstTreeBuilder) bSelectorExpr() (*SelectorExpr, error) {
	return nil, nil
}

func (at *AstTreeBuilder) bStringLit() (*BasicLit, error) {
	at.next()
	crune := []rune(at.clit)
	if crune[0] == '"' {
		return &BasicLit{
			Type:  token.STRING,
			Value: at.clit,
		}, nil
	}
	return nil, expectError("BasicLit Type=String", at.clit)
}

// an identity starts with a lit and followed by lits and decimals
func (at *AstTreeBuilder) bIdent() (*Ident, error) {
	t := at.next()
	if t == token.FILE_END {
		return nil, expectError("Ident", "endOfFile")
	}
	if !token.IsLiteral(at.clit) {
		return nil, expectError("Ident", at.clit)
	}
	return &Ident{
		Name: at.clit,
	}, nil
}

func (at *AstTreeBuilder) bDecl() (Decl, error) {
	// GenDecl starts with import/var/const
	tf, v := at.propv()
	if tf == token.IMPORT || tf == token.CONST || tf == token.VAR {
		return at.bGenDecl()
	}
	if tf == token.FUNC {
		return at.bFuncDecl()
	}
	return nil, unexpectError(v)
}

func (at *AstTreeBuilder) bGenDecl() (*GenDecl, error) {
	tok := at.next()
	group := false
	if at.prop() == token.LPAREN {
		at.next()
		group = true
	}
	specs := make([]Spec, 0)
	if tok == token.IMPORT {
		fsp, err := at.bImportSpec()
		if err != nil {
			return nil, err
		}
		specs = append(specs, fsp)
		if group {
			for at.prop() != token.RPAREN {
				sp, err := at.bImportSpec()
				if err != nil {
					return nil, err
				}
				specs = append(specs, sp)
			}
			// jump over the right parenthesis
			at.next()
		}
	}
	return &GenDecl{
		Type:  tok,
		Specs: specs,
	}, nil
}

func (at *AstTreeBuilder) bImportSpec() (*ImportSpec, error) {
	lit, err := at.bStringLit()
	if err != nil {
		return nil, err
	}
	return &ImportSpec{
		Path: lit,
	}, nil
}

func (at *AstTreeBuilder) bFuncDecl() (*FuncDecl, error) {
	at.next()
	fname, err := at.bIdent()
	if err != nil {
		return nil, err
	}
	// build function params and returns
	ft, err := at.bFuncType()
	if err != nil {
		return nil, err
	}
	body, err := at.bBlockStmt()
	if err != nil {
		return nil, err
	}
	return &FuncDecl{
		Name: fname,
		Type: ft,
		Body: body,
	}, nil
}

func (at *AstTreeBuilder) bFuncType() (*FuncType, error) {
	if at.next() != token.LPAREN {
		return nil, expectError("(", at.clit)
	}
	params, err := at.bFieldList()
	if err != nil {
		return nil, err
	}
	// jump over the right parenthesis
	at.next()

	returns := &FieldList{}
	// results have optional parenthesis
	if at.prop() == token.LPAREN {
		at.next()
		returns, err = at.bFieldList()
		if err != nil {
			return nil, err
		}
		at.next()
	} else if at.prop() != token.LBRACE {
		// single result mode
		retdt, err := at.bIdent()
		if err != nil {
			return nil, err
		}
		field := &Field{
			Names: []*Ident{retdt},
		}
		returns.Fields = []*Field{field}
	}
	return &FuncType{
		Params:  params,
		Returns: returns,
	}, nil
}

func (at *AstTreeBuilder) bFieldList() (*FieldList, error) {
	// build FieldList
	flist := make([]*Field, 0)
	fbuf := make([]*Ident, 0)
	for {
		if at.prop() == token.RPAREN {
			break
		}
		idt, err := at.bIdent()
		if err != nil {
			return nil, err
		}
		fbuf = append(fbuf, idt)
		sep := at.prop()
		if sep == token.COMMA {
			at.next()
			continue
		}
		typ, err := at.bIdent()
		if err != nil {
			return nil, err
		}
		fi := &Field{
			Names: fbuf,
			Type:  typ,
		}
		flist = append(flist, fi)
	}
	return &FieldList{
		Fields: flist,
	}, nil
}

func (at *AstTreeBuilder) bStmt() (Stmt, error) {
	return nil, nil
}

func (at *AstTreeBuilder) bExprStmt() (*ExprStmt, error) {
	expr, err := at.bExpr()
	if err != nil {
		return nil, err
	}
	return &ExprStmt{
		X: expr,
	}, nil
}

func (at *AstTreeBuilder) bBlockStmt() (*BlockStmt, error) {
	// block statement starts with '{' and ends with '}'
	if at.next() != token.LBRACE {
		return nil, expectError("{", at.clit)
	}
	sarr := make([]Stmt, 0)
	for {
		if at.prop() == token.RBRACE {
			at.next()
			break
		}
		stmt, err := at.bStmt()
		if err != nil {
			return nil, err
		}
		sarr = append(sarr, stmt)
	}
	return &BlockStmt{
		List: sarr,
	}, nil
}
