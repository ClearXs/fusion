package ip

import "regexp"

const (
	IpV4Regx = "(?:25[0-5]|2[0-4]\\\\d|1\\\\d\\\\d|[1-9]\\\\d|\\\\d)(?:\\\\.(?:25[0-5]|2[0-4]\\\\d|1\\\\d\\\\d|[1-9]\\\\d|\\\\d)){3}"
)

var (
	IpV4Pattern *regexp.Regexp
)

func Init() {
	IpV4Pattern, _ = regexp.Compile(IpV4Regx)
}

// IsIpV4 test ip address is ip v4
func IsIpV4(ip string) bool {
	return IpV4Pattern.MatchString(ip)
}
