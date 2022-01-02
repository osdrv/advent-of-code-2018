package main

import (
	"os"
)

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	nums := make([]int, 0, len(lines))
	for _, line := range lines {
		nums = append(nums, parseInt(line))
	}

	printf("nums: %+v", nums)

	sum := 0
	for _, num := range nums {
		sum += num
	}
	printf("sum is: %d", sum)

	sum = 0
	sums := make(map[int]int)
Cycle:
	for {
		for _, num := range nums {
			sum += num
			sums[sum]++
			if sums[sum] > 1 {
				printf("first repeating freq: %d", sum)
				break Cycle
			}
		}
	}

}
