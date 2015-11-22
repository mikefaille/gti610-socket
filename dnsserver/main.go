/* UDPDaytimeServer
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mikefaille/gti610-socket/dnsserver/dnsUtil"
)

func main() {

	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		var byt [512]byte

		bytSize, addr, err := conn.ReadFromUDP(byt[0:])

		if err != nil {
			return
		}

		go handleClient(conn, addr, &byt, bytSize)
	}
}

func handleClient(conn *net.UDPConn, srcAddr *net.UDPAddr, byt *[512]byte, bytSise int) {

	packet := new(dnsUtil.QuestionPckt)

	packet.EncodeBytes(byt)
	// //
	// // fmt.Fprint("", byt[0:16])
	// daytime := time.Now().String()

	// buf := bytes.NewReader(byt[0:16])
	// err := binary.Read(buf, binary.LittleEndian, &packet.ID)
	// if err != nil {
	// 	fmt.Println("binary.Read failed:", err)
	// }
	fmt.Println(packet.ID)

	// conn.WriteToUDP(packet.ID, srcAddr)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
