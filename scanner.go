package main

import (
    "regexp"
    "fmt"
)

type Token int 

const (
    EOF Token = iota
    ILLEGAL
    IDENT
    STR
    INT
    OBJ
    FLOAT
    LBRACE
    RBRACE
    LPAREN
    RPAREN
    LBRACK
    RBRACK
    COLON
    SEMICOLON
    COMMA
    PERIOD
    ADD
    SUB
    MUL
    DIV
    MOD
    ASSIGN
    EQ
    NE
    LT
    GT
    LE
    GE
    NOT
    AND
    OR
    IN
    RANGE
    OPERATOR
    KEYWORD
)

var tokens = map[Token]string {
    EOF: "EOF",
    ILLEGAL: "ILLEGAL",
    IDENT: "IDENT",
    INT: "INT",
    STR: "STR",
    FLOAT: "FLOAT",
    OBJ: "OBJ",
    ADD: "ADD",
    SUB: "SUB",
    MUL: "MUL",
    DIV: "DIV",
    MOD: "MOD",
    ASSIGN: "ASSIGN",
    LPAREN: "LPAREN",
    RPAREN: "RPAREN",
    LBRACK: "LBRACK",
    RBRACK: "RBRACK",
    LBRACE: "LBRACE",
    RBRACE: "RBRACE",
    RANGE: "RANGE",
    COLON: "COLON",
    SEMICOLON: "SEMICOLON",
    PERIOD: "PERIOD",
    COMMA: "COMMA",
}

const (
    LowestPrec = 0
    UnaryPrec = 7
    HighestPrec = 8
)

func (t Token) Precedence() int {
    switch t {
    case ASSIGN:
        return 1
    case OR:
        return 2
    case AND:
        return 3
    case EQ, NE, LT, LE, GT, GE:
        return 4
    case ADD, SUB:
        return 5
    case MUL, DIV, MOD:
        return 6
    }
    return LowestPrec
}

func (t Token) String() string {
    if s, ok := tokens[t]; ok {
        return s
    }
    return fmt.Sprintf("%v", int(t))        
}

var keyword = map[string]bool {
    "for": true,
    "endfor": true,
    "in": true,
    "if": true,
    "else": true,
    "endif": true,
    "return": true,
}

var operator = map[string]Token {
    "+": ADD,
    "-": SUB,
    "*": MUL,
    "/": DIV,
    "%": MOD,
    "..": RANGE,
    "==": EQ,
    "=": ASSIGN,
    "<": LT,
    ">": GT,
    "<=": LE,
    ">=": GE,
    "!=": NE,
    "!": NOT,
    "&&": AND,
    "||": OR,
    "in": IN,
}

var symbol = map[string]Token {
    "{": LBRACE,
    "}": RBRACE,
    "[": LBRACK,
    "]": RBRACK,
    "(": LPAREN,
    ")": RPAREN,
    ":": COLON,
    ";": SEMICOLON,
    ".": PERIOD,
    ",": COMMA,
}

type Scanner struct {
    b []byte
    pos int
    tt string
    advance bool
}

func NewScanner(b []byte) *Scanner {
    return &Scanner{b, 0, "", false}
}

func (s *Scanner) Peek() int {
    if s.pos < len(s.b) {
        return int(s.b[s.pos])
    }
    return int(EOF)
}

var reInt = regexp.MustCompile(`^[0-9]+`)
var reFloat = regexp.MustCompile(`^[0-9]+\.[0-9]+`)
var reStr = regexp.MustCompile(`^".*"`)
var reObj = regexp.MustCompile(`^#[0-9\-]+`)
var reIdent = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*`)
var reOp = regexp.MustCompile(`^[\^\.\|\?\*\+\-!%<>=&/]+`)
var reSym = regexp.MustCompile(`^[\{\[\(\}\]\)\$\.@,;:]`)
var reWs = regexp.MustCompile(`^[\n\t\r ]+`)

func (s *Scanner) scanRegexp(re *regexp.Regexp) (ok bool) {
    t := s.b[s.pos:]
    if loc := re.FindIndex(t); len(loc) == 2 {
        ok = true
        s.tt = string(t[loc[0]:loc[1]])
    }
    return
}

func (s *Scanner) scanInt() (ok bool) {
    return s.scanRegexp(reInt)
}

func (s *Scanner) scanFloat() (ok bool) {
    return s.scanRegexp(reFloat)
}

func (s *Scanner) scanStr() (ok bool) {
    return s.scanRegexp(reStr)
}

func (s *Scanner) scanObj() (ok bool) {
    return s.scanRegexp(reObj)
}

func (s *Scanner) scanWs() (ok bool) {
    return s.scanRegexp(reWs)
}

func (s *Scanner) scanIdent() (ok bool) {
    if ok = s.scanRegexp(reIdent); ok {
        ok = !keyword[s.tt]
    }
    return
}

func (s *Scanner) scanKeyword() (ok bool) {
    if ok = s.scanIdent(); !ok {
        s.advance = keyword[s.tt]
    }
    return
}

func (s *Scanner) scanOp() (tok Token) {
    tok = ILLEGAL
    if ok := s.scanRegexp(reOp); ok {
        tok, ok = operator[s.tt]
        if !ok {
            tok = ILLEGAL
        }
        s.advance = ok
    }
    return
}

func (s *Scanner) scanLit() (tok Token) {
    tok = ILLEGAL
    switch {
    case s.scanStr():
        s.advance = true
        tok = STR
        break
    case s.scanIdent():
        s.advance = true
        tok = IDENT
        break
    case s.scanFloat():
        s.advance = true
        tok = FLOAT
        break
    case s.scanInt():
        s.advance = true
        tok = INT
        break
    case s.scanObj():
        s.advance = true
        tok = OBJ
        break
    }
    return
}

func (s *Scanner) scanSym() (tok Token) {
    tok = ILLEGAL
    if ok := s.scanRegexp(reSym); ok {
        tok, s.advance = symbol[s.tt]
    }
    return
}

func (s *Scanner) Seek(pos int) {
    s.pos = pos
}

func (s *Scanner) Next() int {
    if s.pos < len(s.b) {
        r := s.b[s.pos]
        s.pos += 1
        return int(r)        
    }
    return int(EOF)
}

func (s *Scanner) tryAdvance() (tok Token) {
    tok = ILLEGAL
    if s.advance {
        s.pos += len(s.tt)
        s.advance = false
    }    
    if s.scanWs() {
        s.pos += len(s.tt)
    }
    if s.pos >= len(s.b) {
        tok = EOF
    }
    return
}

func (s *Scanner) Scan() (tok Token) {
    tok = s.tryAdvance()
    if tok = s.scanLit(); tok != ILLEGAL {
        return
    }
    if tok = s.scanOp(); tok != ILLEGAL {
        return
    }
    if s.scanKeyword() {
        tok = KEYWORD
        return
    }
    tok = s.scanSym()
    return
}