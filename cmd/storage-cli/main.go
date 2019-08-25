package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	addr = flag.String("addr", "localhost:2094", "Addr")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Command: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Result: " + message)
	}
}
