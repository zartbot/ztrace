package main

import (
	"flag"
	"time"

	"github.com/zartbot/ztrace"
)

var cmd = struct {
	dst      string
	src      string
	protocol string
	maxPath  int
	maxTTL   int
	pps      float64
	wmode    bool
	lat      float64
	long     float64
}{
	"",
	"",
	"udp",
	16,
	64,
	1,
	false,
	31.02,
	121.1,
}

func init() {
	flag.StringVar(&cmd.protocol, "proto", cmd.protocol, "Protocol[icmp|tcp|udp]")
	flag.StringVar(&cmd.dst, "dest", cmd.dst, "Destination ")
	flag.StringVar(&cmd.dst, "src", cmd.dst, "Source ")
	flag.IntVar(&cmd.maxPath, "path", cmd.maxPath, "Max ECMP Number")
	flag.IntVar(&cmd.maxTTL, "ttl", cmd.maxTTL, "Max TTL")
	flag.Float64Var(&cmd.pps, "rate", cmd.pps, "Packet Rate per second")
	flag.BoolVar(&cmd.wmode, "wide", cmd.wmode, "Widescreen mode")
	flag.Parse()
}

func main() {
	if cmd.dst == "" {
		flag.PrintDefaults()
		return
	}
	t := ztrace.New(cmd.protocol, cmd.dst, cmd.src, cmd.maxPath, uint8(cmd.maxTTL), float32(cmd.pps), 0, cmd.wmode, "geoip/asn.mmdb", "geoip/geoip.mmdb")
	t.Latitude = cmd.lat
	t.Longitude = cmd.long

	t.Start()
	go t.Report(time.Second)
	time.Sleep(time.Second * 100)
	t.Stop()
}
