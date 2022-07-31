package main

import (
	"flag"
	"fmt"
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
	flag.StringVar(&cmd.src, "src", cmd.src, "Source ")
	flag.IntVar(&cmd.maxPath, "path", cmd.maxPath, "Max ECMP Number")
	flag.IntVar(&cmd.maxPath, "p", cmd.maxPath, "Max ECMP Number")
	flag.IntVar(&cmd.maxTTL, "ttl", cmd.maxTTL, "Max TTL")
	flag.Float64Var(&cmd.pps, "rate", cmd.pps, "Packet Rate per second")
	flag.Float64Var(&cmd.pps, "r", cmd.pps, "Packet Rate per second")
	flag.BoolVar(&cmd.wmode, "wide", cmd.wmode, "Widescreen mode")
	flag.BoolVar(&cmd.wmode, "w", cmd.wmode, "Widescreen mode")
	flag.Parse()
}

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ./ztrace [-src source] [-proto protocol] [-ttl ttl] [-rate packetRate] [-wide Widescreen mode] [-path NumOfECMPPath] host")
	fmt.Println("Example:")
	fmt.Println(" ./ztrace www.cisco.com")
	fmt.Println(" ./ztrace -ttl 30 -rate 1 -path 8 -wide www.cisco.com")
	fmt.Println("Option:")
	flag.PrintDefaults()
}

func main() {

	if flag.NArg() != 1 {
		PrintUsage()
		return
	} else {
		cmd.dst = flag.Arg(0)
		fmt.Println(cmd.dst)
	}

	t := ztrace.New(cmd.protocol, cmd.dst, cmd.src, cmd.maxPath, uint8(cmd.maxTTL), float32(cmd.pps), 0, cmd.wmode, "geoip/asn.mmdb", "geoip/geoip.mmdb")
	t.Latitude = cmd.lat
	t.Longitude = cmd.long

	t.Start()
	go t.Report(time.Second)
	time.Sleep(time.Second * 100)
	t.Stop()
}
