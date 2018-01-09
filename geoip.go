package main

import (
	"log"
	"net"
	"os"
	"sync"

	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
)

var (
	poolCity sync.Pool
	poolASN  sync.Pool
)

func init() {
	poolCity = sync.Pool{New: geoipOpener("geoip2.city")}
	poolASN = sync.Pool{New: geoipOpener("geoip2.asn")}

	viper.BindEnv("geoip2.city", "LOCATOR_GEOIP2_CITY")
	viper.SetDefault("geoip2.city", "GeoLite2-City.mmdb")

	viper.BindEnv("geoip2.asn", "LOCATOR_GEOIP2_ASN")
	viper.SetDefault("geoip2.asn", "GeoLite2-ASN.mmdb")
}

func geoipOpener(name string) func() interface{} {
	return func() interface{} {
		db, err := geoip2.Open(viper.GetString(name))
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		return db
	}
}

type geoipLocation struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type geoip struct {
	IP        string        `json:"ip"`
	Continent string        `json:"continent"`
	Country   string        `json:"country"`
	Region    string        `json:"region"`
	City      string        `json:"city"`
	Postal    string        `json:"postal"`
	ASN       uint          `json:"asn"`
	Org       string        `json:"org"`
	Location  geoipLocation `json:"location"`
}

func geoipAddr(addr net.IP) (*geoip, error) {
	var (
		dbCity = poolCity.Get().(*geoip2.Reader)
		dbASN  = poolASN.Get().(*geoip2.Reader)

		region string
		city   *geoip2.City
		asn    *geoip2.ASN
		err    error
	)
	defer poolCity.Put(dbCity)
	defer poolASN.Put(dbASN)

	city, err = dbCity.City(addr)
	if err != nil {
		return nil, err
	}

	asn, err = dbASN.ASN(addr)
	if err != nil {
		return nil, err
	}

	if city.Subdivisions != nil && len(city.Subdivisions) != 0 {
		region = city.Subdivisions[0].IsoCode
	}

	g := geoip{
		IP:        addr.String(),
		Continent: city.Continent.Code,
		Country:   city.Country.IsoCode,
		Region:    region,
		City:      city.City.Names["en"],
		Postal:    city.Postal.Code,
		Location: geoipLocation{
			Lon: city.Location.Longitude,
			Lat: city.Location.Latitude,
		},
		ASN: asn.AutonomousSystemNumber,
		Org: asn.AutonomousSystemOrganization,
	}

	return &g, nil
}
