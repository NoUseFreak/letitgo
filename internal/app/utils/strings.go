package utils

import "bytes"

func NormalizeNewlines(s string) string {
	d := []byte(s)
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)

	return string(d)
}
