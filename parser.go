package main

import (
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

    ListExpr struct {
        Lbrace int
        Xs []Expr
        Rbrace int
    }

    PropExpr struct {
        Pos int
        Obj Expr
        Prop Expr
    }

    VerbExpr struct {
        Pos int
        Obj Expr
        Verb Expr 
        Args []Expr
    }

    IndexExpr struct {
        X Expr
        Lbrack int
        Index Expr
        Rbrack int
    }

    SliceExpr struct {
        X Expr
        Lbrack int
        Low Expr
        High Expr
        Rbrack int
    }
)

func (x *Ident) String() string {
    return fmt.Sprintf("Ident('%v')", x.Name)
}

func (x *BasicLit) String() string {
    return fmt.Sprintf("BasicLit(%v %v)", x.Kind, x.Value)
}

func (x *BinaryExpr) String() string {
    return fmt.Sprintf("BinaryExpr(%v %v %v)", x.X, x.Op, x.Y)
}

func (x *ParenExpr) String() string {
    return fmt.Sprintf("ParenExpr(%v)", x.X)
}

func (x *UnaryExpr) String() string {
    return fmt.Sprintf("UnaryExpr(%v %v)", x.Op, x.X)
}

func (x *ListExpr) String() string {
    return fmt.Sprintf("ListExpr(%v)", x.Xs)
}

func (x *PropExpr) String() string {
    return fmt.Sprintf("PropExpr(%v %v)", x.Obj, x.Prop)
}

func (x *VerbExpr) String() string {
    return fmt.Sprintf("VerbExpr(%v %v %v)", x.Obj, x.Verb, x.Args)
}

func (x *IndexExpr) String() string {
    return fmt.Sprintf("IndexExpr(%v %v)", x.X, x.Index)
}

func (x *SliceExpr) String() string {
    return fmt.Sprintf("SliceExpr(%v %v %v)", x.X, x.Low, x.High)
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
        p.errorExpected("'" + tok.String() + "' but got '" + p.tok.String() + "'")
    }
    p.next()
    return pos
}

func (p *Parser) errorExpected(s string) {
    panic("expected " + s)
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

func (p *Parser) parsePropExpr(lhs bool, obj Expr) Expr {
    if p.trace {
        defer un(trace(p, "PropExpr"))
    }
    pos := p.expect(PERIOD)
    var prop Expr
    if p.tok == LPAREN {
        p.next()
        prop = p.parseExpr(lhs)
        p.expect(RPAREN)
    } else {
        prop = p.parseIdent()
    }
    return &PropExpr{pos, obj, prop}
}

func (p *Parser) parseVerbExpr(lhs bool, obj Expr) Expr {
    if p.trace {
        defer un(trace(p, "VerbExpr"))
    }
    pos := p.expect(COLON)
    var y Expr
    if p.tok == LPAREN {
        p.next()
        y = p.parseExpr(lhs)
        p.expect(RPAREN)
    } else {
        y = p.parseIdent()        
    }
    pos = p.expect(LPAREN)
    args := p.parseListExpr(false, RPAREN)
    return &VerbExpr{pos, obj, y, args}
}

func (p *Parser) parseIndexOrSlice(lhs bool, x Expr) Expr {
    if p.trace {
        defer un(trace(p, "IndexOrSlice"))
    }
    lbrack := p.expect(LBRACK)
    p.exprLev++
    var index [2]Expr
    index[0] = p.parseRhs()
    if p.tok == RANGE {
        p.next()
        index[1] = p.parseRhs()
    }
    rbrack := p.expect(RBRACK)
    if index[1] != nil {
        return &SliceExpr{x, lbrack, index[0], index[1], rbrack}
    } else {
        return &IndexExpr{x, lbrack, index[0], rbrack}
    }
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
            x = p.parsePropExpr(lhs, x)
        case COLON:
            x = p.parseVerbExpr(lhs, x)
            p.expect(RPAREN)
        case LBRACK:
            x = p.parseIndexOrSlice(lhs, x)
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

func (p *Parser) parseListExpr(lhs bool, final Token) (xs []Expr) {
    if p.trace {
        defer un(trace(p, "ListExpr"))
    }
    if p.tok == final {
        return
    }
    xs = make([]Expr, 0, 16)
    for {
        xs = append(xs, p.parseExpr(lhs))        
        if p.tok != COMMA {
            return
        }
        p.next()
    }
    return
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
    case INT, FLOAT, STR, OBJ:
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
    case LBRACE:
        lbrace := p.pos
        p.next()
        xs := p.parseListExpr(lhs, RBRACE)
        rbrace := p.expect(RBRACE)
        return &ListExpr{Lbrace: lbrace, Xs: xs, Rbrace: rbrace}
    }
    return nil
}