package utils

import (
	"log"
	"strings"
)

var threePointDomains = []string{
	".gov.cn",
	".net.cn",
	".com.cn",
	".org.cn",
	".co.uk",
}

func ParseDomain(domain string) (rootDomain, levelsDomain string) {
	dotNum := strings.Count(domain, ".")
	if dotNum < 1 {
		log.Fatalln("Domain format error")
		return
	}
	if dotNum == 1 {
		rootDomain = domain
		return
	}

	dotTime := 2
	for _, item := range threePointDomains {
		if strings.HasSuffix(domain, item) {
			dotTime = 3
			break
		}
	}

	strSlice := strings.Split(domain, ".")

	if len(strSlice) == dotTime {
		rootDomain = domain
		return
	}
	for index, value := range strSlice {
		if index == 0 {
			levelsDomain = value
			continue
		}
		rootDomain += "." + value
	}
	rootDomain = strings.Trim(rootDomain, ".")
	return
}