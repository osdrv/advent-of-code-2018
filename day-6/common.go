package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const (
	ALOT    = int(999999999)
	ALOT32u = uint32(4294967295)
	ALOT32  = int32(2147483647)
	ALOT64u = uint64(18446744073709551615)
	ALOT64  = int64(9223372036854775807)
)

var (
	STEPS4 = [][2]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}

	STEPS8 = [][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
)

func noerr(err error) {
	if err != nil {
		panic(fmt.Sprintf("unhandled error: %s", err))
	}
}

func assert(check bool, msg string) {
	if !check {
		panic(fmt.Sprintf("assert %q failed", msg))
	}
}

func parseInt(s string) int {
	num, err := strconv.Atoi(s)
	noerr(err)
	return num
}

func readFile(in io.Reader) string {
	data, err := ioutil.ReadAll(in)
	noerr(err)
	return trim(string(data))
}

func readLines(in io.Reader) []string {
	scanner := bufio.NewScanner(in)
	lines := make([]string, 0, 1)
	for scanner.Scan() {
		lines = append(lines, trim(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scan failed: %s", err))
	}
	return lines
}

func trim(s string) string {
	return strings.TrimRight(s, "\t\n\r")
}

func parseInts(s string) []int {
	chs := strings.FieldsFunc(trim(s), func(r rune) bool {
		return r == ' ' || r == ',' || r == '\t'
	})
	nums := make([]int, 0, len(chs))
	for i := 0; i < len(chs); i++ {
		nums = append(nums, parseInt(chs[i]))
	}
	return nums
}

type Point2 struct {
	x, y int
}

func NewPoint2(x, y int) *Point2 {
	return &Point2{x, y}
}

func (p2 *Point2) String() string {
	return fmt.Sprintf("P2{%d, %d}", p2.x, p2.y)
}

type Point3 struct {
	x, y, z int
}

func (p3 *Point3) String() string {
	return fmt.Sprintf("P3{%d, %d, %d}", p3.x, p3.y, p3.z)
}

func NewPoint3(x, y, z int) *Point3 {
	return &Point3{x, y, z}
}

func makeIntField(h, w int) [][]int {
	res := make([][]int, h)
	for i := 0; i < h; i++ {
		res[i] = make([]int, w)
	}
	return res
}

func makeByteField(h, w int) [][]byte {
	res := make([][]byte, h)
	for i := 0; i < h; i++ {
		res[i] = make([]byte, w)
	}
	return res
}

func sizeIntField(field [][]int) (int, int) {
	rows, cols := len(field), 0
	if rows > 0 {
		cols = len(field[0])
	}
	return rows, cols
}

func sizeByteField(field [][]byte) (int, int) {
	rows, cols := len(field), 0
	if rows > 0 {
		cols = len(field[0])
	}
	return rows, cols
}

func copyIntField(field [][]int) [][]int {
	cp := makeIntField(sizeIntField(field))
	for i := 0; i < len(field); i++ {
		copy(cp[i], field[i])
	}
	return cp
}

func copyByteField(field [][]byte) [][]byte {
	cp := makeByteField(sizeByteField(field))
	for i := 0; i < len(field); i++ {
		copy(cp[i], field[i])
	}
	return cp
}

func printIntField(field [][]int, sep string) string {
	return printIntFieldWithSubs(field, sep, make(map[int]string))
}

func printIntFieldWithSubs(field [][]int, sep string, subs map[int]string) string {
	var buf bytes.Buffer
	rows, cols := sizeIntField(field)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if j > 0 {
				buf.WriteString(sep)
			}
			if sub, ok := subs[field[i][j]]; ok {
				buf.WriteString(sub)
			} else {
				buf.WriteByte('0' + byte(field[i][j]))
			}
		}
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	return buf.String()
}

func printByteField(field [][]byte, sep string) string {
	return printByteFieldWithSubs(field, sep, make(map[byte]string))
}

func printByteFieldWithSubs(field [][]byte, sep string, subs map[byte]string) string {
	var buf bytes.Buffer
	rows, cols := sizeByteField(field)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if j > 0 {
				buf.WriteString(sep)
			}
			if sub, ok := subs[field[i][j]]; ok {
				buf.WriteString(sub)
			} else {
				buf.WriteByte('0' + field[i][j])
			}
		}
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	return buf.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// functions to compute local extremums

func findLocalMin(n int, compute func(i int) int) int {
	a, b := 0, n-1
	leftix, midix, rightix := a, (a+b)/2, b
	left, mid, right := compute(leftix), compute(midix), compute(rightix)
	for rightix-leftix > 1 {
		if left <= mid && mid <= right {
			b = midix
			leftix, midix, rightix = a, (a+midix)/2, midix
			left, mid, right = compute(leftix), compute(midix), mid
		} else if left >= mid && mid >= right {
			a = midix
			leftix, midix, rightix = midix, (midix+b)/2, b
			left, mid, right = right, compute(midix), compute(rightix)
		} else {
			a = leftix
			b = rightix
			leftix, rightix = (leftix+midix)/2, (midix+rightix)/2
			left, right = compute(leftix), compute(rightix)
		}
	}
	return min(left, right)
}

func findLocalMax(n int, compute func(i int) int) int {
	return -1 * findLocalMin(n, func(i int) int {
		return -1 * compute(i)
	})
}

// slice helpers

func mapIntArr(arr []int, mapfn func(int) int) []int {
	res := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = mapfn(arr[i])
	}
	return res
}

func mapByteArr(arr []byte, mapfn func(byte) byte) []byte {
	res := make([]byte, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = mapfn(arr[i])
	}
	return res
}

func reverseIntArr(arr []int) []int {
	res := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		res[len(arr)-1-i] = arr[i]
	}
	return res
}

func reverseByteArr(arr []byte) []byte {
	res := make([]byte, len(arr))
	for i := 0; i < len(arr); i++ {
		res[len(arr)-1-i] = arr[i]
	}
	return res
}

func grepIntArr(arr []int, grepfn func(int) bool) []int {
	res := make([]int, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		if grepfn(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return res
}

func grepByteArr(arr []byte, grepfn func(byte) bool) []byte {
	res := make([]byte, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		if grepfn(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return res
}

// logging function

func printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
