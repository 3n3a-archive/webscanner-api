package scanner

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

const (
	GEOIP_CITY_DB_PATH = "./geodb/GeoLite2-City.mmdb"
	GEOIP_ASN_DB_PATH = "./geodb/GeoLite2-ASN.mmdb"
	LANG_CODE = "en"
)

func (s *ScanClient) GetIPInfo(ip net.IP) (IPInfo, error) {
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