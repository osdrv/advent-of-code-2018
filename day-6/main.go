package main

import (
	"log"
	"os"
)

func parsePoints(lines []string) []*Point2 {
	points := make([]*Point2, 0, len(lines))
	for _, line := range lines {
		ints := parseInts(line)
		points = append(points, NewPoint2(ints[0], ints[1]))
	}
	return points
}

func boundingBox(points []*Point2) (int, int, int, int) {
	minx, miny, maxx, maxy := points[0].x, points[0].y, points[0].x, points[0].y
	for _, point := range points {
		minx = min(minx, point.x)
		miny = min(miny, point.y)
		maxx = max(maxx, point.x)
		maxy = max(maxy, point.y)
	}
	return minx, miny, maxx, maxy
}

func floodFill(p0 *Point2, minx, miny, maxx, maxy int, points []*Point2) ([]Point2, bool) {
	var res []Point2
	ok := true
	var visit func(p *Point2) bool

	visited := make(map[Point2]bool)
	visit = func(p *Point2) bool {
		if p0.x == 5 && p0.y == 5 && p.x == 5 && p.y == 2 {
			//runtime.Breakpoint()
		}
		printf("visiting %s", p)
		if visited[*p] {
			return true
		}
		// ensure this point is closest to p0
		d0 := abs(p0.x-p.x) + abs(p0.y-p.y)
		for _, pp := range points {
			if *pp == *p0 {
				continue
			}
			d := abs(p.x-pp.x) + abs(p.y-pp.y)
			if d <= d0 {
				// this point is at worst as close to p0, return immediately
				return true
			}
		}
		// by now we're sure this is the closest position
		visited[*p] = true
		if p.x <= minx || p.y <= miny || p.x >= maxx || p.y >= maxy {
			ok = false
			return false
		}
		for _, step := range STEPS4 {
			nx, ny := p.x+step[0], p.y+step[1]
			if !visit(NewPoint2(nx, ny)) {
				return false
			}
		}
		return true
	}

	visit(p0)

	printf("visited: %+v", visited)

	if ok {
		res = make([]Point2, 0, len(visited))
		for pp := range visited {
			res = append(res, pp)
		}
	}
	return res, ok
}

const (
	MAXDIST = 10000
)

func dist(p1, p2 *Point2) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

func solve2(points []*Point2) []Point2 {
	minx, miny, maxx, maxy := boundingBox(points)
	var visit func(p *Point2)
	visited := make(map[Point2]bool)
	visit = func(p *Point2) {
		if visited[*p] {
			return
		}
		d := 0
		for _, pp := range points {
			d += dist(p, pp)
		}
		if d < MAXDIST {
			visited[*p] = true
			for _, step := range STEPS4 {
				nx, ny := p.x+step[0], p.y+step[1]
				visit(NewPoint2(nx, ny))
			}
		}
	}
	visit(NewPoint2((maxx-minx)/2, (maxy-miny)/2))

	res := make([]Point2, 0, len(visited))
	for vv := range visited {
		res = append(res, vv)
	}
	return res
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	points := parsePoints(lines)
	log.Printf("points: %+v", points)

	minx, miny, maxx, maxy := boundingBox(points)

	maxp := 0
	for _, point := range points {
		pp, ok := floodFill(point, minx, miny, maxx, maxy, points)
		if ok {
			if len(pp) > maxp {
				printf("new best area: %+v", pp)
				maxp = len(pp)
			}
		}
	}

	printf("the biggest area is: %d", maxp)

	res2 := solve2(points)
	printf("distance within range: %d", len(res2))
}
