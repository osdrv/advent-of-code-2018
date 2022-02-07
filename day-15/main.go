package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
)

type UnitType uint8

const (
	_ UnitType = iota
	Elf
	Goblin
)

var (
	UNIT_TYPE_NAME = map[UnitType]string{
		Elf:    "Elf",
		Goblin: "Goblin",
	}
)

var (
	STEPS4RDORD = [][2]int{
		{0, -1},
		{-1, 0},
		{1, 0},
		{0, 1},
	}
)

const (
	FULLHP     = 200
	HIT_POINTS = 3
)

type Unit struct {
	pos      Point2
	hp       int
	typ      UnitType
	strength int
}

func NewUnit(typ UnitType, x, y, hp int) *Unit {
	return &Unit{
		typ:      typ,
		pos:      Point2{x, y},
		hp:       hp,
		strength: HIT_POINTS,
	}
}

func (u *Unit) IsAlive() bool {
	return u.hp > 0
}

func (u *Unit) Hit(damage int) {
	u.hp -= damage
}

func (u *Unit) String() string {
	typ := UNIT_TYPE_NAME[u.typ]
	return fmt.Sprintf("%s{hp: %d, pos: %s}", typ, u.hp, &u.pos)
}

type Path []Point2

type Game struct {
	elves   map[Point2]*Unit
	goblins map[Point2]*Unit
	terrain map[Point2]bool
}

func NewGame(input []string) *Game {
	elves := make(map[Point2]*Unit)
	goblins := make(map[Point2]*Unit)
	terrain := make(map[Point2]bool)
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[0]); x++ {
			switch input[y][x] {
			case '#':
				terrain[Point2{x, y}] = false
			case '.', 'E', 'G':
				terrain[Point2{x, y}] = true
				if input[y][x] == 'E' {
					elves[Point2{x, y}] = NewUnit(Elf, x, y, FULLHP)
				} else if input[y][x] == 'G' {
					goblins[Point2{x, y}] = NewUnit(Goblin, x, y, FULLHP)
				}
			}
		}
	}
	return &Game{
		elves:   elves,
		goblins: goblins,
		terrain: terrain,
	}
}

func (g *Game) getUnits() []*Unit {
	units := make([]*Unit, 0, len(g.elves)+len(g.goblins))

	for _, elf := range g.elves {
		units = append(units, elf)
	}
	for _, goblin := range g.goblins {
		units = append(units, goblin)
	}

	sort.Slice(units, func(i, j int) bool {
		if units[i].pos.y == units[j].pos.y {
			return units[i].pos.x < units[j].pos.x
		}
		return units[i].pos.y < units[j].pos.y
	})

	return units
}

func (g *Game) getEnemies(unit *Unit) map[Point2]*Unit {
	if unit.typ == Elf {
		return g.goblins
	}
	return g.elves
}

func (g *Game) unitPaths(unit *Unit) []Path {
	enemies := g.getEnemies(unit)
	dests := make([]Point2, 0, 1)
	for pos := range enemies {
		for _, step := range STEPS4RDORD {
			p1 := Point2{pos.x + step[0], pos.y + step[1]}
			if g.terrain[p1] && g.goblins[p1] == nil && g.elves[p1] == nil {
				dests = append(dests, p1)
			}
		}
	}

	sort.Slice(dests, func(i, j int) bool {
		if dests[i].y == dests[j].y {
			return dests[i].x < dests[j].x
		}
		return dests[i].y < dests[j].y
	})

	paths := make([]Path, 0, 1)

	minlen := ALOT
	for _, dest := range dests {
		if path, ok := g.findPath(unit.pos, dest); ok {
			paths = append(paths, path)
			minlen = min(minlen, len(path))
		}
	}

	if len(paths) == 0 {
		return paths
	}

	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) == len(paths[j]) {
			if paths[i][1].y == paths[j][1].y {
				return paths[i][1].x < paths[j][1].x
			}
			return paths[i][1].y < paths[j][1].y
		}
		return len(paths[i]) < len(paths[j])
	})

	minpaths := make([]Path, 0, 1)
	for _, path := range paths {
		if len(path) > minlen {
			break
		}
		minpaths = append(minpaths, path)
	}

	return minpaths
}

func (g *Game) findPath(p1, p2 Point2) (Path, bool) {
	var visit func(p Point2, dist int) (Path, bool)
	visited := make(map[Point2]int)
	visit = func(p Point2, dist int) (Path, bool) {
		visited[p] = dist
		if p2 == p {
			return []Point2{p}, true
		}
		var minPath Path
		found := false
		for _, step := range STEPS4RDORD {
			pnext := Point2{p.x + step[0], p.y + step[1]}
			if !g.terrain[pnext] || g.elves[pnext] != nil || g.goblins[pnext] != nil {
				continue
			}
			if prev, ok := visited[pnext]; !ok || prev > (dist+1) {
				if path, ok1 := visit(pnext, dist+1); ok1 {
					found = found || ok1
					if minPath == nil || len(minPath) > len(path) {
						minPath = path
					}
				}
			}
		}
		if !found {
			return nil, false
		}
		return append([]Point2{p}, minPath...), minPath != nil
	}
	return visit(p1, 0)
}

func (g *Game) Play() bool {
	units := g.getUnits()
	printf("unit order: %+v", units)
UNIT:
	for _, unit := range units {
		printf("===================")
		if !unit.IsAlive() {
			// unit is dead, skip
			continue UNIT
		}
		printf("unit %s is playing", unit)

		enemies := g.getEnemies(unit)
		printf("unit enemies: %+v", enemies)
		if len(enemies) == 0 {
			printf("no enemies found, the game is over")
			// game is over
			return false
		}

		unitenemies := make([]*Unit, 0, 1)
		for _, step := range STEPS4RDORD {
			epoint := Point2{unit.pos.x + step[0], unit.pos.y + step[1]}
			if g.terrain[epoint] && enemies[epoint] != nil {
				unitenemies = append(unitenemies, enemies[epoint])
			}
		}

		var dest Point2
		var paths []Path

		if len(unitenemies) > 0 {
			goto ATACK
		}

		printf("unit will be moving")

		paths = g.unitPaths(unit)
		printf("unit paths: %+v", paths)
		if len(paths) == 0 {
			printf("no paths available, moving on")
			continue UNIT
		}

		sort.Slice(paths, func(i, j int) bool {
			pi, pj := paths[i][len(paths[i])-1], paths[j][len(paths[j])-1]
			if pi.y == pj.y {
				return pi.x < pj.x
			}
			return pi.y < pj.y
		})

		dest = paths[0][1]

		printf("unit %s moving to %s", unit, &dest)

		if unit.typ == Elf {
			delete(g.elves, unit.pos)
			unit.pos = dest
			g.elves[dest] = unit
		} else {
			delete(g.goblins, unit.pos)
			unit.pos = dest
			g.goblins[dest] = unit
		}

		unitenemies = make([]*Unit, 0, 1)
		for _, step := range STEPS4RDORD {
			epoint := Point2{unit.pos.x + step[0], unit.pos.y + step[1]}
			if g.terrain[epoint] && enemies[epoint] != nil {
				unitenemies = append(unitenemies, enemies[epoint])
			}
		}

	ATACK:
		if len(unitenemies) > 0 {
			printf("unit is ready to attack")

			sort.Slice(unitenemies, func(i, j int) bool {
				if unitenemies[i].hp == unitenemies[j].hp {
					if unitenemies[i].pos.y == unitenemies[j].pos.y {
						return unitenemies[i].pos.x < unitenemies[j].pos.x
					}
					return unitenemies[i].pos.y < unitenemies[j].pos.y
				}
				return unitenemies[i].hp < unitenemies[j].hp
			})
			printf("enemies in range: %+v", unitenemies)
			enemy := unitenemies[0]
			printf("unit %s hits %s", unit, enemy)
			enemy.Hit(unit.strength)
			if !enemy.IsAlive() {
				printf("unit %s dies", enemy)
				if enemy.typ == Elf {
					delete(g.elves, enemy.pos)
				} else {
					delete(g.goblins, enemy.pos)
				}
			}
			continue UNIT
		}
	}

	return true
}

func (g *Game) String() string {
	var buf bytes.Buffer
	maxx, maxy := 0, 0
	for point := range g.terrain {
		maxx = max(maxx, point.x)
		maxy = max(maxy, point.y)
	}

	for y := 0; y <= maxy; y++ {
		for x := 0; x <= maxx; x++ {
			p := Point2{x, y}
			cp := '#'
			if g.terrain[p] {
				cp = '.'
			}
			if g.elves[p] != nil {
				cp = 'E'
			}
			if g.goblins[p] != nil {
				cp = 'G'
			}
			buf.WriteRune(cp)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func main() {
	TALKATIVE = false
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	elfstrength := 3
GAME:
	for {
		game := NewGame(lines)
		printf(game.String())

		elfcnt := len(game.elves)
		for _, elf := range game.elves {
			elf.strength = elfstrength
		}

		rounds := 0
		//var input string
		for {
			//fmt.Scanln(&input)
			res := game.Play()
			if len(game.elves) != elfcnt {
				log.Printf("1 or more elves died with strength %d, restarting the game", elfstrength)
				elfstrength++
				continue GAME
			}
			printf(game.String())
			if !res {
				break
			}
			rounds++
		}

		log.Printf("all elves are alive, the strength is: %d", elfstrength)

		hpsum := 0
		for _, unit := range game.elves {
			printf("unit %s +%dhp", unit, unit.hp)
			hpsum += unit.hp
		}
		for _, unit := range game.goblins {
			printf("unit %s +%dhp", unit, unit.hp)
			hpsum += unit.hp
		}

		log.Printf("rounds played: %d", rounds)

		log.Printf("total score: %d * %d = %d", rounds, hpsum, rounds*hpsum)
		break GAME
	}
}
