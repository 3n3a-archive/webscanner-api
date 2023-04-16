package scanner

import (
	"fmt"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

const (
	LANG_CODE = "en"
)

func (s *ScanClient) GetIPInfo(ip net.IP) (IPInfo, error) {
	GEOIP_CITY_DB_PATH := os.Getenv("GEOIP_CITY_DB_PATH")
	GEOIP_ASN_DB_PATH := os.Getenv("GEOIP_ASN_DB_PATH")
	cityDb, err := geoip2.Open(GEOIP_CITY_DB_PATH)
	if err != nil {
		fmt.Println("GeoIP City Database couldn't be read")
		return IPInfo{}, err
	}
	defer cityDb.Close()

	asnDb, err := geoip2.Open(GEOIP_ASN_DB_PATH)
	if err != nil {
		fmt.Println("GeoIP ASN Database couldn't be read")
		return IPInfo{}, err
	}
	defer asnDb.Close()

	city, err := cityDb.City(ip)
	if err != nil {
		fmt.Println("Getting city for ip failed")
		return IPInfo{}, err
	}

	asn, err := asnDb.ASN(ip)
	if err != nil {
		fmt.Println("Getting asn for ip failed")
		return IPInfo{}, err
	}

	return IPInfo{
		City: city.City.Names[LANG_CODE],
		Country: city.Country.Names[LANG_CODE],
		ASN: asn.AutonomousSystemOrganization,
	}, nil
}