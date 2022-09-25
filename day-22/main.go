package main

import "bytes"

const (
	ROCKY int = iota
	WET
	NARROW
)

const (
	MOD = 20183
)

const (
	NEITHER int = 1 << iota
	CLIMB
	TORCH
)

func GenGraph(depth int, target Point2) func(int, int) int {

	var geoIxAt func(x, y int) int
	var errAt func(x, y int) int

	errAt = func(x, y int) int {
		return (geoIxAt(x, y)%MOD + depth%MOD) % MOD
	}

	gix := make(map[Point2]int)
	geoIxAt = func(x, y int) int {
		p := Point2{x, y}
		if v, ok := gix[p]; ok {
			return v
		}
		var g int

		if x == 0 && y == 0 {
			g = 0
		} else if x == target.x && y == target.y {
			g = 0
		} else if y == 0 {
			g = x * 16807
		} else if x == 0 {
			g = y * 48271
		} else {
			g = (errAt(x-1, y) * errAt(x, y-1)) % MOD
		}

		gix[p] = g
		return g
	}

	return func(x, y int) int {
		return errAt(x, y) % 3
	}
}

const (
	SWAP_TIME = 7
	MOVE_TIME = 1
)

func main() {
	depth := 5355
	target := Point2{14, 796}
	//depth := 510
	//target := Point2{10, 10}

	graph := GenGraph(depth, target)

	var EQUIP [3]int
	EQUIP[ROCKY] = CLIMB | TORCH
	EQUIP[WET] = CLIMB | NEITHER
	EQUIP[NARROW] = TORCH | NEITHER

	risk := 0
	var b bytes.Buffer
	for y := 0; y < target.y+1; y++ {
		for x := 0; x < target.x+1; x++ {
			t := graph(x, y)
			risk += t
			if t == ROCKY {
				b.WriteByte('.')
			} else if t == WET {
				b.WriteByte('=')
			} else if t == NARROW {
				b.WriteByte('|')
			}
		}
		b.WriteByte('\n')
	}
	printf("total risk: %d", risk)
	println(b.String())

	start := Point3{0, 0, TORCH}
	finish := Point3{target.x, target.y, TORCH}

	gScore := make(map[Point3]int)
	fScore := make(map[Point3]int)

	h := func(p Point3) int {
		d := abs(p.x-finish.x) + abs(p.y-finish.y)
		if p.z != finish.z {
			d += SWAP_TIME
		}
		return d
	}

	gScore[start] = 0
	fScore[start] = h(start)

	q := NewBinHeap(func(a, b Point3) bool {
		return fScore[a] < fScore[b]
	})
	q.Push(start)

	for q.Size() > 0 {
		curr := q.Pop()
		if curr == finish {
			printf("Min time: %d", gScore[curr])
			break
		}

		currtyp := graph(curr.x, curr.y)

		for _, step := range STEPS4 {
			nx, ny := curr.x+step[0], curr.y+step[1]
			if nx < 0 || ny < 0 {
				continue
			}
			ntyp := graph(nx, ny)
			for i := 0; i < 3; i++ {
				neq := 1 << i
				if neq&EQUIP[currtyp]&EQUIP[ntyp] == 0 {
					continue
				}
				np := Point3{nx, ny, neq}
				dt := MOVE_TIME
				if neq != curr.z {
					dt += SWAP_TIME
				}
				if _, ok := gScore[np]; !ok {
					gScore[np] = ALOT
				}
				ngScore := gScore[curr] + dt
				if ngScore < gScore[np] {
					gScore[np] = ngScore
					fScore[np] = ngScore + h(np)
					q.Push(Point3{nx, ny, neq})
				}
			}
		}
	}
}
