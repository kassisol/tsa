package host

import (
	"net/url"
	"strings"
)

func New(url *url.URL, host string) string {
	hostname1 := url.Hostname()
	hostname2 := strings.Split(host, ":")[0]

	if hostname1 == hostname2 {
		return hostname1
	}

	return hostname2
}
