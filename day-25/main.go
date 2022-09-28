package main

import (
	"os"
)

func parseLine(s string) [4]int {
	nums := parseInts(s)
	return [4]int{nums[0], nums[1], nums[2], nums[3]}
}

func distance(p1, p2 [4]int) int {
	return abs(p1[0]-p2[0]) + abs(p1[1]-p2[1]) + abs(p1[2]-p2[2]) + abs(p1[3]-p2[3])
}

func computeDjs(points [][4]int, dist int) []int {
	edges := make(map[int][]int)
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			if distance(points[i], points[j]) <= dist {
				if _, ok := edges[i]; !ok {
					edges[i] = make([]int, 0, 1)
				}
				if _, ok := edges[j]; ok {
					edges[j] = make([]int, 0, 1)
				}
				edges[i] = append(edges[i], j)
				edges[j] = append(edges[j], i)
			}
		}
	}

	djs := make([]int, len(points))
	for i := 0; i < len(djs); i++ {
		djs[i] = i
	}
	getParent := func(i int) int {
		parent := i
		for djs[parent] != parent {
			parent = djs[parent]
		}
		return parent
	}

	for i := 0; i < len(points); i++ {
		for _, edge := range edges[i] {
			edgeparent := getParent(edge)
			parent := getParent(i)
			if edgeparent != parent {
				djs[edgeparent] = parent
			}
		}
	}

	return djs
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	points := make([][4]int, 0, len(lines))

	for _, line := range lines {
		points = append(points, parseLine(line))
	}

	djs := computeDjs(points, 3)

	cnt := 0
	for edge, parent := range djs {
		if edge == parent {
			cnt++
		}
	}

	printf("djs size: %d", cnt)
}
