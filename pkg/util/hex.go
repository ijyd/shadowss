package util

import "fmt"

//DumpHex dump []byte with hex
func DumpHex(b []byte) string {
	dumpStr := string("")

	for _, v := range b {
		dumpStr += fmt.Sprintf("0x%02x ", v)
	}
	return dumpStr
}
