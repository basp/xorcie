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
    STRING
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
    STRING: "STRING",
    MUL: "MUL",
    DIV: "DIV",
    MOD: "MOD",
    LPAREN: "LPAREN",
    RPAREN: "RPAREN",
    LBRACK: "LBRACK",
    RBRACK: "RBRACK",
    RANGE: "RANGE",
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
}

func NewScanner(b []byte) *Scanner {
    return &Scanner{b, 0, ""}
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
var reOp = regexp.MustCompile(`^[\^\.\|\?\*\+!%<>=&]+`)
var reSym = regexp.MustCompile(`^[\{\[\(\}\]\)\$\.@,;:]`)
var reWs = regexp.MustCompile(`^[\n\t\r ]+`)

func (s *Scanner) scanRegexp(re *regexp.Regexp) (ok bool) {
    t := s.b[s.pos:]
    if loc := re.FindIndex(t); len(loc) == 2 {
        ok = true
        s.tt = string(t[loc[0]:loc[1]])
        s.pos += len(s.tt)
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
        if !ok {
            s.pos -= len(s.tt)
        }
    }
    return
}

func (s *Scanner) scanKeyword() (ok bool) {
    if ok = s.scanIdent(); !ok {
        ok = keyword[s.tt]
        if ok {
            s.pos += len(s.tt)
        }
    }
    return
}

func (s *Scanner) scanOp() (tok Token) {
    tok = ILLEGAL
    if ok := s.scanRegexp(reOp); ok {
        if tok, ok = operator[s.tt]; !ok {
            tok = ILLEGAL
            s.pos -= len(s.tt)
        }
    }
    return
}

func (s *Scanner) scanLit() (tok Token) {
    tok = ILLEGAL
    switch {
    case s.scanStr():
        tok = STRING
        break
    case s.scanIdent():
        tok = IDENT
        break
    case s.scanFloat():
        tok = FLOAT
        break
    case s.scanInt():
        tok = INT
        break
    case s.scanObj():
        tok = OBJ
        break
    }
    return
}

func (s *Scanner) scanSym() (tok Token) {
    tok = ILLEGAL
    if ok := s.scanRegexp(reSym); ok {
        if tok, ok = symbol[s.tt]; ok {
            return tok
        }
        s.pos -= len(s.tt)
    }
    return
}

func (s *Scanner) Next() int {
    if s.pos < len(s.b) {
        r := s.b[s.pos]
        s.pos += 1
        return int(r)        
    }
    return int(EOF)
}

func (s *Scanner) Scan() (tok Token) {
    s.scanWs()
    if s.pos >= len(s.b) {
        return EOF
    }
    tok = ILLEGAL
    if tok = s.scanLit(); tok != ILLEGAL {
        return
    }
    if tok = s.scanOp(); tok != ILLEGAL {
        return
    }
    tok = s.scanSym()
    return
}