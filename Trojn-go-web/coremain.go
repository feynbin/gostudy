package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入命令：help/list/add/get/set/del")
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		help()
	case "list":
		list()
	case "add":
		add()
	case "get":
		get()
	case "set":
		set()
	case "del":
		del()
	default:
		fmt.Println("无效的命令，请输入 help/list/add/get/set/del")
	}
}

func help() {
	fmt.Println("这是一个命令行工具，支持以下命令：")
	fmt.Println("help\t\t打印帮助信息")
	fmt.Println("list\t\t输出list")
	fmt.Println("add\t\t输出add")
	fmt.Println("get\t\t输出get")
	fmt.Println("set\t\t输出set")
	fmt.Println("del\t\t输出del")
}

func list() {
	fmt.Println("这是 list 命令")
}

func add() {
	fmt.Println("这是 add 命令")
}

func get() {
	fmt.Println("这是 get 命令")
}

func set() {
	fmt.Println("这是 set 命令")
}

func del() {
	fmt.Println("这是 del 命令")
}
