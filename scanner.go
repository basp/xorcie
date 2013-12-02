package main

import (
    "regexp"
)

type Token int 

const (
    EOF Token = -iota
    ILLEGAL
    IDENT
    STRING
    INT
    OBJ
    FLOAT
    LPAREN
    RPAREN
    LBRACK
    RBRACK
    OPERATOR
    KEYWORD
)

var keyword = map[string]bool {
    "for": true,
    "endfor": true,
    "in": true,
    "if": true,
    "else": true,
    "endif": true,
    "return": true,
}

var operator = map[string]bool {
    "+": true,
    "-": true,
    "*": true,
    "/": true,
    "%": true,
    "..": true,
    "==": true,
    "=": true,
    "<": true,
    ">": true,
    // etc...
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
var reOp = regexp.MustCompile(`^[\$\^\.\|\?\*\+!%<>=]`)
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

func (s *Scanner) scanOp() (ok bool) {
    return s.scanRegexp(reOp)
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

func (s *Scanner) scanLit() (tok Token) {
    tok = ILLEGAL
    switch {
    case s.scanOp():
        tok = OPERATOR
        break
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

func (s *Scanner) Next() int {
    if s.pos < len(s.b) {
        r := s.b[s.pos]
        s.pos += 1
        return int(r)        
    }
    return int(EOF)
}

func (s *Scanner) Scan() (tok Token) {
    tok = s.scanLit()
    // We just bluntly run by any 
    // whitespace we encounter
    s.scanWs()
    return
}