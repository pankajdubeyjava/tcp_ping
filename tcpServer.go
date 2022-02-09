package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host example: 127.0.0.1:1234")
		return
	}

	PORT := arguments[1] 
        fmt.Println("server is running on: ",PORT)
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c)

	for {

	}
}

