package main

const (
	SERIAL = 1955
	//SERIAL = 42
)

func computePower(x, y int, serial int) int {
	rackId := x + 10
	pow := (rackId*y + serial) * rackId
	pow = (pow % 1000) / 100
	res := pow - 5
	return res
}

func findMax(field [][]int, size int) *Point3 {
	maxPow := -ALOT
	var p *Point3
	for i := 0; i < len(field)-size; i++ {
		for j := 0; j < len(field[0])-size; j++ {
			pow := 0
			for ii := 0; ii < size; ii++ {
				for jj := 0; jj < size; jj++ {
					pow += field[i+ii][j+jj]
				}
			}
			if pow > maxPow {
				maxPow = pow
				p = NewPoint3(j+1, i+1, maxPow)
			}
		}
	}
	printf("max pow: %d", maxPow)
	return p
}

func main() {
	field := makeIntField(300, 300)
	for i := 1; i <= len(field); i++ {
		for j := 1; j <= len(field[0]); j++ {
			field[i-1][j-1] = computePower(j, i, SERIAL)
		}
	}

	//printf("3,5,8: %d", computePower(3, 5, 8))
	//printf("122,79,57: %d", computePower(122, 79, 57))
	//printf("217,196,39: %d", computePower(217, 196, 39))
	//printf("101,153,71: %d", computePower(101, 153, 71))

	res1 := findMax(field, 3)
	printf("res1: %d,%d", res1.x, res1.y)

	maxPow := res1.z
	p := res1
	for s := 1; s < 300; s++ {
		tp := findMax(field, s)
		if tp.z > maxPow {
			maxPow = tp.z
			p = NewPoint3(tp.x, tp.y, s)
		}
	}
	printf("res2: %d,%d,%d", p.x, p.y, p.z)
}
