package main

import (
	"bytes"
	"strconv"
)

type Game struct {
	nums []byte
	i, j int
}

func NewGame() *Game {
	return &Game{
		nums: []byte{3, 7},
		i:    0,
		j:    1,
	}
}

func (g *Game) Play() {
	sum := g.nums[g.i] + g.nums[g.j]
	if sum >= 10 {
		g.nums = append(g.nums, 1)
	}
	g.nums = append(g.nums, sum%10)
	i := (g.i + 1 + int(g.nums[g.i])) % len(g.nums)
	j := (g.j + 1 + int(g.nums[g.j])) % len(g.nums)
	g.i = i
	g.j = j
}

func (g *Game) String() string {
	var buf bytes.Buffer
	for ix, num := range g.nums {
		v := strconv.Itoa(int(num))
		if ix == g.i {
			v = "(" + v + ")"
		} else if ix == g.j {
			v = "[" + v + "]"
		}
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(v)
	}
	return buf.String()
}

func (g *Game) Search(needle []byte, offset int) int {
	i := offset
	for i < len(g.nums) {
		if g.nums[i] == needle[0] {
			j := 1
			for j < len(needle) {
				if len(g.nums) <= i+j {
					return -1
				}
				if g.nums[i+j] != needle[j] {
					break
				}
				j++
			}
			if j == len(needle) {
				return i
			}
		}
		i++
	}
	return -1
}

const (
	DIGITS = 10
	RCPTS  = 290431
)

func part1() {
	game := NewGame()
	for len(game.nums) < RCPTS+DIGITS {
		game.Play()
		//printf("%s", game)
	}
	print("result: ")
	for _, n := range game.nums[RCPTS : RCPTS+DIGITS] {
		print(strconv.Itoa(int(n)))
	}
	println("")
}

func part2() {
	game := NewGame()
	//needle := []byte{5, 1, 5, 8, 9}
	//needle := []byte{0, 1, 2, 4, 5}
	//needle := []byte{9, 2, 5, 1, 0}
	//needle := []byte{5, 9, 4, 1, 4}
	needle := []byte{2, 9, 0, 4, 3, 1}
	off := 0
	for {
		game.Play()
		ix := game.Search(needle, off)
		if ix >= 0 {
			printf("search ix: %d", ix)
			break
		}
		off = max(0, len(game.nums)-1-len(needle))
	}
}

func main() {
	part1()
	part2()
}
