package main

import (
    "log"
)

func main() {
    input := "fubar"
    s := NewScanner([]byte(input))
    log.Println(string(s.Peek()))
}