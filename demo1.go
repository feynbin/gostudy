package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Hello World")
	fmt.Println("Hello World")
	end := time.Now()
	fmt.Println(end.Sub(start))
	current_time := time.Now()
	fmt.Println(current_time)
	fmt.Println(current_time.Format(time.DateTime))
	fmt.Println(time.DateTime)

}
