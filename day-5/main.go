package main

import (
	"bytes"
	"os"
)

func willReact(a, b byte) bool {
	if a >= 'a' && a <= 'z' {
		return b == a-('a'-'A')
	} else if a >= 'A' && a <= 'Z' {
		return b == a+('a'-'A')
	}
	panic("should not happen")
}

func react(s string) string {
	polym := []byte(s)

	ptr := 0
	for ptr < len(polym) {
		if ptr < 0 {
			ptr = 0
		}
		if ptr < len(polym)-1 {
			if willReact(polym[ptr], polym[ptr+1]) {
				polym = append(polym[:ptr], polym[ptr+2:]...)
				//printf("new poly: %s", string(polym))
				ptr--
				continue
			}
		}
		ptr++
	}
	return string(polym)
}

func strip(s string, cut string) string {
	reject := make(map[byte]struct{})
	for i := 0; i < len(cut); i++ {
		reject[cut[i]] = struct{}{}
	}
	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		if _, ok := reject[s[i]]; !ok {
			buf.WriteByte(s[i])
		}
	}
	return buf.String()
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()
	data := trim(readFile(f))

	react1 := react(data)
	printf("final poly: %s (%d)", react1, len(react1))

	maxLen := len(react1)
	for ch := byte('A'); ch <= byte('Z'); ch++ {
		s2 := strip(data, string([]byte{ch, ch + 'a' - 'A'}))
		r2 := react(s2)
		if len(r2) < maxLen {
			maxLen = len(r2)
		}
	}

	printf("maxlen: %d", maxLen)
}
