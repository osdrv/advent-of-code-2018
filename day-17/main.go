package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func parseRangePair(s string) (int, int) {
	ss := strings.Split(s, "..")
	v0 := parseInt(ss[0])
	v1 := v0
	if len(ss) > 1 {
		v1 = parseInt(ss[1])
	}
	return v0, v1
}

func parseRange(s string) [2]Point2 {
	ss := strings.SplitN(s, ", ", 2)
	if ss[0][0] == 'y' {
		ss[0], ss[1] = ss[1], ss[0]
	}
	x0, x1 := parseRangePair(ss[0][2:])
	y0, y1 := parseRangePair(ss[1][2:])
	return [2]Point2{{x0, y0}, {x1, y1}}
}

const (
	SAND int = iota
	CLAY
	WATER
	STREAM
	SOURCE
)

func makeField(rr [][2]Point2) map[Point2]int {
	field := make(map[Point2]int)
	for _, r := range rr {
		for y := min(r[0].y, r[1].y); y <= max(r[0].y, r[1].y); y++ {
			for x := min(r[0].x, r[1].x); x <= max(r[0].x, r[1].x); x++ {
				p := Point2{x, y}
				field[p] = CLAY
			}
		}
	}
	return field
}

func fieldDim(field map[Point2]int) (int, int, int, int) {
	xmin, ymin, xmax, ymax := ALOT, ALOT, -ALOT, -ALOT
	for p := range field {
		xmin = min(xmin, p.x)
		ymin = min(ymin, p.y)
		xmax = max(xmax, p.x)
		ymax = max(ymax, p.y)
	}
	return xmin, ymin, xmax, ymax
}

func printField(field map[Point2]int) string {
	var b bytes.Buffer
	xmin, ymin, xmax, ymax := fieldDim(field)
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			p := Point2{x, y}
			switch field[p] {
			case CLAY:
				b.WriteByte('#')
			case WATER:
				b.WriteByte('~')
			case STREAM:
				b.WriteByte('|')
			case SOURCE:
				b.WriteByte('+')
			default:
				b.WriteByte('.')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

const (
	INF int = iota
	EDGE
	DROP
)

func flood(field map[Point2]int) {
	xmin, ymin, xmax, ymax := fieldDim(field)

	findSideBound := func(p Point2, step int) (int, int) {
		x := p.x
		for x >= xmin-1 && x <= xmax+1 {
			below := field[Point2{x, p.y + 1}]
			if below == CLAY || below == WATER {
				if field[Point2{x + step, p.y}] == CLAY {
					return x, EDGE
				}
				x += step
			} else if below == STREAM || below == SAND {
				return x, DROP
			}
		}
		return x, INF
	}

	q := make([]Point2, 0, 1)
	enq := make(map[Point2]bool)

	enqueue := func(p Point2) {
		if enq[p] {
			printf("point %v has already been enqueued, skipping", p)
			return
		}
		enq[p] = true
		q = append(q, p)
	}

	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			p := Point2{x, y}
			if field[p] == SOURCE {
				enqueue(p)
				break
			}
		}
	}

	var head Point2
	for len(q) > 0 {
		head, q = q[0], q[1:]
		for field[Point2{head.x, head.y + 1}] == SAND && head.y < ymax {
			head.y++
			field[head] = STREAM
		}
		if head.y >= ymax {
			continue
		}
	FloodLine:
		left, nml := findSideBound(head, -1)
		right, nmr := findSideBound(head, 1)

		if DEBUG {
			printf("left: %+v, nml: %d, right: %+v, nmr: %d", left, nml, right, nmr)
		}

		if nml == INF && nmr == INF {
			continue
		}

		if nml == EDGE && nmr == EDGE {
			for x := left; x <= right; x++ {
				field[Point2{x, head.y}] = WATER
			}
			head.y--
			goto FloodLine
		}

		if nml == DROP {
			enqueue(Point2{left, head.y})
		}
		if nmr == DROP {
			enqueue(Point2{right, head.y})
		}
		for x := left; x <= right; x++ {
			field[Point2{x, head.y}] = STREAM
		}
		if DEBUG {
			println(printField(field))
			var s string
			fmt.Scanf("%s", &s)
		}
	}
}

func countAllWater(f map[Point2]int) int {
	cnt := 0
	for _, v := range f {
		if v == WATER || v == STREAM {
			cnt++
		}
	}
	return cnt
}

func countRetWater(f map[Point2]int) int {
	cnt := 0
	for _, v := range f {
		if v == WATER {
			cnt++
		}
	}
	return cnt
}

const DEBUG = false

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	ranges := make([][2]Point2, 0, len(lines))

	for _, s := range lines {
		ranges = append(ranges, parseRange(s))
	}

	printf("ranges: %+v", ranges)

	field := makeField(ranges)
	field[Point2{500, 0}] = SOURCE

	flood(field)

	println(printField(field))

	allWater := countAllWater(field)
	printf("all water tiles: %d", allWater)

	retWater := countRetWater(field)
	printf("retain water: %d", retWater)
}
