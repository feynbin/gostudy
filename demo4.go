package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, 11)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    fmt.Println(string(b))
}
