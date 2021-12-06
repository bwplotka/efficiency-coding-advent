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

func (l *segment) intersectionPoints(other *segment) (ret []point) {
	// Rough check of boundaries within "square".
	sx, ex, isOverlap := overlappedRange(l.x1, l.x2, other.x1, other.x2)
	if !isOverlap {
		return nil
	}
	sy, ey, isOverlap := overlappedRange(l.y1, l.y2, other.y1, other.y2)
	if !isOverlap {
		return nil
	}

	defer func() {
		for _, p := range ret {
			if p.x < sx || p.x > ex {
				panic("found point outside of x")
			}
			if p.y < sy || p.y > ey {
				panic("found point outside of y")
			}
		}
	}()

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
				fmt.Println("outside!!")
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
	xFloat := (other.b - l.b) / (l.a - other.a)
	if xFloat != math.Trunc(xFloat) {
		// Our space is discrete.
		return nil
	}
	p.x = int64(xFloat)
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
		// fmt.Println("got", newSeg.x1, newSeg.y1, "->", newSeg.x2, newSeg.y2, newSeg.a, newSeg.b, newSeg.vertX)
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

// VentsOverlapPart2_V2 is optimized version of VentsOverlapPart2.
// Main offendant are intersection functions as well as mapassign (as we predicted).
// For space it's mainly []points and splits while parsing.
func VentsOverlapPart2_V2(input string) (_ int, err error) {
	var (
		// TODO(bwplotka): Map can be quite large and slow, so idea could be to maintain sorted array
		// of overlaps by its distance from 0,0.
		overlaps = make(map[point]struct{}, 500)
		segments = make([]segment_V2, 0, 500) // 500 is "cheating" - I know max input size is 500.
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

		newSeg := newSegment_V2(x1, y1, x2, y2)

		var p point
		markFn := func(x, y int64) {
			p.x = x
			p.y = y
			if _, ok := overlaps[p]; !ok {
				overlaps[p] = struct{}{}
			}
		}
		for _, seg := range segments {
			seg.markIntersectionPoints(&newSeg, markFn)
		}

		segments = append(segments, newSeg)
	}

	return len(overlaps), nil
}

type segment_V2 struct {
	x1, y1 int64
	x2, y2 int64

	// A linear function that this segment is part of.
	// y = ax + b, unless it's vertical line, then vertX is non MaxInt.
	a, b  float64
	vertX int64
}

func newSegment_V2(x1, y1, x2, y2 int64) segment_V2 {
	if x1 > x2 {
		x1, x2, y1, y2 = x2, x1, y2, y1
	}

	s := segment_V2{
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

func overlappedRange_V2(a1, a2, b1, b2 int64) (s, e int64, potentialOverlap bool) {
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

func (l *segment_V2) markIntersectionPoints(other *segment_V2, markFn func(x, y int64)) {
	// Rough check of boundaries within "square".
	sx, ex, isOverlap := overlappedRange_V2(l.x1, l.x2, other.x1, other.x2)
	if !isOverlap {
		return
	}
	sy, ey, isOverlap := overlappedRange_V2(l.y1, l.y2, other.y1, other.y2)
	if !isOverlap {
		return
	}

	if l.vertX != math.MaxInt {
		if other.vertX != math.MaxInt {
			if l.vertX != other.vertX {
				return
			}

			// Parallel, and same x.
			for i := sy; i <= ey; i++ {
				markFn(l.vertX, i)
			}
			return
		}
		y := int64(other.a*float64(l.vertX) + other.b)
		if y < sy || y > ey {
			return
		}
		markFn(l.vertX, y)
		return
	}

	if other.vertX != math.MaxInt {
		y := int64(l.a*float64(other.vertX) + l.b)
		if y < sy || y > ey {
			return
		}
		markFn(other.vertX, y)
		return
	}

	if l.a == other.a {
		// Parallel, but is b same?
		if l.b != other.b {
			// No intersect point.
			return
		}

		for i := sx; i <= ex; i++ {
			y := int64(l.a*float64(i) + l.b)
			if y < sy || y > ey {
				continue
			}

			markFn(i, y)
		}
		return
	}

	// a1 * x + b1 = y
	// a2 * x + b2 = y
	// so:
	// a1 * x + b1 = a2 * x + b2 -> x * (a1 - a2) = b2 - b1 -> x = (b2 - b1) / (a1 - a2)
	xFloat := (other.b - l.b) / (l.a - other.a)
	if xFloat != math.Trunc(xFloat) {
		// Our space is discrete.
		return
	}
	x := int64(xFloat)
	y := int64(l.a*float64(x) + l.b)

	if x < sx || x > ex || y < sy || y > ey {
		return
	}
	markFn(x, y)
}

// VentsOverlapPart2_V3 is optimized version of VentsOverlapPart2_V2.
// Main offendant is still intersection functions as well as mapassign (as we predicted).
// For space it's mainly splits while parsing.
func VentsOverlapPart2_V3(input string) (_ int, err error) {
	var (
		// Map can be quite large and slow, so idea could be to maintain an array of 1000*1000 elements.
		// That's 1MB, which we need to live with (: Trade-off to win latency.
		overlaps    = make([]bool, 1000*1000)
		newOverlaps int

		segments = make([]segment_V2, 0, 500) // 500 is "cheating" - I know max input size is 500.
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

		newSeg := newSegment_V2(x1, y1, x2, y2)

		markFn := func(x, y int64) {
			i := x + 1000*y
			if overlaps[i] {
				return
			}
			overlaps[i] = true

			newOverlaps++
		}
		for _, seg := range segments {
			seg.markIntersectionPoints(&newSeg, markFn)
		}

		segments = append(segments, newSeg)
	}

	return newOverlaps, nil
}

// VentsOverlapPart2_V4 is optimized version of VentsOverlapPart2_V3.
// Main offendant is still intersection functions and 6% split.
func VentsOverlapPart2_V4(input string) (_ int, err error) {
	var (
		// Map can be quite large and slow, so idea could be to maintain an array of 1000*1000 elements.
		// That's 1MB, which we need to live with (: Trade-off to win latency.
		overlaps    = make([]bool, 1000*1000)
		newOverlaps int

		segments = make([]segment_V2, 0, 500) // 500 is "cheating" - I know max input size is 500.
	)

	for i := 0; i < len(input); {
		j := i
		for input[i] != ',' {
			i++
		}
		x1, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}

		i++
		j = i
		for input[i] != ' ' {
			i++
		}
		y1, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}

		i += 4
		j = i
		for input[i] != ',' {
			i++
		}
		x2, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}

		i++
		j = i
		for input[i] != '\n' {
			i++
		}
		y2, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		i++

		newSeg := newSegment_V2(x1, y1, x2, y2)
		markFn := func(x, y int64) {
			i := x + 1000*y
			if overlaps[i] {
				return
			}
			overlaps[i] = true

			newOverlaps++
		}
		for _, seg := range segments {
			seg.markIntersectionPoints(&newSeg, markFn)
		}

		segments = append(segments, newSeg)
	}

	return newOverlaps, nil
}
