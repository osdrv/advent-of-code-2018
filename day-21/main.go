package main

import (
	"fmt"
	"os"
	"runtime"
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

	BRKP
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
		BRKP: func(regs []int, instr [4]int) {
			runtime.Breakpoint()
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

	case "brkp":
		return BRKP
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

	case BRKP:
		return "brkp"
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

func eval(regs []int, instrs []Instr) int {
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
	cnt := 0
	maxb := -1
	last := -1
	seen := make(map[int]struct{})
	for pc < len(instrs) {

		i := instrs[pc]
		cnt++
		//printf("pc: %d, executing instr %s", pc, i)

		if pc == 28 {
			b := regs[1]
			if b > maxb {
				printf("new max b: %d", b)
				maxb = b
			}
			if _, ok := seen[regs[1]]; ok {
				printf("last: %d", last)
				break
			}
			seen[regs[1]] = struct{}{}
			last = regs[1]
		}

		binder(Operations[i[0]])(regs, i)
		pc++

		//printf("regs: %+v", regs)
		//if regs[0] != prev {
		//	printf("regs[0] changed from %d to %d", prev, regs[0])
		//	prev = regs[0]
		//}
	}

	return cnt
}

func main() {
	f, err := os.Open("INPUT-TST")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	instrs := make([]Instr, 0, len(lines))
	for _, s := range lines {
		instrs = append(instrs, parseInstr(s))
	}

	regs := make([]int, 6)
	regs[0] = 123
	cnt := eval(regs, instrs)
	printf("cnt: %d", cnt)

	//for _, a := range []int{2159153, 6413910, 7723681, 9861226, 16533497, 16573011, 16706847, 16768316, 16776078, 16776274} {
	//	regs := make([]int, 6)
	//	regs[0] = a
	//	cnt := eval(regs, instrs)

	//	printf("a: %d, instrs: %d", a, cnt)
	//}
}
