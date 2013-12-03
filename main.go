package main

import (
    "log"
)

type (
    Expr interface {

    }

    
)

func main() {
    input := "foo[1..3] % (3 * \"foo\")"
    s := NewScanner([]byte(input))
    tok := s.Scan()
    for tok != EOF {
        log.Printf("[%v] %v", s.tt, tok)
        tok = s.Scan()
    }
}