package main

import (
    "log"
    "fmt"
)

type (
    Expr interface {
    }

    Ident struct {
        NamePos int
        Name string
    }

    BasicLit struct {
        ValuePos int
        Kind Token
        Value string
    }

    ParenExpr struct {
        Lparen int
        X Expr
        Rparen int
    }

    UnaryExpr struct {
        OpPos int
        Op Token
        X Expr
    }

    BinaryExpr struct {
        X Expr
        OpPos int
        Op Token
        Y Expr
    }
)

func (x *BasicLit) String() string {
    return fmt.Sprintf("BasicLit(%v %v)", x.Kind, x.Value)
}

func (x *BinaryExpr) String() string {
    return fmt.Sprintf("BinaryExpr(%v %v %v)", x.X, x.Op, x.Y)
}

func (x *ParenExpr) String() string {
    return fmt.Sprintf("ParenExpr(%v)", x.X)
}

type Parser struct {
    scanner *Scanner
    lit string
    pos int
    tok Token
    exprLev int
    inRhs bool
    indent int
    trace bool
}

func NewParser(b []byte) *Parser {
    s := NewScanner(b)
    p := &Parser{
        scanner: s,
        pos: 0,
        tok: s.Scan(),
        lit: s.tt,
        trace: true,
    }
    return p
}

func (p *Parser) printTrace(a ...interface{}) {
    const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
    const n = len(dots)
    pos := p.pos
    fmt.Printf("%5d: ", pos)
    i := 2 * p.indent
    for i > n {
        fmt.Print(dots)
        i -= n
    }
    fmt.Print(dots[0:i])
    fmt.Println(a...)
}

func trace(p *Parser, msg string) *Parser {
    p.printTrace(msg, "(")
    p.indent++
    return p
}

func un(p *Parser) {
    p.indent--
    p.printTrace(")")
}

func (p *Parser) tokPrec() (Token, int) {
    tok := p.tok
    if p.inRhs && tok == ASSIGN {
        tok = EQ
    }
    return tok, tok.Precedence()
}

func (p *Parser) next() {
    p.tok = p.scanner.Scan()
    p.pos = p.scanner.pos
    p.lit = p.scanner.tt
}

func (p *Parser) expect(tok Token) int {
    pos := p.pos
    if p.tok != tok {
        p.errorExpected("'" + tok.String() + "'")
    }
    p.next()
    return pos
}

func (p *Parser) errorExpected(s string) {
    log.Fatalf("expected " + s)
}

func (p *Parser) parseIdent() *Ident {
    pos := p.pos
    name := "_"
    if p.tok == IDENT {
        name = p.lit
        p.next()
    } else {
        p.errorExpected("identifier")
    }
    return &Ident{NamePos: pos, Name: name}
}

func (p *Parser) parsePropRef(lhs Expr) Expr {
    return nil
}

func (p *Parser) parseVerbRef(lhs Expr) Expr {
    return nil
}

func (p *Parser) parsePrimaryExpr(lhs bool) Expr {
    if p.trace {
        defer un(trace(p, "PrimaryExpr"))
    }
    x := p.parseOperand(lhs)
L:
    for {
        switch p.tok {
        case PERIOD:
            x = p.parsePropRef(x)
        case COLON:
            x = p.parseVerbRef(x)
        default:
            break L
        }
    }
    return x
}

func (p *Parser) parseUnaryExpr(lhs bool) Expr {
    if p.trace {
        defer un(trace(p, "UnaryExpr"))
    }
    switch p.tok {
    case ADD, SUB, NOT:
        pos, op := p.pos, p.tok
        p.next()
        x := p.parseUnaryExpr(false)
        return &UnaryExpr{OpPos: pos, Op: op, X: x}
    }
    return p.parsePrimaryExpr(lhs)
}

func (p *Parser) parseRhs() Expr {
    old := p.inRhs
    p.inRhs = true
    x := p.parseExpr(false)
    p.inRhs = old
    return x
}

func (p *Parser) parseBinaryExpr(lhs bool, prec1 int) Expr {
    if p.trace {
        defer un(trace(p, "BinaryExpr"))
    }
    x := p.parseUnaryExpr(lhs)
    for _, prec := p.tokPrec(); prec >= prec1; prec-- {
        for {
            op, oprec := p.tokPrec()
            if oprec != prec {
                break
            }
            pos := p.expect(op)
            if lhs {
                lhs = false
            }
            y := p.parseBinaryExpr(false, prec+1)
            x = &BinaryExpr{X: x, OpPos: pos, Op: op, Y: y}
        }
    }
    return x
}

func (p *Parser) parseExpr(lhs bool) Expr {
    if p.trace {
        defer un(trace(p, "Expression"))
    }
    return p.parseBinaryExpr(lhs, LowestPrec+1)
}

func (p *Parser) parseOperand(lhs bool) Expr {
    if p.trace {
        defer un(trace(p, "Operand"))
    }
    switch p.tok {
    case IDENT:
        return p.parseIdent()
    case INT, FLOAT, STRING, OBJ:
        x := &BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
        p.next()
        return x
    case LPAREN:
        lparen := p.pos
        p.next()
        p.exprLev++
        x := p.parseRhs()
        p.exprLev--
        rparen := p.expect(RPAREN)
        return &ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
    }
    return nil
}

func main() {
    input := "2 - 3 * (2 / 5)"
    p := NewParser([]byte(input))
    var r Expr
    r = p.parseExpr(false)
    log.Printf("%v", r)
}