package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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

var (
	DEBUG = os.Getenv("DEBUG") == "1"
)

type Number interface {
	byte | int | int32 | int64 | uint32 | uint64 | float64
}

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

func makeNumField[N Number](h, w int) [][]N {
	res := make([][]N, h)
	for i := 0; i < h; i++ {
		res[i] = make([]N, w)
	}
	return res
}

func makeIntField(h, w int) [][]int {
	return makeNumField[int](h, w)
}

func makeByteField(h, w int) [][]byte {
	return makeNumField[byte](h, w)
}

func sizeNumField[N Number](field [][]N) (int, int) {
	rows, cols := len(field), 0
	if rows > 0 {
		cols = len(field[0])
	}
	return rows, cols
}

// Deprecated: please use `sizeNumField` instead.
func sizeIntField(field [][]int) (int, int) {
	return sizeNumField(field)
}

// Deprecated: please use `sizeNumField` instead.
func sizeByteField(field [][]byte) (int, int) {
	return sizeNumField(field)
}

func copyNumField[N Number](field [][]N) [][]N {
	cp := makeNumField[N](sizeNumField(field))
	for i := 0; i < len(field); i++ {
		copy(cp[i], field[i])
	}
	return cp
}

// Deprecated: please use `copyNumField` instead.
func copyIntField(field [][]int) [][]int {
	return copyNumField(field)
}

// Deprecated: please use `copyNumField` instead.
func copyByteField(field [][]byte) [][]byte {
	return copyNumField(field)
}

func printNumField[N Number](field [][]N, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[N]string))
}

// Deprecated: please use `printNumField` instead.
func printIntField(field [][]int, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[int]string))
}

// Deprecated: please use `printNumField` instead.
func printByteField(field [][]byte, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[byte]string))
}

func printNumFieldWithSubs[N Number](field [][]N, sep string, subs map[N]string) string {
	var buf bytes.Buffer
	rows, cols := sizeNumField(field)
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

func printIntFieldWithSubs(field [][]int, sep string, subs map[int]string) string {
	return printNumFieldWithSubs(field, sep, subs)
}

func printByteFieldWithSubs(field [][]byte, sep string, subs map[byte]string) string {
	return printNumFieldWithSubs(field, sep, subs)
}

func min[N Number](a, b N) N {
	if a < b {
		return a
	}
	return b
}

func max[N Number](a, b N) N {
	if a > b {
		return a
	}
	return b
}

func abs[N Number](v N) N {
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

func reverseNumArr[N Number](arr []N) []N {
	res := make([]N, len(arr))
	for i := 0; i < len(arr); i++ {
		res[len(arr)-1-i] = arr[i]
	}
	return res
}

// Deprecated: please use `reverseNumArr` instead.
func reverseIntArr(arr []int) []int {
	return reverseNumArr(arr)
}

// Deprecated: please use `reverseNumArr` instead.
func reverseByteArr(arr []byte) []byte {
	return reverseByteArr(arr)
}

func grepNumArr[N Number](arr []N, grepfn func(N) bool) []N {
	res := make([]N, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		if grepfn(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return res
}

// Deprecated: please use `grepNumArr` instead.
func grepIntArr(arr []int, grepfn func(int) bool) []int {
	return grepNumArr(arr, grepfn)
}

// Deprecated: please use `grepNumArr` instead.
func grepByteArr(arr []byte, grepfn func(byte) bool) []byte {
	return grepNumArr(arr, grepfn)
}

// logging function

func debugf(format string, v ...interface{}) {
	if DEBUG {
		log.Printf(format, v...)
	}
}

func printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Data types

type BinHeap[T comparable] struct {
	items []T
	index map[T]int
	cmp   func(a, b T) bool
}

func NewBinHeap[T comparable](cmp func(a, b T) bool) *BinHeap[T] {
	return &BinHeap[T]{
		items: make([]T, 0, 1),
		index: make(map[T]int),
		cmp:   cmp,
	}
}

func (h *BinHeap[T]) Size() int {
	return len(h.items)
}

func (h *BinHeap[T]) Push(item T) {
	last := len(h.items)
	if _, ok := h.index[item]; !ok {
		h.items = append(h.items, item)
		h.index[item] = last
	}
	ptr := h.index[item]
	h.reheapAt(ptr)
}

func (h *BinHeap[T]) Pop() T {
	last := len(h.items) - 1
	h.swap(0, last)
	item := h.items[last]
	h.items = h.items[:last]
	delete(h.index, item)
	h.reheapAt(0)

	return item
}

func (h *BinHeap[T]) swap(i, j int) {
	h.index[h.items[i]], h.index[h.items[j]] = h.index[h.items[j]], h.index[h.items[i]]
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *BinHeap[T]) reheapAt(ptr int) {
	for ptr > 0 {
		parent := (ptr - 1) / 2
		if h.cmp(h.items[ptr], h.items[parent]) {
			h.swap(ptr, parent)
			ptr = parent
		} else {
			break
		}
	}

	for ptr < len(h.items) {
		ch1, ch2 := ptr*2+1, ptr*2+2
		next := ptr
		if ch1 < len(h.items) && h.cmp(h.items[ch1], h.items[next]) {
			next = ch1
		}
		if ch2 < len(h.items) && h.cmp(h.items[ch2], h.items[next]) {
			next = ch2
		}
		if next != ptr {
			h.swap(ptr, next)
			ptr = next
		} else {
			break
		}
	}
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
