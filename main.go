package main

import (
    "log"
)

func parse(s string) Expr {
    p := NewParser([]byte(s))
    return p.parseExpr(false)
}

func main() {
    expr := parse("foo:quux(x,y,{1,2,3}) = bar.(baz + 123) * 12.5)")
    log.Printf("%v", expr)
}