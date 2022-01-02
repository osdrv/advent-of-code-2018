package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Rect struct {
	id     int
	lt, rb *Point2
}

func NewRect(s string) *Rect {
	chs := strings.FieldsFunc(s[1:], func(r rune) bool {
		return r == ' ' || r == '@' || r == ',' || r == ':' || r == 'x'
	})
	id, j, i, w, h := parseInt(chs[0]), parseInt(chs[1]), parseInt(chs[2]), parseInt(chs[3]), parseInt(chs[4])
	return &Rect{
		id: id,
		lt: NewPoint2(i, j),
		rb: NewPoint2(i+h-1, j+w-1),
	}
}

func (r *Rect) String() string {
	return fmt.Sprintf("R{id: %d, lt: %s, rb: %s}", r.id, r.lt, r.rb)
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	rects := make([]*Rect, 0, len(lines))
	rectids := make(map[int]struct{})
	for _, line := range lines {
		rect := NewRect(line)
		rects = append(rects, rect)
		rectids[rect.id] = struct{}{}
	}

	log.Printf("rects: %+v", rects)

	width, height := 0, 0
	for _, rect := range rects {
		if rect.rb.x > height {
			height = rect.rb.x
		}
		if rect.rb.y > width {
			width = rect.rb.y
		}
	}

	printf("width: %d, height: %d", width, height)

	field := makeIntField(height+1, width+1)
	lastid := makeIntField(height+1, width+1)
	for _, rect := range rects {
		for i := rect.lt.x; i <= rect.rb.x; i++ {
			for j := rect.lt.y; j <= rect.rb.y; j++ {
				field[i][j]++
				if lastid[i][j] > 0 {
					delete(rectids, lastid[i][j])
					delete(rectids, rect.id)
				}
				lastid[i][j] = rect.id
			}
		}
	}

	overlap := 0
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			if field[i][j] > 1 {
				overlap++
			}
		}
	}

	printf("overlap: %d", overlap)

	printf("remaining ids: %+v", rectids)
}
