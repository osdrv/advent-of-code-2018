package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	INITIAL     = "####....#...######.###.#...##....#.###.#.###.......###.##..##........##..#.#.#..##.##...####.#..##.#"
	INITIAL_TST = "#..#.#..##......###...###"
)

type Rule struct {
	pat uint8
	set bool
}

func parseRule(s string) *Rule {
	var pat uint8
	var set bool

	for i := 0; i < 5; i++ {
		if s[i] == '#' {
			pat |= (1 << (4 - i))
		}
	}

	if s[9] == '#' {
		set = true
	}
	return &Rule{
		pat: pat,
		set: set,
	}
}

func (r *Rule) String() string {
	return fmt.Sprintf("{pat: %05b, set: %t}", r.pat, r.set)
}

type Game struct {
	state        map[int]bool
	rules        map[uint8]bool
	minix, maxix int
}

func NewGame(s string, rs []*Rule) *Game {
	rules := make(map[uint8]bool)
	for _, rule := range rs {
		rules[rule.pat] = rule.set
	}
	state := make(map[int]bool)
	maxix := 0
	for i := 0; i < len(s); i++ {
		state[i] = (s[i] == '#')
		maxix = max(maxix, i)
	}
	return &Game{
		rules: rules,
		state: state,
		minix: 0,
		maxix: maxix,
	}
}

func getPat(ix int, state map[int]bool) uint8 {
	var pat uint8
	for i := ix - 2; i <= ix+2; i++ {
		if state[i] {
			pat |= 1 << (4 - (i - (ix - 2)))
		}
	}
	return pat
}

func (g *Game) Play() {
	ns := make(map[int]bool)
	for ix := g.minix - 2; ix <= g.maxix+2; ix++ {
		pat := getPat(ix, g.state)
		set := g.rules[pat]
		if set {
			ns[ix] = set
		}
	}
	minix, maxix := ALOT, -ALOT
	for ix, set := range ns {
		if set {
			minix = min(minix, ix)
			maxix = max(maxix, ix)
		}
	}
	g.state = ns
	g.minix = minix
	g.maxix = maxix
}

func (g *Game) String() string {
	var buf bytes.Buffer
	for ix := g.minix; ix <= g.maxix; ix++ {
		if g.state[ix] {
			buf.WriteByte('#')
		} else {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}

func (g *Game) Count() int {
	res := 0
	for ix, set := range g.state {
		if set {
			res += ix
		}
	}
	return res
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	rules := make([]*Rule, 0, len(lines))
	for _, line := range lines {
		rules = append(rules, parseRule(line))
	}
	printf("rules: %+v", rules)

	game := NewGame(INITIAL, rules)
	println(game.String())

	states := make(map[string]bool)
	states[game.String()] = true
	CAP := 50000000000
	ii := 0
	for i := 0; i < CAP; i++ {
		//for i := 0; i < 20; i++ {
		game.Play()
		s := game.String()
		if _, ok := states[s]; ok {
			printf("match after %d steps, minix: %d, count: %d", i, game.minix, game.Count())
			ii++
			if ii > 10 {
				break
			}
		}
		states[s] = true
		println(s)
	}

	// 2022/01/09 19:03:45 match after 96 steps, minix: 25, count: 3432
	// 2022/01/09 19:03:45 match after 97 steps, minix: 26, count: 3464
	// from now on, we only progress forward with a static population
	// at a rate 1 step per iteration, which sums up in 32 extra
	// points on every step
	// the final formula is:
	// 3400+(50000000000-96)*32 == 1600000000328

	//println(game.String())
	printf("pot sum: %d", game.Count())
}
