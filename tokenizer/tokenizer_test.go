package main

import (
    "log"
    "fmt"
    "testing"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
)

func TestTokenizer(t *testing.T) {

    // 1. test the top-level index handler
    ts := httptest.NewServer(http.HandlerFunc(index))
    defer ts.Close()

    res, err := http.Get(ts.URL)
    if err != nil {
        log.Fatal(err)
    }

    greeting, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s", greeting)

    // 2. continue verifying authorized tokens
    // ...

}
