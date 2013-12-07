package main

import (
    "log"
)

func parse(s string) Expr {
    p := NewParser([]byte(s))
    return p.parseIfStmt()
}

func main() {
    // expr := parse("v = {foo:bar(quux) * 5, #123}")
    // expr := parse("foo = {foo:quux()}")
    // expr := parse("foo[0..5][bar];")
    // s := parse("foo = (3 * (2 * \"bar\")); return 5;")
    s := parse("if (foo == bar) quux; elseif (quux < 3) zoz; {1,2,3}; elseif (nox > #123) fuz; else foz; endif")
    log.Printf("%v", s)
}