package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

func parseConn(s string) [2]byte {
	return [2]byte{s[5], s[36]}
}

func solve1(conns [][2]byte) string {
	unvisited := make(map[byte]bool)
	visited := make(map[byte]bool)
	deps := make(map[byte][]byte)
	q := make([]byte, 0, 1)
	enq := make(map[byte]bool)
	for _, conn := range conns {
		to, from := conn[0], conn[1]
		unvisited[from] = true
		unvisited[to] = true
		deps[from] = append(deps[from], to)
	}
	var buf bytes.Buffer
	for len(unvisited) > 0 {
	Unvisited:
		for v := range unvisited {
			if enq[v] {
				continue
			}
			for _, dep := range deps[v] {
				if _, ok := visited[dep]; !ok {
					continue Unvisited
				}
			}
			q = append(q, v)
			enq[v] = true
		}
		sort.Slice(q, func(i, j int) bool {
			return q[i] < q[j]
		})
		if len(q) == 0 {
			break
		}
		var head byte
		head, q = q[0], q[1:]
		visited[head] = true
		delete(unvisited, head)
		buf.WriteByte(head)
	}
	return buf.String()
}

func printWorkers(time int, workers [][2]int, done []byte) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%04d", time))
	for i := 0; i < len(workers); i++ {
		buf.WriteString(fmt.Sprintf("\t%c", byte(workers[i][0])))
	}
	buf.WriteString("\t" + string(done))
	buf.WriteByte('\n')
	return buf.String()
}

func solve2(conns [][2]byte, nWorkers int, stepDur int) int {
	workers := make([][2]int, nWorkers)
	deps := make(map[byte][]byte)
	visited := make(map[byte]bool)
	enq := make(map[byte]bool)
	unvisited := make(map[byte]bool)
	q := make([]byte, 0, 1)

	for _, conn := range conns {
		to, from := conn[0], conn[1]
		unvisited[from] = true
		unvisited[to] = true
		deps[from] = append(deps[from], to)
	}

	getJob := func() (byte, bool) {
		if len(unvisited) == 0 {
			return 0, false
		}
	Next:
		for v := range unvisited {
			if enq[v] {
				continue
			}
			for _, dep := range deps[v] {
				if _, ok := visited[dep]; !ok {
					continue Next
				}
			}
			q = append(q, v)
			enq[v] = true
		}
		if len(q) == 0 {
			return 0, false
		}
		sort.Slice(q, func(i, j int) bool { return q[i] < q[j] })
		var head byte
		head, q = q[0], q[1:]
		return head, true
	}

	time := 0
	done := make([]byte, 0, 1)
	for {
		active := false
		for i := 0; i < nWorkers; i++ {
			if workers[i][0] > 0 && workers[i][1] == 0 {
				b := byte(workers[i][0])
				visited[b] = true
				delete(unvisited, b)
				done = append(done, b)
				workers[i][0] = 0
			}
		}
		for i := 0; i < nWorkers; i++ {
			if workers[i][1] == 0 {
				if job, ok := getJob(); ok {
					workers[i][0] = int(job)
					workers[i][1] = stepDur + int(job-'A')
					active = true
				}
			} else {
				workers[i][1]--
				active = true
			}
		}
		print(printWorkers(time, workers, done))
		if len(unvisited) == 0 {
			break
		}
		if !active {
			panic("all workers are stuck")
		}
		time++
	}
	return time
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	conns := make([][2]byte, 0, len(lines))
	for _, line := range lines {
		conns = append(conns, parseConn(line))
	}

	res1 := solve1(conns)
	printf("part1 result: %s", res1)

	printf("conns: %+v", conns)

	res2 := solve2(conns, 5, 60)
	printf("time it takes: %d", res2)
}
