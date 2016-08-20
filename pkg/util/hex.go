package util

import "fmt"

//DumpHex dump []byte with hex
func DumpHex(b []byte) string {
	dumpStr := string("")

	i := 0
	for _, v := range b {
		i++
		dumpStr += fmt.Sprintf("0x%02x ", v)
		if 0 == (i % 16) {
			dumpStr += fmt.Sprintf("\r\n")
		}

	}
	return dumpStr
}
