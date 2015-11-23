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
	qPckt.parseQuestion(b)
}

type GenericPacket interface {
	EncodeBytes(b []byte)
}

func (h header) parseHeader(b *[512]byte) {
	// fmt.Println("Binary : ", b[0:64])
	parseBytes(b, 0, 16, &h.ID)
	// fmt.Println("Header id : ", h.ID)

	parseBytes(b, 15, 1, &h.QR)
	// fmt.Println("Header QR : ", h.QR)

	parseBytes(b, 16, 4, &h.OPCODE)
	// fmt.Println("Header AA ")
	parseBytes(b, 20, 1, &h.AA)
	// fmt.Println("Header TC ")
	parseBytes(b, 21, 1, &h.TC)
	// fmt.Println("Header RD ")
	parseBytes(b, 22, 1, &h.RD)
	// fmt.Println("Header RA ")
	parseBytes(b, 23, 1, &h.RA)
	parseBytes(b, 24, 4, &h.Z)
	parseBytes(b, 28, 4, &h.RCODE)
	parseBytes(b, 32, 16, &h.QDCOUNT)
	parseBytes(b, 48, 16, &h.ANCOUNT)
	parseBytes(b, 64, 16, &h.NSCOUNT)
	parseBytes(b, 80, 16, &h.ARCOUNT)

}

func (q QuestionPckt) parseQuestion(b *[512]byte) {

	qNameLen := parseBytesName(b, 96, &q.QNAME)
	fmt.Println("QNAME : ", q.QNAME)
	parseBytes(b, 96+qNameLen, 16, &q.QTYPE)
	fmt.Println("QTYPE : ", q.QTYPE)
	parseBytes(b, 96+qNameLen+16, 16, &q.QCLASS)
	fmt.Println("QCLASS : ", q.QCLASS)

	return
}

func (r ReponsePckt) parseReponse(b *[512]byte) {

	rNameLen := parseBytesName(b, 96, &r.NAME)
	parseBytes(b, 96+rNameLen, 16, &r.TYPE)
	parseBytes(b, 96+rNameLen+16, 16, &r.CLASS)
	parseBytes(b, 96+rNameLen+16+16, 32, &r.TTL)
	parseBytes(b, 96+rNameLen+16+16+32, 16, &r.RDLENGTH)
	parseBytesRdata(b, 96+rNameLen+16+16+32+16, int(r.RDLENGTH), r.RDATA)

	return
}

func parseBytes(b *[512]byte, offset int, size int, data *uint16) {

	byteOffset := uint16((offset + 1) / 8)
	// fmt.Println("offset : ", offset)         // int(math.Ceil(float64((offset+1)/8))))
	// fmt.Println("byteoffset : ", byteOffset) // int(math.Ceil(float64((offset+1)/8))))
	// fmt.Println("size : ", size)
	byteSize := uint16(math.Ceil(float64(size) / 8.0))
	// fmt.Println("byteSize : ", byteSize)

	byteToAnalysis := b[byteOffset : byteOffset+byteSize]
	// fmt.Println("byteToAnalysis : ", byteToAnalysis)

	// Special offset for byte offset % 8

	*data = bitFilterToInt(byteToAnalysis, uint16(offset%8+1), uint16(size), byteSize)
	// fmt.Println("data : ", data)
	// fmt.Println("")
}

func parseBytesName(b *[512]byte, offset uint16, data *DomainName) int {

	byteOffset := uint16((offset + 1) / 8)

	fmt.Println("offset : ", byteOffset) // int(math.Ceil(float64((offset+1)/8))))

	isZeroParsed := false
	strIndex := byteOffset
	strLenght := uint16(b[strIndex])
	strIndex++

	for isZeroParsed == false {

		*data += DomainName(b[strIndex : strIndex+strLenght])

		// mise à jour du prochain index
		strIndex = uint16(strIndex) + uint16(strLenght)
		strLenght = uint16(b[strIndex])

		// strLenght++ // il y a un point
		fmt.Println("domain name : ", *data)
		fmt.Println("strIndex", strIndex)
		fmt.Println("Name lenght : ", strLenght)
		if uint(b[strIndex]) == uint(0) {
			isZeroParsed = true
			break
		} else {
			*data = *data + "."
		}
		strIndex++

	}
	// Note :  0 is counted
	fmt.Println("Name lenght ", strIndex-byteOffset)

	// Note :  0 is counted
	return int(strIndex - byteOffset)
}

func parseBytesRdata(b *[512]byte, offset int, size int, data string) (nameLenght int) {
	fmt.Println("Rdata need to be parsed")
	return
}

//   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
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
func buildBitMask(size uint16) (mask uint16) {

	for i := uint16(0); i < size; i++ {
		mask = mask | (1 << uint16(i))

	}
	return
}

func bitFilterToInt(b []byte, offset uint16, size uint16, byteSize uint16) (data uint16) {

	mask := buildBitMask(size)
	// fmt.Println("size : ", size, " mask : ", mask)
	totalBitFromByte := byteSize * 8
	bForBitShift := byteToInt(b, byteSize)
	// fmt.Println("bForBitShift : ", bForBitShift)
	if size%8 != 0 {
		// fmt.Println("offset : ", offset)
		bitShift := totalBitFromByte - (offset) - size
		// fmt.Println("bitshift : ", bitShift)

		byteShifted := (bForBitShift >> bitShift)
		data = mask & byteShifted

	} else {
		byteShifted := (bForBitShift >> (totalBitFromByte - size))
		data = mask & byteShifted
	}
	// fmt.Println("data :", data)
	return data
}

func byteToInt(b []byte, size uint16) (data uint16) {

	for i := size; i > 0; i-- {
		data = data | (uint16(b[i-1]) << uint(8*(size-i)))

	}
	return
}
