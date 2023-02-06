package scanner

import (
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
	"github.com/temoto/robotstxt"
)

type ScanReport struct {
	SecurityTxt SecurityTxtParser.SecurityTxt
	RobotsTxt robotstxt.RobotsData
}