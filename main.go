package main

import (
    "log"
)

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