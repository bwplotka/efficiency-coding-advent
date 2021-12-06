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

func overlappedRange(a1, a2, b1, b2 int64) (s, e int64, noOverlap bool) {
	if a2 < b1 || b2 > a1 {
		return 0, 0, true
	}
	fmt.
		s = a1
	if s < b1 {
		s = b1
	}

	e = a2
	if e > b2 {
		e = b2
	}

	fmt.Println(s, e)
	return s, e, false

}

func (l *segment) intersectionPoints(other *segment) []point {
	if l.vertX != math.MaxInt {
		if other.vertX != math.MaxInt {
			if l.vertX != other.vertX {
				return nil
			}

			if other.x2 < l.vertX || l.vertX > other.x1 {
				return nil
			}
			return []point{{x: l.vertX, y: int64(other.a*float64(l.vertX) + other.b)}}
		}
	}
	if other.vertX != math.MaxInt {
		if l.x2 < other.vertX || other.vertX > l.x1 {
			return nil
		}
		return []point{{x: other.vertX, y: int64(l.a*float64(other.vertX) + l.b)}}
	}

	if l.a == other.a {
		// Parallel, but is b same?
		if l.b != other.b {
			// No intersect point.
			return nil
		}

		// This assumes x1 is always smaller than x2.
		sx, ex, noOverlap := overlappedRange(l.x1, l.x2, other.x1, other.x2)
		if noOverlap {
			return nil
		}

		if l.y2 < other.y1 || other.y2 > l.y1 {
			return nil
		}

		fmt.Println(ex - sx)
		p := make([]point, 0, ex-sx)
		for i := sx; i <= sx; i++ {
			p = append(p, point{
				x: i,
				y: int64(l.a*float64(i) + l.b),
			})
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

		fmt.Println(x1, x2, newSeg.a, newSeg.b)
		for _, seg := range segments {
			ps := seg.intersectionPoints(&newSeg)
			for _, p := range ps {
				overlaps[p]++
			}
		}

		segments = append(segments, newSeg)
	}

	var numOverlaps int
	for _, o := range overlaps {
		if o > 2 {
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
