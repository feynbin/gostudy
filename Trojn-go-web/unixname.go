package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// 生成随机字符串
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 11)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	fmt.Println(string(b))
}
