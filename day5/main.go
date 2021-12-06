package day5

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type point struct {
	x, y int64
}

type segment struct {
	x1, y1 int64
	x2, y2 int64

	// A linear function that this segment is part of.
	// y = ax + b, unless it's vertical line, then vertX is non MaxInt.
	a, b  float64
	vertX int64
}

func newSegment(x1, y1, x2, y2 int64) segment {
	if x1 > x2 {
		x1, x2, y1, y2 = x2, x1, y2, y1
	}

	s := segment{
		x1: x1, y1: y1, x2: x2, y2: y2,
		vertX: math.MaxInt,
	}

	if x2-x1 == 0 {
		// Vertical line.
		s.vertX = x2
		return s
	}

	// a * x1 + b = y1
	// a * x2 + b = y2
	// so:
	// b = y1 - a * x1
	// so:
	// a * x2 + y1 - a * x1 = y2 -> a (x2 - x1) = y2 - y1 -> a = (y2 - y1) / (x2 -x1)
	s.a = float64(y2-y1) / float64(x2-x1)
	s.b = float64(y1) - s.a*float64(x1)
	return s
}

func overlappedRange(a1, a2, b1, b2 int64) (s, e int64, potentialOverlap bool) {
	if a1 > a2 {
		a1, a2 = a2, a1
	}
	if b1 > b2 {
		b1, b2 = b2, b1
	}

	if a2 < b1 || b2 < a1 {
		return 0, 0, false
	}
	s = a1
	if s < b1 {
		s = b1
	}

	e = a2
	if e > b2 {
		e = b2
	}

	return s, e, true
}

func (l *segment) intersectionPoints(other *segment) []point {
	// Rough check of boundaries within "square".
	sx, ex, isOverlap := overlappedRange(l.x1, l.x2, other.x1, other.x2)
	if !isOverlap {
		return nil
	}
	sy, ey, isOverlap := overlappedRange(l.y1, l.y2, other.y1, other.y2)
	if !isOverlap {
		return nil
	}

	if l.vertX != math.MaxInt {
		if other.vertX != math.MaxInt {
			if l.vertX != other.vertX {
				return nil
			}

			// Parallel, and same x.
			p := make([]point, 0, ey-sy)
			for i := sy; i <= ey; i++ {
				p = append(p, point{
					x: l.vertX,
					y: i,
				})
			}
			return p
		}
		y := int64(other.a*float64(l.vertX) + other.b)
		if y < sy || y > ey {
			return nil
		}
		return []point{{x: l.vertX, y: y}}
	}

	if other.vertX != math.MaxInt {
		y := int64(l.a*float64(other.vertX) + l.b)
		if y < sy || y > ey {
			return nil
		}
		return []point{{x: other.vertX, y: y}}
	}

	if l.a == other.a {
		// Parallel, but is b same?
		if l.b != other.b {
			// No intersect point.
			return nil
		}

		p := make([]point, 0, ex-sx)
		for i := sx; i <= ex; i++ {
			y := int64(l.a*float64(i) + l.b)
			if y < sy || y > ey {
				continue
			}

			p = append(p, point{x: i, y: y})
		}
		return p
	}

	// a1 * x + b1 = y
	// a2 * x + b2 = y
	// so:
	// a1 * x + b1 = a2 * x + b2 -> x * (a1 - a2) = b2 - b1 -> x = (b2 - b1) / (a1 - a2)
	var p point
	p.x = int64((other.b - l.b) / (l.a - other.a))
	p.y = int64(l.a*float64(p.x) + l.b)

	if p.x < sx || p.x > ex {
		return nil
	}
	if p.y < sy || p.y > ey {
		return nil
	}
	return []point{p}

}

func VentsOverlapPart1(input string) (_ int, err error) {
	var (
		// TODO(bwplotka): Map can be quite large and slow, so idea could be to maintain sorted array
		// of overlaps by its distance from 0,0.
		overlaps = map[point]int{}
		segments = make([]segment, 0, 500) // 500 is "cheating" - I know max input size is 500.
		line     string
	)
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]

		s := strings.Split(line, " -> ")
		start := strings.Split(s[0], ",")
		end := strings.Split(s[1], ",")

		x1, err := strconv.ParseInt(start[0], 10, 64)
		if err != nil {
			return 0, err
		}
		y1, err := strconv.ParseInt(start[1], 10, 64)
		if err != nil {
			return 0, err
		}

		x2, err := strconv.ParseInt(end[0], 10, 64)
		if err != nil {
			return 0, err
		}
		y2, err := strconv.ParseInt(end[1], 10, 64)
		if err != nil {
			return 0, err
		}

		// Part1.
		if x1 != x2 && y1 != y2 {
			continue
		}

		newSeg := newSegment(x1, y1, x2, y2)

		fmt.Println("got", newSeg.x1, newSeg.y1, "->", newSeg.x2, newSeg.y2, newSeg.a, newSeg.b, newSeg.vertX)
		for _, seg := range segments {
			ps := seg.intersectionPoints(&newSeg)
			if len(ps) > 0 {
				fmt.Println("intersections against", seg.x1, seg.y1, "->", seg.x2, seg.y2, ps)
			}

			for _, p := range ps {
				overlaps[p]++
			}
		}

		segments = append(segments, newSeg)
	}

	var numOverlaps int
	for _, o := range overlaps {
		if o > 0 {
			numOverlaps++
		}
	}

	return numOverlaps, nil
}

func VentsOverlapPart2(input string) (_ int, err error) {
	var (
		// TODO(bwplotka): Map can be quite large and slow, so idea could be to maintain sorted array
		// of overlaps by its distance from 0,0.
		overlaps = map[point]int{}
		segments = make([]segment, 0, 500) // 500 is "cheating" - I know max input size is 500.
		line     string
	)
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]

		s := strings.Split(line, " -> ")
		start := strings.Split(s[0], ",")
		end := strings.Split(s[1], ",")

		x1, err := strconv.ParseInt(start[0], 10, 64)
		if err != nil {
			return 0, err
		}
		y1, err := strconv.ParseInt(start[1], 10, 64)
		if err != nil {
			return 0, err
		}

		x2, err := strconv.ParseInt(end[0], 10, 64)
		if err != nil {
			return 0, err
		}
		y2, err := strconv.ParseInt(end[1], 10, 64)
		if err != nil {
			return 0, err
		}

		newSeg := newSegment(x1, y1, x2, y2)

		// Useful debug log (:
		//fmt.Println("got", newSeg.x1, newSeg.y1, "->", newSeg.x2, newSeg.y2, newSeg.a, newSeg.b, newSeg.vertX)
		for _, seg := range segments {
			ps := seg.intersectionPoints(&newSeg)
			// Useful debug log.
			//if len(ps) > 0 {
			//	fmt.Println("intersections against", seg.x1, seg.y1, "->", seg.x2, seg.y2, ps)
			//}

			for _, p := range ps {
				overlaps[p]++
			}
		}

		segments = append(segments, newSeg)
	}

	var numOverlaps int
	for _, o := range overlaps {
		if o > 0 {
			numOverlaps++
		}
	}

	return numOverlaps, nil
}

type overlap struct {
	x, y int

	occurrences int
}

//type overlaps []overlap
//
//func (o *overlaps) add(x, y int) {
//	// TODO: Use binary search at some point.
//	for i := range *o {
//		if (*o)[i].x == x && (*o)[i].y == y {
//			(*o)[i].occurrences++
//			return
//		}
//	}
//
//	*o = append(*o, overlap{x: x, y: y, occurrences: 1})
//}
