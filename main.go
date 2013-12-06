package main

import (
    "log"
)

func parse(s string) Expr {
    p := NewParser([]byte(s))
    return p.parseStmtList()
}

func main() {
    // expr := parse("v = {foo:bar(quux) * 5, #123}")
    // expr := parse("foo = {foo:quux()}")
    // expr := parse("foo[0..5][bar];")
    s := parse("foo = (3 * (2 * \"bar\")); return 5;")
    log.Printf("%v", s)
}