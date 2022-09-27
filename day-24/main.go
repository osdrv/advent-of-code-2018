package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Group struct {
	units   int
	hp      int
	atkdmg  int
	atktyp  string
	init    int
	weaks   []string
	immunes []string

	boost int
}

func NewGroup(units int, hp int, atkdmg int, atktyp string, init int, weaks []string, immunes []string) *Group {
	return &Group{
		units:   units,
		hp:      hp,
		atkdmg:  atkdmg,
		atktyp:  atktyp,
		init:    init,
		weaks:   weaks,
		immunes: immunes,

		boost: 0,
	}
}

func (g *Group) StringExt() string {
	var buf bytes.Buffer
	// 17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2
	buf.WriteString(fmt.Sprintf("%d units each with %d hit points ", g.units, g.hp))
	if len(g.weaks) > 0 || len(g.immunes) > 0 {
		buf.WriteByte('(')
		if len(g.immunes) > 0 {
			buf.WriteString(fmt.Sprintf("immune to %s", strings.Join(g.immunes, ", ")))
			if len(g.weaks) > 0 {
				buf.WriteString("; ")
			}
		}
		if len(g.weaks) > 0 {
			buf.WriteString(fmt.Sprintf("weak to %s", strings.Join(g.weaks, ", ")))
		}
		buf.WriteByte(')')
		buf.WriteByte(' ')
	}
	buf.WriteString(fmt.Sprintf("with an attack that does %d %s damage at initiative %d", g.atkdmg, g.atktyp, g.init))
	return buf.String()
}

func (g *Group) String() string {
	return fmt.Sprintf("Group containinig %d units", g.units)
}

func (g *Group) Power() int {
	return g.units * (g.atkdmg + g.boost)
}

func (g *Group) Initiative() int {
	return g.init
}

func (g *Group) IsEmpty() bool {
	return g.units == 0
}

func (g *Group) IsImmutableTo(atktyp string) bool {
	for _, im := range g.immunes {
		if atktyp == im {
			return true
		}
	}
	return false
}

func (g *Group) isWeakTo(atktyp string) bool {
	for _, wk := range g.weaks {
		if atktyp == wk {
			return true
		}
	}
	return false
}

func (g *Group) DealDamage(dmg int) {
	g.units = max(g.units-(dmg/g.hp), 0)
}

func (g *Group) Boost(boost int) {
	g.boost = boost
}

func (g *Group) Copy() *Group {
	return NewGroup(g.units, g.hp, g.atkdmg, g.atktyp, g.init, g.weaks, g.immunes)
}

func parseGroup(s string) *Group {
	ptr := 0
	var units, hp, atkdmg, init int
	var atktyp string
	var weaks, immunes []string

	debugf("Parsing group str: %q", s)

	units, ptr = readNumber(s, ptr)
	ptr = readWhiteSpace(s, ptr)
	ptr = readStaticString(s, ptr, "units each with")
	ptr = readWhiteSpace(s, ptr)
	hp, ptr = readNumber(s, ptr)
	ptr = readWhiteSpace(s, ptr)
	ptr = readStaticString(s, ptr, "hit points")
	ptr = readWhiteSpace(s, ptr)

	if match(s, ptr, '(') {
		ptr = readStaticString(s, ptr, "(")
		immunes, weaks, ptr = readImmunesAndWeaks(s, ptr)
		ptr = readStaticString(s, ptr, ")")
	}
	ptr = readWhiteSpace(s, ptr)
	ptr = readStaticString(s, ptr, "with an attack that does")
	ptr = readWhiteSpace(s, ptr)
	atkdmg, ptr = readNumber(s, ptr)
	ptr = readWhiteSpace(s, ptr)
	atktyp, ptr = readString(s, ptr)
	ptr = readWhiteSpace(s, ptr)
	ptr = readStaticString(s, ptr, "damage at initiative")
	ptr = readWhiteSpace(s, ptr)
	init, ptr = readNumber(s, ptr)

	// 989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3
	return NewGroup(units, hp, atkdmg, atktyp, init, weaks, immunes)
}

func readImmunesAndWeaks(s string, ptr int) ([]string, []string, int) {
	immunes, weaks := make([]string, 0, 1), make([]string, 0, 1)
	if matchStr(s, ptr, "immune to") {
		ptr = readStaticString(s, ptr, "immune to")
		ptr = readWhiteSpace(s, ptr)
		var immune string
		for {
			immune, ptr = readString(s, ptr)
			immunes = append(immunes, immune)
			if !match(s, ptr, ',') {
				break
			}
			ptr = readStaticString(s, ptr, ", ")
		}
	} else if matchStr(s, ptr, "weak to") {
		ptr = readStaticString(s, ptr, "weak to")
		ptr = readWhiteSpace(s, ptr)
		var weak string
		for {
			weak, ptr = readString(s, ptr)
			weaks = append(weaks, weak)
			if !match(s, ptr, ',') {
				break
			}
			ptr = readStaticString(s, ptr, ", ")
		}
	}
	if match(s, ptr, ';') {
		ptr = readStaticString(s, ptr, "; ")
		var ii, ww []string
		ii, ww, ptr = readImmunesAndWeaks(s, ptr)
		immunes = append(immunes, ii...)
		weaks = append(weaks, ww...)
	}
	return immunes, weaks, ptr
}

func readString(s string, ptr int) (string, int) {
	from := ptr
	for ptr < len(s) && isAlpha(s[ptr]) {
		ptr++
	}
	return s[from:ptr], ptr
}

func readStaticString(s string, ptr int, exp string) int {
	off := 0
	for off < len(exp) {
		assert(s[ptr+off] == exp[off], fmt.Sprintf("string mismatch after: %s...", exp[:off]))
		off++
	}
	return ptr + off
}

func readNumber(s string, ptr int) (int, int) {
	from := ptr
	for ptr < len(s) && isNumber(s[ptr]) {
		ptr++
	}
	v := parseInt(s[from:ptr])
	return v, ptr
}

func readWhiteSpace(s string, ptr int) int {
	for ptr < len(s) && isWhiteSpace(s[ptr]) {
		ptr++
	}
	return ptr
}

func match(s string, ptr int, b byte) bool {
	return s[ptr] == b
}

func matchStr(s string, ptr int, exp string) bool {
	return s[ptr:ptr+len(exp)] == exp
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isWhiteSpace(b byte) bool {
	return b == ' '
}

func parseGroups(lines []string) []*Group {
	groups := make([]*Group, 0, len(lines))
	for _, line := range lines {
		groups = append(groups, parseGroup(line))
	}
	return groups
}

func Combat(immune, infect []*Group) bool {
	isImmune := make(map[*Group]bool)
	for _, grp := range immune {
		isImmune[grp] = true
	}
	for len(immune) > 0 && len(infect) > 0 {

		debugf("New Round!")
		// target selection
		allGroups := make([]*Group, 0, len(immune)+len(infect))
		allGroups = append(allGroups, immune...)
		allGroups = append(allGroups, infect...)
		sort.Slice(allGroups, func(i, j int) bool {
			if allGroups[i].Power() != allGroups[j].Power() {
				return allGroups[i].Power() > allGroups[j].Power()
			}
			return allGroups[i].Initiative() > allGroups[j].Initiative()
		})

		targets := make(map[*Group]*Group)
		acquired := make(map[*Group]bool)

		for _, grp := range allGroups {
			target := infect
			if !isImmune[grp] {
				target = immune
			}
			options := make([][2]int, 0, 1)
			for ix, tt := range target {
				if acquired[tt] {
					continue
				}
				dmg := grp.Power()
				if tt.IsImmutableTo(grp.atktyp) {
					dmg = 0
				} else if tt.isWeakTo(grp.atktyp) {
					dmg *= 2
				}
				if dmg > 0 {
					pref := "infection"
					if isImmune[grp] {
						pref = "immune"
					}
					debugf("%s group %q would deal group %q %d damage", pref, grp, target[ix], dmg)
					options = append(options, [2]int{ix, dmg})
				}
			}
			sort.Slice(options, func(i, j int) bool {
				if options[i][1] != options[j][1] {
					return options[i][1] > options[j][1]
				}
				if target[options[i][0]].Power() != target[options[j][0]].Power() {
					return target[options[i][0]].Power() > target[options[j][0]].Power()
				}
				return target[options[i][0]].Initiative() > target[options[j][0]].Initiative()
			})
			if len(options) > 0 {
				targets[grp] = target[options[0][0]]
				acquired[target[options[0][0]]] = true
			}
		}

		if len(targets) == 0 {
			printf("detected a deadlock, returning from the combat")
			return false
		}

		debugf("Attack!")

		sort.Slice(allGroups, func(i, j int) bool {
			return allGroups[i].Initiative() > allGroups[j].Initiative()
		})

		total := 0
		for _, grp := range allGroups {
			if grp.IsEmpty() {
				continue
			}
			target, ok := targets[grp]
			if !ok {
				continue
			}
			dmg := grp.Power()
			if target.isWeakTo(grp.atktyp) {
				dmg *= 2
			}
			assert(dmg > 0, "dmg should be greater than 0")
			before := target.units
			target.DealDamage(dmg)
			after := target.units
			total += before - after
			debugf("group %q attacks group %q and kills %d units", grp, target, before-after)
		}
		if total == 0 {
			printf("no damage was made during the round, returning false")
			return false
		}

		newimmune, newinfect := make([]*Group, 0, 1), make([]*Group, 0, 1)
		for _, grp := range allGroups {
			if grp.IsEmpty() {
				continue
			}
			if isImmune[grp] {
				newimmune = append(newimmune, grp)
			} else {
				newinfect = append(newinfect, grp)
			}
		}
		immune, infect = newimmune, newinfect
	}

	return true
}

func checksum(grps []*Group) int {
	sum := 0
	for _, grp := range grps {
		sum += grp.units
	}
	return sum
}

func copyGroups(grps []*Group) []*Group {
	grpscp := make([]*Group, 0, len(grps))
	for _, grp := range grps {
		grpscp = append(grpscp, grp.Copy())
	}
	return grpscp
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	immune := parseGroups(lines[1 : len(lines)/2])
	infect := parseGroups(lines[len(lines)/2+2:])

	printf("immune: %+v", immune)
	printf("infect: %+v", infect)

	immune1, infect1 := copyGroups(immune), copyGroups(infect)

	Combat(immune1, infect1)

	printf("===== part1 =====")
	printf("winner checksum = %d + %d = %d", checksum(immune1), checksum(infect1), checksum(immune1)+checksum(infect1))

	printf("===== part2 =====")

	minboost := 0
	for boost := 1; boost < 1_000_000; boost++ {
		immuneX, infectX := copyGroups(immune), copyGroups(infect)
		for _, grp := range immuneX {
			grp.Boost(boost)
		}
		printf("probe boost: %d", boost)
		if ok := Combat(immuneX, infectX); !ok {
			continue
		}
		res := checksum(immuneX) > 0
		printf("probe boost %d, res: %t", boost, res)
		if res {
			minboost = boost
			break
		}
	}

	//minboost := sort.Search(1_000_000_000, func(boost int) bool {
	//	immuneX, infectX := copyGroups(immune), copyGroups(infect)
	//	for _, grp := range immuneX {
	//		grp.Boost(boost)
	//	}
	//	printf("probe boost: %d", boost)
	//	Combat(immuneX, infectX)
	//	res := checksum(immuneX) > 0
	//	printf("probe boost %d, res: %t", boost, res)
	//	return res
	//})
	printf("smallest boost: %d", minboost)
	immune2, infect2 := copyGroups(immune), copyGroups(infect)
	for _, grp := range immune2 {
		grp.Boost(minboost)
	}
	Combat(immune2, infect2)
	printf("immune system left with %d units", checksum(immune2))
}
