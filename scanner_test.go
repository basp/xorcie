package main

import (
    "testing"
)

func TestNext(t *testing.T) {
    s := NewScanner([]byte("foo"))
    var r int
    r = s.Next()
    if r != 'f' {
        t.Error("Expected f but got", r)
    }
    r = s.Next()
    if r != 'o' {
        t.Error("Expected o but got", r)
    }
    r = s.Next()
    if r != 'o' {
        t.Error("Expected o but got", r)
    }
    r = s.Next()
    if r != int(EOF) {
        t.Error("Expected EOF but got", r)
    }
}

type simpleScanCase struct {
    input string
    expected Token
}

var simpleScanCases = []simpleScanCase {
    simpleScanCase{"foo", IDENT},
    simpleScanCase{"123.45", FLOAT},
    simpleScanCase{`"foo"`, STRING},
    simpleScanCase{"123", INT},
    simpleScanCase{"#123", OBJ},
}

func TestSimpleScan(t *testing.T) {
    for _, c := range simpleScanCases {
        s := NewScanner([]byte(c.input))
        actual := s.Scan()
        if actual != c.expected {
            t.Error("Expected", c.expected, "but got", actual)
        }
    }
}

func TestScanVarious(t *testing.T) {
    input := "foo * 123.45"
    s := NewScanner([]byte(input))
    var tok Token
    tok = s.Scan()
    if tok != IDENT {
        t.Error("Expected IDENT but got", tok)
    }
    tok = s.Scan()
    if tok != MUL {
        t.Error("Expected MUL but got", tok)
    }
    tok = s.Scan()
    if tok != FLOAT {
        t.Error("Expected FLOAT but got", tok)
    }
}

var symScanCases = []simpleScanCase {
    simpleScanCase{"{", LBRACE},
    simpleScanCase{"]", RBRACK},
    simpleScanCase{"(", LPAREN},
    simpleScanCase{";", SEMICOLON},
    simpleScanCase{",", COMMA},
    simpleScanCase{".", PERIOD},
}

func TestScanSymbol(t *testing.T) {
    for _, c := range symScanCases {
        s := NewScanner([]byte(c.input))
        actual := s.Scan()
        if actual != c.expected {
            t.Error("Expected", c.expected, "but got", actual)
        }   
    }
}