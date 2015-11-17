package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	service := ":1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleClient(conn)
		conn.Close() // we're finished
	}
}

func handleClient(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		line := string(buf[0:n])

		//*******LOGIC*************************************************************************************
		m := map[string]string{"google.com": "8.8.8.8", "cedille.etsmtl.ca": "142.137.247.120"}

		for key, value := range m {
			// for each pair in the map, print key and value
			if line == key {
				fmt.Println(value)
			} else {
				fmt.Println("not mapped")
			}
		}
		//*******END OF LOGIC******************************************************************************
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
