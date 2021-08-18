package ztrace

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zartbot/ztrace/geoip"
	"github.com/zartbot/ztrace/stats/describe"
	"github.com/zartbot/ztrace/stats/quantile"
)

type ServerRecord struct {
	TTL             uint8
	Addr            string
	Name            string
	Session         string
	GeoLocation     geoip.GeoLocation
	LatencyDescribe *describe.Item
	Quantile        *quantile.Stream
	RecvCnt         uint64
	Lock            *sync.Mutex
}

func (s *ServerRecord) LookUPAddr() {
	rA, _ := net.LookupAddr(s.Addr)
	var buf bytes.Buffer
	for _, item := range rA {
		if len(item) > 0 {
			//some platform may add dot in suffix
			item = strings.TrimSuffix(item, ".")
			if !strings.HasSuffix(item, ".in-addr.arpa") {
				buf.WriteString(item)
			}
		}
	}
	s.Name = buf.String()
}
func (t *TraceRoute) NewServerRecord(ipaddr string, ttl uint8, key string) *ServerRecord {
	r := &ServerRecord{
		TTL:             ttl,
		Addr:            ipaddr,
		LatencyDescribe: describe.New(),
		Session:         key,
		Quantile: quantile.NewTargeted(map[float64]float64{
			0.50: 0.005,
			0.90: 0.001,
			0.99: 0.0001,
		}),
		RecvCnt: 0,
		Lock:    &sync.Mutex{},
	}
	r.GeoLocation = t.geo.Lookup(ipaddr)
	return r
}

func (t *TraceRoute) Stats() {
	for {
		select {
		case v := <-t.SendChan:

			tdb, ok := t.DB.Load(v.FlowKey)
			if !ok {
				continue
			}
			db := tdb.(*StatsDB)
			db.Cache.Store(v.ID, v, v.TimeStamp)

		case v := <-t.RecvChan:

			tdb, ok := t.DB.Load(v.FlowKey)
			if !ok {
				continue
			}
			db := tdb.(*StatsDB)
			tsendInfo, valid := db.Cache.Load(v.ID)
			if !valid {
				continue
			}
			sendInfo := tsendInfo.(*SendMetric)
			server, valid := t.Metric[sendInfo.TTL][v.RespAddr]
			//create server
			if !valid {
				server = t.NewServerRecord(v.RespAddr, uint8(sendInfo.TTL), sendInfo.FlowKey)
				t.Metric[sendInfo.TTL][v.RespAddr] = server
			}

			server.Lock.Lock()
			server.RecvCnt++
			latency := float64(v.TimeStamp.Sub(sendInfo.TimeStamp) / time.Microsecond)
			//logrus.Info(v.RespAddr, ":", latency)
			server.LatencyDescribe.Append(latency, 2)
			server.Quantile.Insert(latency)
			server.Lock.Unlock()
			if server.Name == "" {
				go server.LookUPAddr()
			}

		}
		if atomic.LoadInt32(t.stopSignal) == 1 {
			return
		}
	}
}

func (t *TraceRoute) Print() {
	fmt.Printf("\033[H\033[2J")
	fmt.Printf("\n   Traceroute Report\n\n")
	table := tablewriter.NewWriter(os.Stdout)

	if t.WideMode {
		table.SetHeader([]string{"TTL ", "Server", "Name", "City", "Country", "ASN", "SP", "p95", "Latency", "Jitter", "Loss"})
	} else {
		table.SetHeader([]string{"TTL ", "Server", "Name", "Country", "SP", "p95", "Latency", "Jitter", "Loss"})

	}
	table.SetAutoFormatHeaders(false)
	for ttl := 1; ttl <= int(t.MaxTTL); ttl++ {
		StopFlag := false
		firstLine := fmt.Sprintf("%3d", ttl)
		for _, v := range t.Metric[ttl] {

			latency := fmt.Sprintf("%-8.2fms", v.LatencyDescribe.Mean/1000)
			jitter := fmt.Sprintf("%-8.2fms", v.LatencyDescribe.Std()/1000)
			p95 := fmt.Sprintf("%12.2fms", v.Quantile.Query(0.95)/1000)
			tdb, ok := t.DB.Load(v.Session)

			loss := float32(0)
			if ok {
				statsDB := tdb.(*StatsDB)
				sendCnt := atomic.LoadUint64(statsDB.SendCnt)
				if sendCnt != 0 {
					loss = (1 - float32(v.RecvCnt)/float32(sendCnt)) * 100
				}

				if v.RecvCnt > sendCnt {
					loss = 0
				}
			}

			city := fmt.Sprintf("%-16.16s", v.GeoLocation.City)
			country := fmt.Sprintf("%-16.16s", v.GeoLocation.Country)
			asn := fmt.Sprintf("%-10d", v.GeoLocation.ASN)

			sp := fmt.Sprintf("%-16.16s", v.GeoLocation.SPName)
			saddr := fmt.Sprintf("%-16.16s", v.Addr)
			sname := fmt.Sprintf("%-26.26s", v.Name)
			if t.WideMode {
				table.Append([]string{firstLine, saddr, sname, city, country, asn, sp, p95, latency, jitter, fmt.Sprintf("%-3.1f%%", loss)})

			} else {
				table.Append([]string{firstLine, saddr, sname, country, sp, p95, latency, jitter, fmt.Sprintf("%-3.1f%%", loss)})
			}
			if firstLine != "" {
				firstLine = ""
			}
			if v.Addr == t.netDstAddr.String() {
				StopFlag = true
				break
			}

		}
		if StopFlag {
			break
		}

	}
	table.Render()

}

func (t *TraceRoute) Report() {
	for {
		t.Print()
		if atomic.LoadInt32(t.stopSignal) == 1 {
			return
		}
		time.Sleep(time.Second)
	}
}

/*

fmt.Printf("\033[H\033[2J")
	fmt.Printf("\n   SDWAN Performance Test Report\n\n")

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Stats ", "Latency(ms)", "Bandwidth(Per Session)"})
	table.SetAutoFormatHeaders(false)
	table.Append([]string{"mean", fmt.Sprintf("%20.2fms", LatencyStats.Mean), fmt.Sprintf("%20.2fMbps", BWStats.Mean)})
	table.Append([]string{"Jitter", fmt.Sprintf("%20.2fms", LatencyStats.Std()), ""})
	table.Append([]string{"", "", ""})
	table.Append([]string{"Min", fmt.Sprintf("%20.2fms", LatencyStats.Min), fmt.Sprintf("%20.2fMbps", BWStats.Min)})
	table.Append([]string{"p25", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.25)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.25))})
	table.Append([]string{"p75", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.75)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.75))})
	table.Append([]string{"p90", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.90)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.90))})
	table.Append([]string{"p95", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.95)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.95))})
	table.Append([]string{"p99", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.99)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.99))})
	table.Append([]string{"Max", fmt.Sprintf("%20.2fms", LatencyStats.Max), fmt.Sprintf("%20.2fMbps", BWStats.Max)})
	table.SetFooter([]string{fmt.Sprintf("Count: %d", LatencyQuantile.Count()), fmt.Sprintf("Error: %d | Timeout: %d", errors, timeouts), fmt.Sprintf("Total-BW: %10.2fMbps", BWStats.Mean*float64(cli.ClientNum))})

	table.Render()

func (p *Proxy) TableRender(name string, addrList map[string]string) {
	fmt.Printf("[%s] DNS Lookup Result\n\n", name)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Addresss ", "ASN", "City", "Region", "Country", "Location", "Distance(KM)", "DNS Server"})
	table.SetAutoFormatHeaders(false)

	for k, v := range addrList {
		result := p.geo.Lookup(k)
		distance := geoip.ComputeDistance(31.02, 121.26, result.Latitude, result.Longitude)
		table.Append([]string{k, fmt.Sprintf("%-30.30s", result.SPName), fmt.Sprintf("%-16.16s", result.City), fmt.Sprintf("%-16.16s", result.Region), fmt.Sprintf("%-16.16s", result.Country), fmt.Sprintf("%6.2f , %6.2f", result.Latitude, result.Longitude), fmt.Sprintf("%8.0f", distance), v})
	}
	table.Render()
}
*/
