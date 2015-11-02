package main

import "net"
import "fmt"

func main() {

	netListen, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := netListen.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}

func handleConnection(con net.Conn) {

	fmt.Println("out")
}
