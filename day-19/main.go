package main

import (
	"fmt"
	"os"
)

type Instr [4]int

func (i Instr) String() string {
	return fmt.Sprintf("[%s %d %d %d]", opcodeToStr(i[0]), i[1], i[2], i[3])
}

const (
	BIND int = iota

	ADDR
	ADDI

	MULR
	MULI

	BANR
	BANI

	BORR
	BORI

	SETR
	SETI

	GTIR
	GTRI
	GTRR

	EQIR
	EQRI
	EQRR
)

var (
	Operations = map[int]func(regs []int, instr [4]int){
		ADDR: func(regs []int, instr [4]int) {
			regs[instr[3]] = regs[instr[1]] + regs[instr[2]]
		},
		ADDI: func(regs []int, instr [4]int) {
			regs[instr[3]] = regs[instr[1]] + instr[2]
		},
		MULR: func(regs []int, instr [4]int) {
			regs[instr[3]] = regs[instr[1]] * regs[instr[2]]
		},
		MULI: func(regs []int, instr [4]int) {
			regs[instr[3]] = regs[instr[1]] * instr[2]
		},
		BANR: func(regs []int, instr [4]int) {
			regs[instr[3]] = int(uint32(regs[instr[1]]) & uint32(regs[instr[2]]))
		},
		BANI: func(regs []int, instr [4]int) {
			regs[instr[3]] = int(uint32(regs[instr[1]]) & uint32(instr[2]))
		},
		BORR: func(regs []int, instr [4]int) {
			regs[instr[3]] = int(uint32(regs[instr[1]]) | uint32(regs[instr[2]]))
		},
		BORI: func(regs []int, instr [4]int) {
			regs[instr[3]] = int(uint32(regs[instr[1]]) | uint32(instr[2]))
		},
		SETR: func(regs []int, instr [4]int) {
			regs[instr[3]] = regs[instr[1]]
		},
		SETI: func(regs []int, instr [4]int) {
			regs[instr[3]] = instr[1]
		},
		GTIR: func(regs []int, instr [4]int) {
			if instr[1] > regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
		GTRI: func(regs []int, instr [4]int) {
			if regs[instr[1]] > instr[2] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
		GTRR: func(regs []int, instr [4]int) {
			if regs[instr[1]] > regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
		EQIR: func(regs []int, instr [4]int) {
			if instr[1] == regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
		EQRI: func(regs []int, instr [4]int) {
			if regs[instr[1]] == instr[2] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
		EQRR: func(regs []int, instr [4]int) {
			if regs[instr[1]] == regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
		},
	}
)

func parseOpCode(s string) int {
	switch s {
	case "addr":
		return ADDR
	case "addi":
		return ADDI

	case "mulr":
		return MULR
	case "muli":
		return MULI

	case "banr":
		return BANR
	case "bani":
		return BANI

	case "borr":
		return BORR
	case "bori":
		return BORI

	case "setr":
		return SETR
	case "seti":
		return SETI

	case "gtir":
		return GTIR
	case "gtri":
		return GTRI
	case "gtrr":
		return GTRR

	case "eqir":
		return EQIR
	case "eqri":
		return EQRI
	case "eqrr":
		return EQRR
	default:
		panic("unknown OpCode")
	}
}

func parseBind(s string) Instr {
	assert(s[0] == '#', "Bind instruction starts with '#'")
	return [4]int{BIND, parseInt(s[4:]), 0, 0}
}

func opcodeToStr(op int) string {
	switch op {
	case BIND:
		return "bind"
	case ADDR:
		return "addr"
	case ADDI:
		return "addi"
	case MULR:
		return "mulr"
	case MULI:
		return "muli"

	case BANR:
		return "banr"
	case BANI:
		return "bani"

	case BORR:
		return "borr"
	case BORI:
		return "bori"

	case SETR:
		return "setr"
	case SETI:
		return "seti"

	case GTIR:
		return "gtir"
	case GTRI:
		return "gtri"
	case GTRR:
		return "gtrr"

	case EQIR:
		return "eqir"
	case EQRI:
		return "eqri"
	case EQRR:
		return "eqrr"
	default:
		panic("oops")
	}
}

func parseInstr(s string) Instr {
	if s[0] == '#' {
		return parseBind(s)
	}
	opcode := parseOpCode(s[:4])
	argv := parseInts(s[5:])
	return [4]int{opcode, argv[0], argv[1], argv[2]}
}

func eval(regs []int, instrs []Instr) {
	binder := func(fn func(regs []int, instr [4]int)) func(regs []int, instr [4]int) {
		return fn
	}
	pc := 0
	if instrs[0][0] == BIND {
		i := instrs[0]
		var bb *int = &regs[i[1]]
		binder = func(fn func(regs []int, instr [4]int)) func(regs []int, instr [4]int) {
			return func(regs []int, instr [4]int) {
				*bb = pc
				fn(regs, instr)
				pc = *bb
			}
		}
		instrs = instrs[1:]
	}
	//prev := regs[0]
	for pc < len(instrs) {
		i := instrs[pc]
		//printf("pc: %d, executing instr %s", pc, i)
		binder(Operations[i[0]])(regs, i)
		pc++

		//printf("regs: %+v", regs)
		//if regs[0] != prev {
		//	printf("regs[0] changed from %d to %d", prev, regs[0])
		//	prev = regs[0]
		//}
	}
}

func prog1() {
	a, b, c, d, e, f := 0, 0, 0, 0, 0, 0
	c = 2 * 2 * 19 * 11
	d = (6 * 22) + 9
	c += d

	printf("c=%d", c)

	for e = 1; e <= c; e++ {
		for b = 1; b <= c; b++ {
			if e*b == c {
				printf("found factors: %d and %d", e, b)
				a += e
			}
		}
	}
	printf("a: %d, b: %d, c: %d, d: %d, e: %d, f: %d", a, b, c, d, e, f)
}

func prog2() {
	a, b, c, d, e, f := 0, 0, 0, 0, 0, 0

	c = 2 * 2 * 19 * 11
	d = (6 * 22) + 9
	c += d

	d = 27
	d *= 28
	d += 29
	d *= 30
	d *= 14
	d *= 32
	c += d

	printf("c=%d (expect: 10551377)", c)

	for e = 1; e <= c; e++ {
		if c%e == 0 {
			a += e
		}
	}

	//for e = 1; e <= c; e++ {
	//	for b = 1; b <= c; b++ {
	//		if e*b == c {
	//			printf("found factors: %d and %d", e, b)
	//			a += e
	//		}
	//	}
	//}

	printf("a: %d, b: %d, c: %d, d: %d, e: %d, f: %d", a, b, c, d, e, f)
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	instrs := make([]Instr, 0, len(lines))
	for _, s := range lines {
		instrs = append(instrs, parseInstr(s))
	}

	printf("instrs: %+v", instrs)

	regs := make([]int, 6)
	//regs[0] = 1
	eval(regs, instrs)

	printf("regs: %+v", regs)

	prog2()
}
