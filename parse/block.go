package parse

import (
	"github.com/socialmachines/soma/ast"
	"github.com/socialmachines/soma/rt"
	"github.com/socialmachines/soma/scan"
)

// block :=
//	LBRACE [arguments] [statements] RBRACE
//
func (p *Parser) parseBlock() (b *ast.Block) {
	p.expect(scan.LBRACE)

	b = &ast.Block{}
	if p.tok == scan.BINARY && p.lit == "|" {
		b.Args = p.parseBlockArguments()
	} else {
		b.Args = []string{}
	}
	var stmts []rt.Expr
	b.Statements = p.parseStatements(stmts)

	p.expect(scan.RBRACE)
	return
}

// arguments :=
//   '|' IDENT (PERIOD IDENT)* '|'
//
func (p *Parser) parseBlockArguments() []string {
	p.expect(scan.BINARY)

	args := []string{p.expect(scan.IDENT)}
	for p.tok != scan.BINARY && p.lit != "|" {
		period := p.expect(scan.PERIOD)
		if period != "." {
			break
		}
		args = append(args, p.expect(scan.IDENT))
	}
	p.expect(scan.BINARY)
	return args
}

// statements :=
//     [expression (',' statements)*]
//
func (p *Parser) parseStatements(stmts []rt.Expr) []rt.Expr {
	stmts = append(stmts, p.parseExpr())
	switch p.tok {
	case scan.RBRACE, scan.PERIOD:
		return stmts
	case scan.COMMA:
		p.expect(scan.COMMA)
		stmts = p.parseStatements(stmts)
	default:
		p.error(p.pos, "expected termination to the block or statement. found '%s'", p.lit)
		p.next()
	}
	return stmts
}
