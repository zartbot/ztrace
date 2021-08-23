package ztrace

import (
	"encoding/binary"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func (t *TraceRoute) SendIPv4UDP() {
	dport := uint16(33434 + rand.Int31n(64))
	sport := uint16(1000 + t.PortOffset + rand.Int31n(500))
	key := GetHash(t.netSrcAddr.To4(), t.netDstAddr.To4(), sport, dport, 17)
	db := NewStatsDB(key)

	t.DB.Store(key, db)
	go db.Cache.Run()
	conn, err := net.ListenPacket("ip4:udp", t.netSrcAddr.String())
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	rSocket, err := ipv4.NewRawConn(conn)
	if err != nil {
		logrus.Fatal("can not create raw socket:", err)
	}
	id := uint16(1)
	mod := uint16(1 << 15)
	for {
		for ttl := 1; ttl <= int(t.MaxTTL); ttl++ {
			hdr, payload := t.BuildIPv4UDPkt(sport, dport, uint8(ttl), id, 0)
			id = (id + 1) % mod
			rSocket.WriteTo(hdr, payload, nil)
			report := &SendMetric{
				FlowKey:   key,
				ID:        uint32(hdr.ID),
				TTL:       uint8(ttl),
				TimeStamp: time.Now(),
			}
			t.SendChan <- report
		}
		atomic.AddUint64(db.SendCnt, 1)
		if atomic.LoadInt32(t.stopSignal) == 1 {
			break
		}
		time.Sleep(time.Microsecond * time.Duration(1000000/t.PacketRate))
	}
}

func (t *TraceRoute) ListenIPv4UDP() {
	laddr := &net.IPAddr{IP: t.netSrcAddr}
	var err error
	t.recvICMPConn, err = net.ListenIP("ip4:icmp", laddr)
	if err != nil {
		logrus.Fatal("bind failure:", err)
	}
	for {
		buf := make([]byte, 1500)
		n, raddr, err := t.recvICMPConn.ReadFrom(buf)
		if err != nil {
			break
		}
		icmpType := buf[0]
		//logrus.Info(raddr, "|", icmpType, "|", n)
		if (icmpType == 11 || (icmpType == 3 && buf[1] == 3)) && (n >= 36) { //TTL Exceeded or Port Unreachable
			id := binary.BigEndian.Uint16(buf[12:14])
			dstip := net.IP(buf[24:28])
			srcip := net.IP(buf[20:24])
			srcPort := binary.BigEndian.Uint16(buf[28:30])
			dstPort := binary.BigEndian.Uint16(buf[30:32])
			if dstip.Equal(t.netDstAddr) { // && dstPort == t.dstPort {
				key := GetHash(srcip, dstip, srcPort, dstPort, 17)

				m := &RecvMetric{
					FlowKey:   key,
					ID:        uint32(id),
					RespAddr:  raddr.String(),
					TimeStamp: time.Now(),
				}
				t.RecvChan <- m
			}
		}
	}
}
