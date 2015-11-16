package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}
	host := os.Args[1]
	conn, err := net.Dial("tcp", host+":1201")
	checkError(err)
	_, err = conn.Write([]byte("HEAD"))
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		//		fmt.Println(err)
		line = strings.TrimRight(line, " \t\r\n")
		_, err = conn.Write([]byte(line))

		if err != nil {
			conn.Close()
			break

		}
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
