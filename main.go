package main

import (
    "log"
)

func parse(s string) Expr {
    p := NewParser([]byte(s))
    return p.parseExpr(false)
}

func main() {
    expr := parse("v = {foo:bar(quux) * 5, #123}")
    // expr := parse("foo = {foo:quux()}")
    log.Printf("%v", expr)
}