package main

import (
	"fmt"
	"os"
	"regexp"
)

type Object struct {
	pos *Point2
	vel *Point2
}

func (o *Object) String() string {
	return fmt.Sprintf("<pos={%d, %d}, vel={%d, %d}>", o.pos.x, o.pos.y, o.vel.x, o.vel.y)
}

var (
	rr = regexp.MustCompile(`position=<\s*([\-\d]+)\,\s*([\-\d]+)>\svelocity=<\s*([\-\d]+)\,\s*([\-\d]+)>`)
)

func parseObject(s string) *Object {
	match := rr.FindAllStringSubmatch(s, -1)
	assert(match != nil, "no match")
	x, y, vx, vy := parseInt(match[0][1]), parseInt(match[0][2]), parseInt(match[0][3]), parseInt(match[0][4])
	return &Object{
		pos: NewPoint2(x, y),
		vel: NewPoint2(vx, vy),
	}
}

func render(objects []*Object) string {
	minx, miny, maxx, maxy := objects[0].pos.x, objects[0].pos.y, objects[0].pos.x, objects[0].pos.y
	for _, obj := range objects {
		if obj.pos.x < minx {
			minx = obj.pos.x
		}
		if obj.pos.y < miny {
			miny = obj.pos.y
		}
		if obj.pos.x > maxx {
			maxx = obj.pos.x
		}
		if obj.pos.y > maxy {
			maxy = obj.pos.y
		}
	}

	if maxx-minx > 400 || maxy-miny > 400 {
		return ""
	}

	field := makeByteField(maxy-miny+1, maxx-minx+1)
	for _, obj := range objects {
		field[obj.pos.y-miny][obj.pos.x-minx] = 1
	}
	return printByteFieldWithSubs(field, "", map[byte]string{
		0: " ",
		1: "#",
	})
}

func evolve(objects []*Object) {
	for _, obj := range objects {
		obj.pos.x += obj.vel.x
		obj.pos.y += obj.vel.y
	}
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)
	objects := make([]*Object, 0, len(lines))
	for _, line := range lines {
		objects = append(objects, parseObject(line))
	}

	printf("objects: %+v", objects)

	if r := render(objects); len(r) > 0 {
		print(r)
	}

	var input string
	t := 0
	for {
		t++
		evolve(objects)
		println(t)
		if r := render(objects); len(r) > 0 {
			print(r)
			fmt.Scanln(&input)
		}
	}
}
