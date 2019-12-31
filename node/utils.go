package node

import "strings"

func InListSubstring(x string, arr []string) bool {
	for _, y := range arr {
		if strings.Contains(x, y) {
			return true
		}
	}
	return false
}

func InList(x string, arr []string) bool {
	for _, y := range arr {
		if x == y {
			return true
		}
	}
	return false
}
