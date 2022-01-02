package main

import (
	"bytes"
	"os"
)

func countLetters(s string) [26]int {
	var res [26]int
	for i := 0; i < len(s); i++ {
		res[int(s[i]-'a')]++
	}
	return res
}

func hasNLetters(s string, n int) bool {
	cnt := countLetters(s)
	for i := 0; i < len(cnt); i++ {
		if cnt[i] == n {
			return true
		}
	}
	return false
}

func computeCommon(s1, s2 string) (string, int) {
	var buf bytes.Buffer
	dist := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			dist++
		} else {
			buf.WriteByte(s1[i])
		}
	}
	return buf.String(), dist
}

func findCommonId(lines []string) string {
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			common, dist := computeCommon(lines[i], lines[j])
			if dist == 1 {
				return common
			}
		}
	}
	return "<NONE>"
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	cnt2, cnt3 := 0, 0
	for _, line := range lines {
		if hasNLetters(line, 2) {
			cnt2++
		}
		if hasNLetters(line, 3) {
			cnt3++
		}
	}

	printf("%d * %d = %d", cnt2, cnt3, cnt2*cnt3)

	id := findCommonId(lines)
	printf("the id is: %s", id)
}
