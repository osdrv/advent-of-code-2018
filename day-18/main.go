package main

import (
	"os"
)

const (
	OPEN int = iota
	TREES
	LUMBER
)

func evolve(field [][]int) [][]int {
	newfield := makeIntField(len(field), len(field[0]))

	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			lumber, trees := 0, 0
			for _, step := range STEPS8 {
				ni, nj := i+step[0], j+step[1]
				if ni >= 0 && ni < len(field) && nj >= 0 && nj < len(field[0]) {
					switch field[ni][nj] {
					case TREES:
						trees++
					case LUMBER:
						lumber++
					}
				}
			}
			switch field[i][j] {
			case OPEN:
				if trees >= 3 {
					newfield[i][j] = TREES
				} else {
					newfield[i][j] = field[i][j]
				}
			case TREES:
				if lumber >= 3 {
					newfield[i][j] = LUMBER
				} else {
					newfield[i][j] = field[i][j]
				}
			case LUMBER:
				if lumber >= 1 && trees >= 1 {
					newfield[i][j] = LUMBER
				} else {
					newfield[i][j] = OPEN
				}
			}
		}
	}

	return newfield
}

func printField(field [][]int) string {
	return printIntFieldWithSubs(field, "", map[int]string{
		OPEN:   ".",
		TREES:  "|",
		LUMBER: "#",
	})
}

func readSym(ch rune) int {
	switch ch {
	case '.':
		return OPEN
	case '|':
		return TREES
	case '#':
		return LUMBER
	default:
		panic("unknown symbol")
	}
}

func countItems(field [][]int) (int, int, int) {
	open, trees, lumber := 0, 0, 0
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			switch field[i][j] {
			case OPEN:
				open++
			case TREES:
				trees++
			case LUMBER:
				lumber++
			}
		}
	}
	return open, trees, lumber
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	field := makeIntField(len(lines), len(lines[0]))
	for i, s := range lines {
		for j, ch := range s {
			field[i][j] = readSym(ch)
		}
	}

	println(printField(field))

	memo := make(map[string]int)
	memo[printField(field)] = 1

	// 204369
	// 206624
	// 209451

	for i := 1; i <= 430+(1_000_000_000-430)%28; i++ {
		//for i := 1; i < 1_000_000_000; i++ {
		field = evolve(field)
		pf := printField(field)
		if v := memo[pf]; v == 1 {
			printf("cycle after %d minutes", i+1)
			memo[pf] = i
		} else if v > 1 {
			printf("cycle is %d", i-memo[pf])
			break
		} else {
			memo[pf]++
		}
		//printf("after %d minutes", i)
		//println(pf)
	}

	_, trees, lumber := countItems(field)
	printf("trees * lumber = %d", trees*lumber)
}
