package config

import (
	"os"
	"strings"
)

var allowDomains = []string{"localhost"}

func IsAllowedDomain(domain string) bool {
	return true // fixme for test
	/*
		for _, v := range allowDomains {
			if v == domain {
				return true
			}
		}
		return false
	*/
}

func init() {
	v, ok := os.LookupEnv("ALLOW_REDIRECT_DOMAIN")
	if ok {
		var domains []string
		for _, vv := range strings.Split(v, ",") {
			domains = append(domains, vv)
		}
		allowDomains = domains
	}
}
