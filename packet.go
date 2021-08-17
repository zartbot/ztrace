package ztrace

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func GetHash(src net.IP, dst net.IP, srcPort uint16, dstPort uint16) string {
	h := sha1.New()
	h.Write(src)
	h.Write(dst)
	p := make([]byte, 2)
	binary.BigEndian.PutUint16(p, srcPort)
	h.Write(p)
	binary.BigEndian.PutUint16(p, dstPort)
	h.Write(p)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type UDPHeader struct {
	Src    uint16
	Dst    uint16
	Length uint16
	Chksum uint16
}

type TCPHeader struct {
	Src        uint16
	Dst        uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8 // only use high 4 bits
	Flags      uint8 // only use low 6 bits
	Window     uint16
	Checksum   uint16
	Urgent     uint16
}

// pseudo header used for checksum calculation
type pseudohdr struct {
	ipsrc   []byte
	ipdst   []byte
	zero    uint8
	ipproto uint8
	plen    uint16
}

//checksum function for ip header
func checkSum(buf []byte) uint16 {
	sum := uint32(0)
	for n := 1; n < len(buf)-1; n += 2 {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	cksum := ^uint16(sum)
	return cksum
}

func (u *UDPHeader) checksum(ip *ipv4.Header, payload []byte) {
	u.Chksum = 0
	phdr := pseudohdr{
		ipsrc:   ip.Src,
		ipdst:   ip.Dst,
		zero:    0,
		ipproto: uint8(ip.Protocol),
		plen:    u.Length,
	}
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, &phdr)
	binary.Write(&b, binary.BigEndian, u)
	binary.Write(&b, binary.BigEndian, &payload)
	u.Chksum = checkSum(b.Bytes())
}

func (t *TraceRoute) BuildIPv4UDPkt(srcPort uint16, dstPort uint16, ttl uint8, id uint16, tos int) (*ipv4.Header, []byte) {
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		TOS:      tos,
		Len:      ipv4.HeaderLen,
		TotalLen: 60,
		ID:       int(id),
		Flags:    0,
		FragOff:  0,
		TTL:      int(ttl),
		Protocol: 17,
		Checksum: 0,
		Src:      t.netSrcAddr,
		Dst:      t.netDstAddr,
	}

	h, err := iph.Marshal()
	if err != nil {
		logrus.Fatal(err)
	}
	iph.Checksum = int(checkSum(h))

	udp := UDPHeader{
		Src:    srcPort,
		Dst:    dstPort,
		Length: 40, //hardcode payload length 32
	}

	payload := make([]byte, 32)
	for i := 0; i < 32; i++ {
		payload[i] = uint8(i + 64)
	}
	udp.checksum(iph, payload)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &udp)
	binary.Write(&buf, binary.BigEndian, &payload)
	return iph, buf.Bytes()
}
