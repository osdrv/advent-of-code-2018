package main

/*

	N-th score: N * 23 + K
	K - ?

*/

type Game struct {
	marbles []int
	curr    int
}

func NewGame() *Game {
	return &Game{
		curr:    -1,
		marbles: make([]int, 0, 1),
	}
}

func (g *Game) Play(marble int) int {
	if len(g.marbles) < 2 {
		g.marbles = append(g.marbles, marble)
		g.curr = len(g.marbles) - 1
		return 0
	}
	res := 0
	if marble%23 == 0 {
		res += marble
		ix := g.curr - 7
		if ix < 0 {
			ix += len(g.marbles)
		}
		res += g.marbles[ix]
		ms := make([]int, 0, len(g.marbles)-1)
		ms = append(ms, g.marbles[:ix]...)
		ms = append(ms, g.marbles[ix+1:]...)
		g.marbles = ms
		g.curr = ix
	} else {
		ix := g.curr + 2
		if ix > len(g.marbles) {
			ix %= len(g.marbles)
		}
		ms := make([]int, 0, len(g.marbles)+1)
		ms = append(ms, g.marbles[:ix]...)
		ms = append(ms, marble)
		ms = append(ms, g.marbles[ix:]...)

		g.marbles = ms
		g.curr = ix
	}

	//printf("game: %+v, curr: %d", g.marbles, g.curr)

	return res
}

func play(nPlayers int, maxMarble int) int {
	scores := make([]int, nPlayers)
	marble := -1
	nextMarble := func() (int, bool) {
		marble++
		if marble > maxMarble {
			return -1, false
		}
		return marble, true
	}
	ptr := 0
	game := NewGame()
	for {
		m, ok := nextMarble()
		if !ok {
			break
		}
		scores[ptr] += game.Play(m)
		ptr++
		if ptr >= nPlayers {
			ptr -= nPlayers
		}
	}
	maxScore := scores[0]
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	printf("%d players, %d marbles, max score: %d", nPlayers, maxMarble, maxScore)
	return maxScore
}

func main() {
	play(9, 25)
	play(10, 1618)
	play(13, 7999)
	play(17, 1104)
	play(21, 6111)
	play(30, 5807)
	play(452, 7125000)
}
