package main

import (
	"bytes"
	"os"
)

type Path interface {
	Traverse(*Point2, map[Point2]int)
	String() string
}

type OrPath struct {
	paths []Path
}

var _ Path = (*OrPath)(nil)

func NewOrPath(paths []Path) *OrPath {
	return &OrPath{paths: paths}
}

func (p *OrPath) Traverse(cur *Point2, m map[Point2]int) {
	for _, path := range p.paths {
		curcp := NewPoint2(cur.x, cur.y)
		path.Traverse(curcp, m)
	}
}

func (p *OrPath) String() string {
	var b bytes.Buffer
	b.WriteString("(OR: ")
	for ix, path := range p.paths {
		if ix > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(path.String())
	}
	b.WriteByte(')')
	return b.String()
}

type AndPath struct {
	paths []Path
}

var _ Path = (*AndPath)(nil)

func NewAndPath(paths []Path) *AndPath {
	return &AndPath{paths: paths}
}

func (p *AndPath) Traverse(cur *Point2, m map[Point2]int) {
	for _, path := range p.paths {
		path.Traverse(cur, m)
	}
}

func (p *AndPath) String() string {
	var b bytes.Buffer
	b.WriteString("(AND: ")
	for ix, path := range p.paths {
		if ix > 0 {
			b.WriteString(" & ")
		}
		b.WriteString(path.String())
	}
	b.WriteByte(')')
	return b.String()
}

type DirPath struct {
	dirs string
}

var _ Path = (*DirPath)(nil)

func NewDirPath(dirs string) *DirPath {
	return &DirPath{dirs: dirs}
}

func (p *DirPath) Traverse(cur *Point2, m map[Point2]int) {
	lookAround(cur, m)
	ptr := 0
	printf("traverse dir: %s from point: %s", p.dirs, cur)
	for ptr < len(p.dirs) {
		switch p.dirs[ptr] {
		case 'N':
			cur.y--
			m[*cur] = DOOR_H
			cur.y--
		case 'E':
			cur.x++
			m[*cur] = DOOR_V
			cur.x++
		case 'S':
			cur.y++
			m[*cur] = DOOR_H
			cur.y++
		case 'W':
			cur.x--
			m[*cur] = DOOR_V
			cur.x--
		default:
			panic("oopsie doopsie")
		}
		m[*cur] = ROOM
		lookAround(cur, m)
		ptr++
	}
}

func (p *DirPath) String() string {
	return p.dirs
}

func match(s string, ptr int, b byte) bool {
	return s[ptr] == b
}

func consume(s string, ptr int, b byte) int {
	if s[ptr] != b {
		panic("unexpected rune")
	}
	return ptr + 1
}

func isDir(s string, ptr int) bool {
	b := s[ptr]
	return b == 'N' || b == 'E' || b == 'S' || b == 'W'
}

func parseDirPath(s string, ptr int) (Path, int) {
	from := ptr
	for ptr < len(s) && isDir(s, ptr) {
		ptr++
	}
	return NewDirPath(s[from:ptr]), ptr
}

func parseAndPath(s string, ptr int) (Path, int) {
	paths := make([]Path, 0, 1)
	var p Path
	for ptr < len(s) {
		if match(s, ptr, '(') {
			ptr = consume(s, ptr, '(')
			p, ptr = parseOrPath(s, ptr)
			ptr = consume(s, ptr, ')')
		} else if isDir(s, ptr) {
			p, ptr = parseDirPath(s, ptr)
		} else {
			break
		}
		paths = append(paths, p)
	}
	if len(paths) == 1 {
		return paths[0], ptr
	}
	return NewAndPath(paths), ptr
}

func parseOrPath(s string, ptr int) (Path, int) {
	paths := make([]Path, 0, 1)
	var p Path
	for ptr < len(s) {
		p, ptr = parseAndPath(s, ptr)
		paths = append(paths, p)
		if !match(s, ptr, '|') {
			break
		}
		ptr = consume(s, ptr, '|')
	}
	if len(paths) == 1 {
		return paths[0], ptr
	}
	return NewOrPath(paths), ptr
}

func parsePath(s string) Path {
	var p Path
	var ptr int = 0
	ptr = consume(s, ptr, '^')
	p, ptr = parseOrPath(s, ptr)
	consume(s, ptr, '$')
	return p
}

const (
	UNKNWN int = iota
	WALL
	DOOR_H
	DOOR_V
	ROOM
	CURR
)

func lookAround(p *Point2, m map[Point2]int) {
	for _, step := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		np := Point2{p.x + step[0], p.y + step[1]}
		if _, ok := m[np]; !ok {
			m[np] = UNKNWN
		}
	}

	for _, step := range [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
		np := Point2{p.x + step[0], p.y + step[1]}
		if _, ok := m[np]; !ok {
			m[np] = WALL
		}
	}
}

func refineMap(m map[Point2]int) {
	for p, v := range m {
		if v == UNKNWN {
			m[p] = WALL
		}
	}
}

func buildMap(path Path) map[Point2]int {
	cur := NewPoint2(0, 0)
	m := make(map[Point2]int)
	m[*cur] = CURR
	path.Traverse(cur, m)

	refineMap(m)

	return m
}

func printMap(m map[Point2]int) string {
	xmin, ymin, xmax, ymax := ALOT, ALOT, -ALOT, -ALOT
	for p := range m {
		xmin = min(xmin, p.x)
		ymin = min(ymin, p.y)
		xmax = max(xmax, p.x)
		ymax = max(ymax, p.y)
	}

	var b bytes.Buffer
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			switch m[Point2{x, y}] {
			case UNKNWN:
				b.WriteByte('?')
			case WALL:
				b.WriteByte('#')
			case DOOR_H:
				b.WriteByte('-')
			case DOOR_V:
				b.WriteByte('|')
			case ROOM:
				b.WriteByte('.')
			case CURR:
				b.WriteByte('X')
			default:
				panic("oops")
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func visitMap(cur *Point2, m map[Point2]int) map[Point2]int {
	dd := make(map[Point2]int)

	q := make([]Point3, 0, 1)
	q = append(q, Point3{cur.x, cur.y, 0}) // the z-component keeps the current number of steps

	var head Point3
	for len(q) > 0 {
		head, q = q[0], q[1:]
		p := Point2{head.x, head.y}
		if v, ok := dd[p]; ok {
			if v <= head.z {
				continue
			}
		}
		dd[p] = head.z
		for _, step := range STEPS4 {
			np := Point2{p.x + step[0], p.y + step[1]}
			if m[np] == DOOR_H || m[np] == DOOR_V {
				q = append(q, Point3{np.x + step[0], np.y + step[1], head.z + 1})
			}
		}
	}

	return dd
}

func maxDoors(m map[Point2]int) int {
	md := 0
	for _, v := range m {
		md = max(md, v)
	}
	return md
}

func cntRoomsWithin(m map[Point2]int, dist int) int {
	cnt := 0
	for _, v := range m {
		if v >= dist {
			cnt++
		}
	}
	return cnt
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	for _, s := range lines {
		path := parsePath(s)
		m := buildMap(path)
		printf("input: %s", s)
		printf("path: %s", path.String())
		println(printMap(m))

		dd := visitMap(NewPoint2(0, 0), m)
		printf("distances: %+v", dd)
		printf("max number of doors to enter: %d", maxDoors(dd))
		printf("number of rooms with at least 1000 doors: %d", cntRoomsWithin(dd, 1000))
	}
}
