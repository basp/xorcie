package main

import (
    "log"
    "regexp"
)

type Token int 

const (
    EOF Token = -(iota + 1)
    ILLEGAL
    IDENT
    STRING
    INT
    OBJ
    FLOAT
    OPERATOR
    LBRACK
    RBRACK
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

type Scanner struct {
    b []byte
    pos int
    tt string
}

func NewScanner(b []byte) *Scanner {
    return &Scanner{b, 0, ""}
}

var reInt = regexp.MustCompile(`^[0-9]+`)
var reFloat = regexp.MustCompile(`^[0-9]+\.[0-9]+`)
var reStr = regexp.MustCompile(`^".*"`)
var reObj = regexp.MustCompile(`^#[0-9\-]+`)
var reIdent = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*`)
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

func (s *Scanner) ScanLiteral() (tok Token) {
    tok = ILLEGAL
    switch {
    case s.scanStr():
        tok = STRING
    case s.scanIdent():
        tok = IDENT
    case s.scanFloat():
        tok = FLOAT
    case s.scanInt():
        tok = INT
    case s.scanObj():
        tok = OBJ
    }
    return
}

func main() {
    b := []byte(`#-123endfor"quux"`)
    s := NewScanner(b)
    var ok bool
    if ok = s.scanObj(); ok {
        log.Println(s.tt)
    }
    if ok = s.scanKeyword(); ok {
        log.Println(s.tt)
    }
    if ok = s.scanWs(); ok {
        log.Printf("'%v'", s.tt)
    }
    if ok = s.scanStr(); ok {
        log.Println(s.tt)
    }
}