package cmn

import "strings"

func IndexAt(s, sep string, n int) int {
	idx := strings.Index(s[n:], sep)
	if idx > -1 {
		idx += n
	}
	return idx
}

func LastIndexAt(s, subStr string, n int) int {
	idx := strings.LastIndex(s[:n], subStr)
	return idx
}

func LastIndexAtStr(s, subStr string, atStr string) int {
	n := strings.LastIndex(s, atStr)
	if n > -1 {
		return LastIndexAt(s, subStr, n)
	} else {
		return -1
	}
}
