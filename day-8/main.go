package main

import (
	"os"
)

type Node struct {
	children []*Node
	meta     []int
}

func (n *Node) Meta() []int {
	res := make([]int, 0, 1)
	res = append(res, n.meta...)
	for _, ch := range n.children {
		res = append(res, ch.Meta()...)
	}
	return res
}

func (n *Node) Value() int {
	res := 0
	if len(n.children) == 0 {
		for _, num := range n.meta {
			res += num
		}
	} else {
		for _, ix := range n.meta {
			if ix > 0 && ix <= len(n.children) {
				res += n.children[ix-1].Value()
			}
		}
	}
	return res
}

func solve1(nums []int) int {
	root, _ := parseNode(nums, 0)
	meta := root.Meta()
	sum := 0
	for _, num := range meta {
		sum += num
	}
	//TODO
	return sum
}

func solve2(nums []int) int {
	root, _ := parseNode(nums, 0)
	return root.Value()
}

func parseNode(nums []int, ptr int) (*Node, int) {
	if ptr >= len(nums) {
		return nil, ptr
	}
	nch := nums[ptr]
	ptr++
	nmeta := nums[ptr]
	ptr++
	children := make([]*Node, 0, nch)
	var ch *Node
	for i := 0; i < nch; i++ {
		ch, ptr = parseNode(nums, ptr)
		children = append(children, ch)
	}
	meta := nums[ptr : ptr+nmeta]
	ptr += nmeta
	return &Node{
		children: children,
		meta:     meta,
	}, ptr
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	data := readFile(f)
	nums := parseInts(data)

	res1 := solve1(nums)
	printf("part 1 answer is: %d", res1)

	res2 := solve2(nums)
	printf("part 2 answer is: %d", res2)
}
