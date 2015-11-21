/* UDPDaytimeServer
 */
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		var buf [512]byte
		_, addr, err := conn.ReadFromUDP(buf[0:])

		if err != nil {
			return
		}
		go handleClient(conn, addr)
	}
}

func handleClient(conn *net.UDPConn, srcAddr *net.UDPAddr) {

	daytime := time.Now().String()

	conn.WriteToUDP([]byte(daytime), srcAddr)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

// $ 6g echo.go && 6l -o echo echo.6
// $ ./echo
//
//  ~ in another terminal ~
//
// $ nc localhost 3540

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strconv"
// )

// const PORT = 1202

// func main() {
// 	server, err := net.Listen("udp", ":"+strconv.Itoa(PORT))
// 	if server == nil {
// 		log.Panic("couldn't start listening: " + err.Error())
// 	}
// 	conns := clientConns(server)
// 	for {
// 		handleConn(<-conns)
// 	}
// }

// func clientConns() chan net.Conn {
// 	ch := make(chan net.Conn)
// 	i := 0
// 	go func() {
// 		for {
// 			client, err := listener.Accept()
// 			if client == nil {
// 				fmt.Printf("couldn't accept: " + err.Error())
// 				continue
// 			}
// 			i++
// 			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
// 			ch <- client
// 		}
// 	}()
// 	return ch
// }

// func handleConn(client net.Conn) {
// 	b := bufio.NewReader(client)
// 	for {
// 		line, err := b.ReadBytes('\n')
// 		if err != nil { // EOF, or worse
// 			break
// 		}
// 		client.Write(line)
// 	}
// }
