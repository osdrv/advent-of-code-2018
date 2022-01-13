package main

import (
	"bytes"
	"os"
	"sort"
)

type Path func(Cart) Cart

func NewPath(kind byte) Path {
	switch kind {
	case '-', '|':
		return STRAIGHT(kind)
	case '/', '\\':
		return CURVE(kind)
	case '+':
		return CROSS(kind)
	}
	panic("wtf")
}

var (
	STRAIGHT = func(kind byte) Path {
		return func(c Cart) Cart {
			return Cart{
				id:   c.id,
				pos:  c.pos,
				dir:  c.dir,
				turn: c.turn,
			}
		}
	}
	CURVE = func(kind byte) Path {
		return func(c Cart) Cart {
			newcart := Cart{
				id:   c.id,
				pos:  c.pos,
				turn: c.turn,
			}
			var dir int
			switch kind {
			case '/':
				dir = map[int]int{
					RIGHT: UP,
					UP:    RIGHT,
					LEFT:  DOWN,
					DOWN:  LEFT,
				}[c.dir]
			case '\\':
				dir = map[int]int{
					UP:    LEFT,
					LEFT:  UP,
					RIGHT: DOWN,
					DOWN:  RIGHT,
				}[c.dir]
			}
			newcart.dir = dir
			return newcart
		}
	}
	CROSS = func(ch byte) Path {
		return func(c Cart) Cart {
			newcart := Cart{
				id:  c.id,
				pos: c.pos,
			}

			var dir, turn int
			switch c.dir {
			case LEFT:
				switch c.turn {
				case LEFT:
					dir, turn = DOWN, UP
				case UP:
					dir, turn = LEFT, RIGHT
				case RIGHT:
					dir, turn = UP, LEFT
				}
			case UP:
				switch c.turn {
				case LEFT:
					dir, turn = LEFT, UP
				case UP:
					dir, turn = UP, RIGHT
				case RIGHT:
					dir, turn = RIGHT, LEFT
				}
			case RIGHT:
				switch c.turn {
				case LEFT:
					dir, turn = UP, UP
				case UP:
					dir, turn = RIGHT, RIGHT
				case RIGHT:
					dir, turn = DOWN, LEFT
				}
			case DOWN:
				switch c.turn {
				case LEFT:
					dir, turn = RIGHT, UP
				case UP:
					dir, turn = DOWN, RIGHT
				case RIGHT:
					dir, turn = LEFT, LEFT
				}
			}
			newcart.dir = dir
			newcart.turn = turn

			return newcart
		}
	}
)

const (
	_ int = iota
	UP
	RIGHT
	DOWN
	LEFT
)

var (
	DIR = map[byte]int{
		'<': LEFT,
		'>': RIGHT,
		'^': UP,
		'v': DOWN,
	}
	DIR_R = map[int]byte{
		LEFT:  '<',
		RIGHT: '>',
		UP:    '^',
		DOWN:  'v',
	}
)

type Cart struct {
	id   int
	pos  Point2
	dir  int
	turn int
}

func readGraph(lines []string) (map[Point2]Path, []Cart) {
	graph := make(map[Point2]Path)
	carts := make([]Cart, 0, 1)
	cartid := 1
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[0]); x++ {
			switch ch := lines[y][x]; ch {
			case '|', '-', '/', '\\', '+':
				graph[Point2{x, y}] = NewPath(ch)
			case '<', '^', 'v', '>':
				graph[Point2{x, y}] = NewPath('-') // assume a cart will never start off on a cross or a curve
				carts = append(carts, Cart{
					id:   cartid,
					pos:  Point2{x, y},
					dir:  DIR[ch],
					turn: LEFT,
				})
				cartid++
			}
		}
	}
	return graph, carts
}

func tick(graph map[Point2]Path, carts []Cart) []Cart {
	nextcarts := make([]Cart, 0, len(carts))
	coords := make(map[Point2]Cart)
	for _, cart := range carts {
		coords[cart.pos] = cart
	}
	hit := make(map[int]bool)
	sort.Slice(carts, func(i, j int) bool {
		return carts[i].pos.y < carts[j].pos.y || (carts[i].pos.y == carts[j].pos.y && carts[i].pos.x < carts[j].pos.x)
	})
	for _, cart := range carts {
		if hit[cart.id] {
			continue
		}
		x, y := cart.pos.x, cart.pos.y
		nx, ny := x, y
		switch cart.dir {
		case LEFT:
			nx--
		case RIGHT:
			nx++
		case UP:
			ny--
		case DOWN:
			ny++
		}

		np := Point2{nx, ny}
		assert(graph[np] != nil, "graph node is missing")

		if another, ok := coords[np]; ok {
			// we hit another cart
			hit[cart.id] = true
			hit[another.id] = true
			printf("collision at %d,%d", nx, ny)
			delete(coords, cart.pos)
			delete(coords, np)
			continue
		}

		nextcart := graph[np](cart)
		nextcart.pos = Point2{nx, ny}

		delete(coords, cart.pos)
		coords[np] = nextcart

		nextcarts = append(nextcarts, nextcart)
	}

	res := make([]Cart, 0, len(nextcarts))
	for _, cart := range nextcarts {
		if !hit[cart.id] {
			res = append(res, cart)
		}
	}

	return res
}

func printGraph(graph map[Point2]Path, carts []Cart) string {
	var buf bytes.Buffer

	var maxx, maxy int
	for p := range graph {
		maxx = max(maxx, p.x)
		maxy = max(maxy, p.y)
	}

	cmap := make(map[Point2]Cart)
	for _, cart := range carts {
		cmap[cart.pos] = cart
	}

	for y := 0; y <= maxy; y++ {
		for x := 0; x <= maxx; x++ {
			p := Point2{x, y}
			if cart, ok := cmap[p]; ok {
				buf.WriteByte(DIR_R[cart.dir])
			} else if _, ok := graph[p]; ok {
				buf.WriteByte('.')
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}

func part1(graph map[Point2]Path, carts []Cart) {
TICK:
	for {
		carts = tick(graph, carts)
		coords := make(map[Point2]bool)
		for _, cart := range carts {
			p := cart.pos
			if _, ok := coords[p]; ok {
				printf("first collision at: %d,%d", p.x, p.y)
				break TICK
			}
			coords[p] = true
		}
	}
}

func part2(graph map[Point2]Path, carts0 []Cart) {
	carts := carts0
	for {
		if len(carts) < 1 {
			panic("wtf")
		}
		carts = tick(graph, carts)
		if len(carts) == 1 {
			printf("last cart is at: %d,%d", carts[0].pos.x, carts[0].pos.y)
			break
		}
	}
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	graph, carts := readGraph(lines)

	//printf("graph: %+v", graph)
	//printf("carts: %+v", carts)

	//part1(graph, cartscp)

	part2(graph, carts)
}
