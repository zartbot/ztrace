package ztrace

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func (t *TraceRoute) SendIPv4TCP(dport uint16) {
	sport := uint16(1000 + t.PortOffset + rand.Int31n(500))
	key := GetHash(t.netSrcAddr.To4(), t.netDstAddr.To4(), sport, dport, 6)
	db := NewStatsDB(key)

	t.DB.Store(key, db)
	go db.Cache.Run()
	conn, err := net.ListenPacket("ip4:tcp", t.netSrcAddr.String())
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	rSocket, err := ipv4.NewRawConn(conn)
	if err != nil {
		logrus.Fatal("can not create raw socket:", err)
	}
	seq := uint32(1000)
	mod := uint32(1 << 30)
	for {
		for ttl := 1; ttl <= int(t.MaxTTL); ttl++ {
			hdr, payload := t.BuildIPv4TCPSYN(sport, dport, uint8(ttl), seq, 0)
			rSocket.WriteTo(hdr, payload, nil)
			report := &SendMetric{
				FlowKey:   key,
				ID:        seq,
				TTL:       uint8(ttl),
				TimeStamp: time.Now(),
			}
			t.SendChan <- report
			seq = (seq + 4) % mod
		}
		atomic.AddUint64(db.SendCnt, 1)
		if atomic.LoadInt32(t.stopSignal) == 1 {
			break
		}
		time.Sleep(time.Microsecond * time.Duration(200000/t.PacketRate))
	}
}

//TODO add more on ICMP handle logic
func (t *TraceRoute) ListenIPv4TCP() {
	laddr := &net.IPAddr{IP: t.netSrcAddr}
	var err error
	t.recvTCPConn, err = net.ListenIP("ip4:tcp", laddr)
	if err != nil {
		logrus.Fatal("bind TCP failure:", err)
	}
	for {
		buf := make([]byte, 1500)
		n, raddr, err := t.recvTCPConn.ReadFrom(buf)
		if err != nil {
			break
		}
		if (n >= 20) && (n <= 100) {
			if (buf[13] == TCP_ACK+TCP_SYN) && (raddr.String() == t.netDstAddr.String()) {
				//no need to generate RST message, Linux will automatically send rst
				sport := binary.BigEndian.Uint16(buf[0:2])
				dport := binary.BigEndian.Uint16(buf[2:4])
				ack := binary.BigEndian.Uint32(buf[8:12]) - 1
				key := GetHash(t.netSrcAddr.To4(), t.netDstAddr.To4(), dport, sport, 6)
				m := &RecvMetric{
					FlowKey:   key,
					ID:        ack,
					RespAddr:  fmt.Sprintf("tcp:%s:%d", raddr.String(), sport),
					TimeStamp: time.Now(),
				}
				t.RecvChan <- m
			}

		}
	}
}

func (t *TraceRoute) ListenIPv4TCP_ICMP() {
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
		if (icmpType == 11 || (icmpType == 3 && buf[1] == 3)) && (n >= 36) { //TTL Exceeded or Port Unreachable
			seq := binary.BigEndian.Uint32(buf[32:36])
			dstip := net.IP(buf[24:28])
			srcip := net.IP(buf[20:24])
			srcPort := binary.BigEndian.Uint16(buf[28:30])
			dstPort := binary.BigEndian.Uint16(buf[30:32])
			if dstip.Equal(t.netDstAddr) { // && dstPort == t.dstPort {
				key := GetHash(srcip, dstip, srcPort, dstPort, 6)

				m := &RecvMetric{
					FlowKey:   key,
					ID:        seq,
					RespAddr:  raddr.String(),
					TimeStamp: time.Now(),
				}
				t.RecvChan <- m
			}
		}
	}
}
