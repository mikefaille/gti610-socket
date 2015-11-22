package dnsUtil

import (
	"fmt"
	"math"
	"net"
)

type DomainName string

type packet struct {
	Ip net.IP
	header
}

type header struct {
	ID      uint16
	QR      uint16
	OPCODE  uint16
	AA      uint16
	TC      uint16
	RD      uint16
	RA      uint16
	Z       uint16
	RCODE   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type ReponsePckt struct {
	packet
	NAME     DomainName
	TYPE     uint16
	CLASS    uint16
	TTL      uint16
	RDLENGTH uint16
	RDATA    string
}

type QuestionPckt struct {
	packet
	QNAME  DomainName
	QTYPE  uint16
	QCLASS uint16
}

func (qPckt QuestionPckt) EncodeBytes(b *[512]byte) {
	qPckt.parseHeader(b)

}

type GenericPacket interface {
	EncodeBytes(b []byte)
}

func (header) parseHeader(b *[512]byte) (h header) {
	fmt.Println("Binary : ", b[0:64])
	parseBytes(b, 0, 16, h.ID)
	parseBytes(b, 15, 1, h.QR)
	parseBytes(b, 16, 4, h.OPCODE)
	parseBytes(b, 20, 1, h.AA)
	parseBytes(b, 21, 1, h.TC)
	parseBytes(b, 22, 1, h.RD)
	parseBytes(b, 23, 1, h.RA)
	parseBytes(b, 24, 4, h.Z)
	parseBytes(b, 28, 4, h.RCODE)
	parseBytes(b, 32, 16, h.QDCOUNT)
	parseBytes(b, 48, 16, h.ANCOUNT)
	parseBytes(b, 64, 16, h.NSCOUNT)
	parseBytes(b, 80, 16, h.ARCOUNT)

	return
}

func parseBytes(b *[512]byte, offset int, size int, data interface{}) {

	byteOffset := uint((offset + 1) / 8)

	fmt.Println("offset : ", byteOffset) // int(math.Ceil(float64((offset+1)/8))))
	fmt.Println("size : ", size)
	byteSize := uint(math.Ceil(float64(size) / 8.0))
	fmt.Println("byteSize : ", byteSize)

	switch data.(type) {
	case uint16:
		byteToAnalysis := b[byteOffset : byteOffset+byteSize]
		fmt.Println("byteToAnalysis : ", byteToAnalysis)

		bitMask := buildBitMask(uint(offset), uint(size))
		fmt.Println("bitmask : ", bitMask)
		data := byteToInt(byteToAnalysis, byteSize)
		fmt.Println("data : ", data&bitMask)
		fmt.Println("")
		break
	case DomainName:

		// TODO À TESTER !
		data := data.(string)
		isZeroParsed := false
		strLenght := uint(b[byteOffset])

		for isZeroParsed == true {

			data += string(b[byteOffset : uint(byteOffset)+strLenght])

			if uint(b[byteOffset]) == uint(0) {
				isZeroParsed = false
				break
			}
			byteOffset = uint(byteOffset) + strLenght - 1
		}

		fmt.Println(data)

		break
	default:
		fmt.Println("unparsed data :/")
		fmt.Println(data)
		fmt.Println("")
		break
	}

}

//////// 0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
/////// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|                      ID                       |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|                    QDCOUNT                    |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|                    ANCOUNT                    |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|                    NSCOUNT                    |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// 	|                    ARCOUNT                    |
// 	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

func buildBitMask(offset uint, size uint) (mask uint16) {

	for i := offset; i < size-1; i++ {
		mask = mask | (1 << uint(i))

	}
	return mask
}

func byteToInt(b []byte, size uint) (data uint16) {

	for i := size; i > 0; i-- {
		data = data | (uint16(b[i-1]) << uint(8*(size-i)))

	}
	return
}
