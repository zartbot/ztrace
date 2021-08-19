package main

import (
	"flag"
	"time"

	"github.com/zartbot/ztrace"
)

var cmd = struct {
	dst     string
	maxPath int
	maxTTL  int
	pps     float64
	wmode   bool
}{
	"8.8.8.8",
	16,
	64,
	1,
	false,
}

func init() {
	flag.StringVar(&cmd.dst, "d", cmd.dst, "Destination ")
	flag.IntVar(&cmd.maxPath, "p", cmd.maxPath, "Max ECMP Number")
	flag.IntVar(&cmd.maxTTL, "t", cmd.maxTTL, "Max TTL")
	flag.Float64Var(&cmd.pps, "r", cmd.pps, "Packet Rate per second")
	flag.BoolVar(&cmd.wmode, "w", cmd.wmode, "Widescreen mode")
	flag.Parse()
}

func main() {
	t := ztrace.New(cmd.dst, "", cmd.maxPath, uint8(cmd.maxTTL), float32(cmd.pps), 0, cmd.wmode, "geoip/asn.mmdb", "geoip/geoip.mmdb")
	t.Start()
	go t.Report(time.Second)
	time.Sleep(time.Second * 100)
	t.Stop()
}
