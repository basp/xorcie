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

    RangeExpr struct {
        Lbrack int
        Low Expr
        High Expr
        Rbrack int
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

    AssignStmt struct {
        Lhs Expr
        TokPos int
        Tok Token
        Rhs Expr
    }

    Stmt interface {        
    }

    ExprStmt struct {
        X Expr
    }    

    EmptyStmt struct {
        Semicolon int
    }

    IfStmt struct {
        If int
        Cond Expr
        List []Stmt
        Else Stmt
    }

    ReturnStmt struct {
        Return int
        Result Expr
    }

    ForRangeStmt struct {
        For int
        Var *Ident
        In int
        Range Expr
        List []Stmt
        EndFor int
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

func (x *RangeExpr) String() string {
    return fmt.Sprintf("Range(%v %v)", x.Low, x.High)
}

func (x *IndexExpr) String() string {
    return fmt.Sprintf("IndexExpr(%v %v)", x.X, x.Index)
}

func (x *SliceExpr) String() string {
    return fmt.Sprintf("SliceExpr(%v %v %v)", x.X, x.Low, x.High)
}

func (x *EmptyStmt) String() string {
    return fmt.Sprintf("EmptyStmt()")
}

func (x *ExprStmt) String() string {
    return fmt.Sprintf("ExprStmt(%v)", x.X)
}

func (x *ReturnStmt) String() string {
    return fmt.Sprintf("ReturnStmt(%v)", x.Result)
}

func (x *IfStmt) String() string {
    return fmt.Sprintf("IfStmt(%v %v %v)", x.Cond, x.List, x.Else)    
}

func (x *ForRangeStmt) String() string {
    return fmt.Sprintf("ForRange(%v %v %v)", x.Var, x.Range, x.List) 
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
        trace: false,
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

func (p *Parser) parseSimpleStmt() Stmt {
    if p.trace {
        defer un(trace(p, "SimpleStmt"))
    }
    if p.tok == SEMICOLON {
        return &EmptyStmt{Semicolon: p.pos}
    }
    x := p.parseExpr(false)
    p.expect(SEMICOLON)
    return &ExprStmt{X: x}
}

func (p *Parser) parseReturnStmt() Stmt {
    if p.trace {
        defer un(trace(p, "ReturnStmt"))
    }
    pos := p.pos
    p.expect(RETURN)
    x := p.parseExpr(false)
    p.expect(SEMICOLON)
    return &ReturnStmt{Return: pos, Result: x}
}

func (p *Parser) parseIfStmt() Stmt {
    if p.trace {
        defer un(trace(p, "IfStmt"))
    }
    if p.tok == IF || p.tok == ELSEIF {
        var _if int
        if p.tok == IF || p.tok == ELSEIF {
            _if = p.expect(p.tok)
        }
        c := p.parseExpr(true)
        l := p.parseStmtList()
        if p.tok == ENDIF {
            return &IfStmt{_if, c, l, nil}
        }
        var _else Stmt
        if p.tok == ELSE {
            p.expect(ELSE)
            _else = p.parseStmtList()
        }
        if p.tok == ELSEIF {
            _else = p.parseIfStmt()
        }
        return &IfStmt{_if, c, l, _else}
    }
    return nil
}

func (p *Parser) parseForRangeStmt() Stmt {
    if p.trace {
        defer un(trace(p, "ForRange"))
    }   
    _for := p.expect(FOR)
    id := p.parseIdent()
    in := p.expect(IN)
    r := p.parseRange()        
    body := p.parseStmtList()
    endfor := p.expect(ENDFOR)
    return &ForRangeStmt{_for, id, in, r, body, endfor}
}

func (p *Parser) parseStmtList() (list []Stmt) {
    if p.trace {
        defer un(trace(p, "StatementList"))
    }
    for {
        switch p.tok {
        case ENDIF, ENDFOR, ELSE, ELSEIF, EOF:
            return
        }
        list = append(list, p.parseStmt())
    }
    return
}

func (p *Parser) parseStmt() (s Stmt) {
    if p.trace {
        defer un(trace(p, "Statement"))
    }
    switch p.tok {
    case IDENT, INT, FLOAT, STR, LPAREN, LBRACE, ADD, SUB, NOT:
        s = p.parseSimpleStmt()
    case RETURN:
        s = p.parseReturnStmt()
    case FOR:
        s = p.parseForRangeStmt()
    case IF:
        s = p.parseIfStmt()
    }
    return
}

func (p *Parser) parseIdent() *Ident {
    if p.trace {
        defer un(trace(p, "Ident"))
    }
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

func (p *Parser) parseRange() Expr {
    if p.trace {
        defer un(trace(p, "Range"))
    }
    lbrack := p.expect(LBRACK)
    var index [2]Expr
    index[0] = p.parseRhs()
    p.expect(RANGE)
    index[1] = p.parseRhs()
    rbrack := p.expect(RBRACK)
    return &RangeExpr{lbrack, index[0], index[1], rbrack}
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
            if x != nil {
                x = p.parseIndexOrSlice(lhs, x)                
            } else {
                x = p.parseRange()
            }
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