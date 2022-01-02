package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Shift struct {
	gid    int
	sleeps [][2]int
}

func NewShift(gid int, sleeps [][2]int) *Shift {
	return &Shift{
		gid:    gid,
		sleeps: sleeps,
	}
}

func (s *Shift) SleepTime() int {
	time := 0
	for _, sleep := range s.sleeps {
		time += sleep[1] - sleep[0]
	}
	return time
}

func (s *Shift) String() string {
	return fmt.Sprintf("shift{gid: %d, sleeps: %+v}", s.gid, s.sleeps)
}

func parseShifts(lines []string) []*Shift {
	shifts := make([]*Shift, 0, 1)
	ix := 0
	var gid, start, end int
	sleeps := make([][2]int, 0, 1)
	for ix < len(lines) {
		line := lines[ix]
		hh, mm := parseInt(line[12:14]), parseInt(line[15:17])
		if hh == 23 {
			mm = mm - 60
		}
		line = line[19:]
		printf("line: %q", line)
		if "falls asleep" == line {
			start = mm
		} else if "wakes up" == line {
			end = mm - 1
			sleeps = append(sleeps, [2]int{start, end})
		} else {
			if ix > 0 {
				shifts = append(shifts, NewShift(gid, sleeps))
				sleeps = make([][2]int, 0, 1)
			}
			chs := strings.FieldsFunc(line, func(r rune) bool {
				return r == ' ' || r == '#'
			})
			gid = parseInt(chs[1])
		}
		ix++
	}
	shifts = append(shifts, NewShift(gid, sleeps))
	return shifts
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	log.Printf("file data: %+v", lines)

	shifts := parseShifts(lines)
	printf("shifts: %+v", shifts)

	sleeps := make(map[int]int)
	maxsleep := 0
	maxgid := 0
	for _, shift := range shifts {
		sleeps[shift.gid] += shift.SleepTime()
		if sleeps[shift.gid] > maxsleep {
			maxsleep = sleeps[shift.gid]
			maxgid = shift.gid
		}
	}

	var minutes [60]int
	var maxix, maxminutes int
	for _, shift := range shifts {
		if shift.gid != maxgid {
			continue
		}
		for _, sleep := range shift.sleeps {
			for i := sleep[0]; i <= sleep[1]; i++ {
				minutes[i]++
				if minutes[i] > maxminutes {
					maxminutes = minutes[i]
					maxix = i
				}
			}
		}
	}

	printf("%d * %d = %d", maxgid, maxix, maxgid*maxix)

	maxgid, maxsleep, maxmin := 0, 0, 0

	perGuard := make(map[int][60]int)
	for _, shift := range shifts {
		mm := perGuard[shift.gid]
		for _, sleep := range shift.sleeps {
			for i := sleep[0]; i <= sleep[1]; i++ {
				mm[i]++
				if mm[i] > maxsleep {
					maxsleep = mm[i]
					maxgid = shift.gid
					maxmin = i
				}
			}
		}
		perGuard[shift.gid] = mm
	}

	printf("%+v", perGuard)

	printf("strategy 2")
	printf("guard id: %d, maxmin: %d, maxsleep: %d", maxgid, maxmin, maxsleep)
	printf("the answer is: %d", maxgid*maxmin)

}
