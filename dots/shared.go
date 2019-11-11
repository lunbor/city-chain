package dots

import "strings"

func Remove0x(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return strings.TrimPrefix(hex, "0x")
	}
	return hex
}
