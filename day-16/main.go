package main

import (
	"fmt"
	"os"
)

type Instruction uint8

type Sample struct {
	instr, before, after [4]int
}

const (
	BEFORE = "Before"
	AFTER  = "After"
)

const (
	_ Instruction = iota
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
	Operations = map[Instruction]func(regs, instr [4]int) [4]int{
		ADDR: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = regs[instr[1]] + regs[instr[2]]
			return regs
		},
		ADDI: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = regs[instr[1]] + instr[2]
			return regs
		},
		MULR: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = regs[instr[1]] * regs[instr[2]]
			return regs
		},
		MULI: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = regs[instr[1]] * instr[2]
			return regs
		},
		BANR: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = int(uint32(regs[instr[1]]) & uint32(regs[instr[2]]))
			return regs
		},
		BANI: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = int(uint32(regs[instr[1]]) & uint32(instr[2]))
			return regs
		},
		BORR: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = int(uint32(regs[instr[1]]) | uint32(regs[instr[2]]))
			return regs
		},
		BORI: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = int(uint32(regs[instr[1]]) | uint32(instr[2]))
			return regs
		},
		SETR: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = regs[instr[1]]
			return regs
		},
		SETI: func(regs, instr [4]int) [4]int {
			regs[instr[3]] = instr[1]
			return regs
		},
		GTIR: func(regs, instr [4]int) [4]int {
			if instr[1] > regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
		GTRI: func(regs, instr [4]int) [4]int {
			if regs[instr[1]] > instr[2] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
		GTRR: func(regs, instr [4]int) [4]int {
			if regs[instr[1]] > regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
		EQIR: func(regs, instr [4]int) [4]int {
			if instr[1] == regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
		EQRI: func(regs, instr [4]int) [4]int {
			if regs[instr[1]] == instr[2] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
		EQRR: func(regs, instr [4]int) [4]int {
			if regs[instr[1]] == regs[instr[2]] {
				regs[instr[3]] = 1
			} else {
				regs[instr[3]] = 0
			}
			return regs
		},
	}
)

func NewSample(lines []string) *Sample {
	assert(lines[0][:6] == BEFORE, "Expect a sample to start with `Before`")
	before := parseInts(lines[0][9:19])
	instr := parseInts(lines[1])

	assert(lines[2][:5] == AFTER, "Expect a sample line to start with `After`")
	after := parseInts(lines[2][9:19])

	sample := &Sample{}
	copy(sample.before[:], before[:4])
	copy(sample.instr[:], instr[:4])
	copy(sample.after[:], after[:4])

	return sample
}

func (s *Sample) String() string {
	return fmt.Sprintf("sample{before: %+v, instr: %+v, after: %+v}", s.before, s.instr, s.after)
}

func part1(samples []*Sample) {
	cnt := 0
	for _, sample := range samples {
		instrs := make([]Instruction, 0, 1)
		for instr, apply := range Operations {
			if sample.after == apply(sample.before, sample.instr) {
				instrs = append(instrs, instr)
			}
		}
		if len(instrs) >= 3 {
			cnt++
		}
	}
	printf("part 1: %d", cnt)
}

func part2(samples []*Sample, commands [][4]int) {
	candidates := make(map[int]map[Instruction]struct{})
	for _, sample := range samples {
		for instr, apply := range Operations {
			if sample.after == apply(sample.before, sample.instr) {
				opcode := sample.instr[0]
				if _, ok := candidates[opcode]; !ok {
					candidates[opcode] = make(map[Instruction]struct{})
				}
				candidates[opcode][instr] = struct{}{}
			}
		}
	}
	printf("candidates: %+v", candidates)

	opcodes := make(map[int]Instruction)
	opcodesinv := make(map[Instruction]int)

	for len(candidates) > 0 {
		newcandidates := make(map[int]map[Instruction]struct{})
		for opcode, instrmap := range candidates {
			newinstrmap := make(map[Instruction]struct{})
			for instr := range instrmap {
				if _, ok := opcodesinv[instr]; ok {
					continue
				}
				newinstrmap[instr] = struct{}{}
			}
			if len(newinstrmap) == 1 {
				instr := firstKey(newinstrmap)
				opcodes[opcode] = instr
				opcodesinv[instr] = opcode
				continue
			}
			newcandidates[opcode] = newinstrmap
		}
		candidates = newcandidates
	}

	printf("opcodes: %+v", opcodes)

	var regs [4]int
	for _, command := range commands {
		instr := opcodes[command[0]]
		regs = Operations[instr](regs, command)
	}

	printf("part 2 regs: %+v", regs)
}

func firstKey(m map[Instruction]struct{}) Instruction {
	for instr := range m {
		return instr
	}
	return 0
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	samples := make([]*Sample, 0, 1)

	ix := 0
	for {
		if len(lines[ix]) == 0 {
			// end sample listing
			ix += 2
			break
		}
		sample := NewSample(lines[ix : ix+3])
		samples = append(samples, sample)
		ix += 4
	}

	commands := make([][4]int, 0, 1)

	for ix < len(lines) {
		var command [4]int
		copy(command[:], parseInts(lines[ix]))
		commands = append(commands, command)
		ix++
	}

	part1(samples)

	part2(samples, commands)
}
