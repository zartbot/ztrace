package geoip

import (
	"fmt"
	"math"
	"net"

	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

func ComputeDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R float64 = 6371
	Deg2Rad := math.Pi / 180.0
	dLat := (lat2 - lat1) * Deg2Rad
	dLon := (lon2 - lon1) * Deg2Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*Deg2Rad)*math.Cos(lat2*Deg2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * R * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

//GeoIPDB is the main struct of ip lookup engine
type GeoIPDB struct {
	CityDB *geoip2.Reader
	ASNDB  *geoip2.Reader
}

//New is used to create database
func New(cityDBPrefix string, asnDBPrefix string) *GeoIPDB {
	var r = GeoIPDB{}
	var err error
	r.CityDB, err = geoip2.Open(cityDBPrefix)
	r.ASNDB, err = geoip2.Open(asnDBPrefix)
	if err != nil {
		logrus.Fatal(err)
	}
	return &r
}

/* Use the following address test hongkong/taiwan/Macau
test_ip = "14.0.207.94"//hongkong :location_region_name == "Hong Kong"
test_ip = "140.112.110.1"//taiwan :location_region_name == "Taiwan"
test_ip = "122.100.160.253" //Macau :location_region_name == "Macau"
*/

//GeoLocation is the response type for location lookup
type GeoLocation struct {
	City      string
	Region    string
	Country   string
	ASN       uint
	SPName    string
	Latitude  float64
	Longitude float64
}

func (g GeoLocation) String() string {
	return fmt.Sprintf("City:  %-15s|Region:  %-15s|Country: %-15s |ASN: %8d::%-40s| Lat: %4.3f,%4.3f", g.City, g.Region, g.Country, g.ASN, g.SPName, g.Latitude, g.Longitude)
}

//Lookup is used to find IP location in GeoIPDB
func (g *GeoIPDB) Lookup(ipAddr string) GeoLocation {
	var r GeoLocation
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return r
	}
	c, _ := g.CityDB.City(ip)
	asn, _ := g.ASNDB.ASN(ip)
	if c.City.GeoNameID != 0 {
		r.City = c.City.Names["en"]
	}
	if len(c.Subdivisions) > 0 {
		if c.Subdivisions[0].GeoNameID != 0 {
			r.Region = c.Subdivisions[0].Names["en"]
		}
	}
	if c.Country.GeoNameID != 0 {
		r.Country = c.Country.Names["en"]
	}

	if r.Country == "Hong Kong" {
		r.Country = "China"
		r.Region = "Hong Kong"
		r.City = "Hong Kong"
	}
	if r.Country == "Macau" || r.Country == "Macao" {
		r.Country = "China"
		r.Region = "Macau"
		r.City = "Macau"
	}

	if r.Country == "Taiwan" {
		r.Country = "China"
		r.Region = "Taiwan"
	}

	r.Latitude = c.Location.Latitude
	r.Longitude = c.Location.Longitude
	r.ASN = asn.AutonomousSystemNumber
	r.SPName = asn.AutonomousSystemOrganization
	return r
}
