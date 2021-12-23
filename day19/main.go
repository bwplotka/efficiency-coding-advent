package day19

import (
	"math"

	"github.com/bwplotka/efficiency-advent-2021/day6"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/pkg/errors"
)

type vec3 struct {
	x, y, z int64
}

type beacon struct {
	mgl64.Vec3

	// Manhattan distances to others within the scanner beacon is in.
	neighbours map[int64]int
}

func abs(x float64) int64 {
	if x < 0 {
		return int64(-1 * x)
	}
	return int64(x)
}

func fillNeighbours(beacons []beacon) {
	for j := 0; j < len(beacons)-1; j++ {
		for k := j + 1; k < len(beacons); k++ {
			dist := beacons[j].manhattanDist(&beacons[k])
			beacons[j].neighbours[dist]++
			beacons[k].neighbours[dist]++
		}
	}
}

func manhattanDist(a, b mgl64.Vec3) int64 {
	return abs(a.X()-b.X()) + abs(a.Y()-b.Y()) + abs(a.Z()-b.Z())
}

func (b *beacon) manhattanDist(o *beacon) int64 {
	return manhattanDist(b.Vec3, o.Vec3)
}

func (b *beacon) isPossiblySame(o *beacon, overlapThreshold int) bool {
	var overlaps int
	for n := range b.neighbours {
		overlaps += o.neighbours[n]
		if overlaps >= (overlapThreshold - 1) { // -1 because of the `b` is accounted already.
			return true
		}
	}
	return false
}

func parse(input string) ([][]beacon, error) {
	var (
		currScanner    = -1
		scannerBeacons [][]beacon
	)
	for i := 0; i < len(input); {
		if input[i] == '\n' {
			i++
			continue
		}

		if input[i] == '-' && input[i+1] == '-' {
			currScanner++
			scannerBeacons = append(scannerBeacons, nil)
			for input[i] != '\n' {
				i++
			}
			i++
			continue
		}

		j := i
		for ; input[i] == '-' || (input[i] >= '0' && input[i] <= '9'); i++ {
		}

		x, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		j = i
		for ; input[i] == '-' || (input[i] >= '0' && input[i] <= '9'); i++ {
		}

		y, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		j = i
		for ; input[i] == '-' || (input[i] >= '0' && input[i] <= '9'); i++ {
		}

		z, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		b := beacon{Vec3: mgl64.Vec3{float64(x), float64(y), float64(z)}, neighbours: map[int64]int{}}
		scannerBeacons[currScanner] = append(scannerBeacons[currScanner], b)
	}
	return scannerBeacons, nil
}

type transform struct {
	rotation mgl64.Mat3
	move     mgl64.Vec3
}

func (t *transform) transformInPlace(bs []beacon) {
	for i := range bs {
		bs[i].Vec3 = round(t.rotation.Mul3x1(bs[i].Vec3)).Add(t.move)
	}
}
func (t *transform) transformVec3(pos mgl64.Vec3) mgl64.Vec3 {
	return round(t.rotation.Mul3x1(pos)).Add(t.move)
}

func round(v mgl64.Vec3) mgl64.Vec3 {
	v[0] = mgl64.Round(v.X(), 0)
	v[1] = mgl64.Round(v.Y(), 0)
	v[2] = mgl64.Round(v.Z(), 0)
	return v
}

func findTransform(scanA, scanB []beacon, overlapThreshold int) (t *transform, overlap bool) {
	// Check if there are at least overlapThreshold pair of "same" beacons.
	// Gather four of those pairs, since we need to solve 4 vars equation to find rot and translation.
	// NOTE: Those four cannot have the same normalized vector (same function). Otherwise our equation won't be solvable.
	//which in normalized are not the same formFind four the same pair beacons from each overlapping scanners (if any).
	var (
		bsScanA [4]mgl64.Vec3
		bsScanB [4]mgl64.Vec3

		i           = -1
		normalized1 = map[struct{ x, y, z float64 }]struct{}{}
		normalized2 = map[struct{ x, y, z float64 }]struct{}{}
	)

outer:
	for sAb := range scanA {
		for sBb := range scanB {

			bA := &(scanA[sAb])
			bB := &(scanB[sBb])

			if !bA.isPossiblySame(bB, overlapThreshold) {
				continue
			}

			// We can't use pair of points for matching that when normalized is exactly the same as another pair.
			x, y, z := round(bA.Normalize()).Elem()
			k := struct{ x, y, z float64 }{x, y, z}
			_, ok1 := normalized1[k]

			x, y, z = round(bB.Normalize()).Elem()
			k2 := struct{ x, y, z float64 }{x, y, z}
			_, ok2 := normalized2[k2]
			if ok1 && ok2 {
				continue
			}

			normalized1[k2] = struct{}{}
			normalized2[k2] = struct{}{}

			i++
			bsScanA[i] = bA.Vec3
			bsScanB[i] = bB.Vec3
			if i >= 3 {
				break outer
			}
		}
	}
	if i < 3 {
		return nil, false
	}

	// We need to find what rotation and offset has to be applied to transform scan B points to scan A coordinates..
	// Useful docs:
	//  * https://towardsdatascience.com/essential-math-for-data-science-basis-and-change-of-basis-f7af2348d463
	//  * https://medium.com/@vasanth260m12/solving-linear-regression-using-linear-algebra-in-golang-b0d66b7056ff
	//
	// Overall we need to find change basis + translation, so. We can do it since we now translation looks like this:
	//
	// [v]A = CB * [b]B + [off], where CB is change basis matrix 3x3 and off is offset vector.
	//
	// Now we "just" need to solve three times four equations with four variables to find CB and off.
	//
	// bsScanA[0].X = x1*bsScanB[0].X + y1*bsScanB[0].Y + z1*bsScanB[0].Z + c1
	// bsScanA[1].X = x1*bsScanB[1].X + y1*bsScanB[1].Y + z1*bsScanB[1].Z + c1
	// bsScanA[2].X = x1*bsScanB[2].X + y1*bsScanB[2].Y + z1*bsScanB[2].Z + c1
	// bsScanA[3].X = x1*bsScanB[3].X + y1*bsScanB[3].Y + z1*bsScanB[3].Z + c1
	// .. and the same for x2, x3, y2, c2 and y3, z2, z3, c3 (:

	// NOTE: Within our task we can assume all x, y and z can only be integers -1, 0 or 1, still this algo should work for other cases too.

	invCoefficients := mgl64.Mat4{
		bsScanB[0].X(), bsScanB[1].X(), bsScanB[2].X(), bsScanB[3].X(), // Col 0.
		bsScanB[0].Y(), bsScanB[1].Y(), bsScanB[2].Y(), bsScanB[3].Y(),
		bsScanB[0].Z(), bsScanB[1].Z(), bsScanB[2].Z(), bsScanB[3].Z(),
		1, 1, 1, 1,
	}.Inv()

	xyzc1 := invCoefficients.Mul4x1(mgl64.Vec4{
		bsScanA[0].X(),
		bsScanA[1].X(),
		bsScanA[2].X(),
		bsScanA[3].X(),
	})

	xyzc2 := invCoefficients.Mul4x1(mgl64.Vec4{
		bsScanA[0].Y(),
		bsScanA[1].Y(),
		bsScanA[2].Y(),
		bsScanA[3].Y(),
	})

	xyzc3 := invCoefficients.Mul4x1(mgl64.Vec4{
		bsScanA[0].Z(),
		bsScanA[1].Z(),
		bsScanA[2].Z(),
		bsScanA[3].Z(),
	})

	tr := &transform{
		// Inverse of change CB will get us from scanner B to A.
		rotation: mgl64.Mat3FromCols(xyzc1.Vec3(), xyzc2.Vec3(), xyzc3.Vec3()).Inv(),
		move:     round(mgl64.Vec3{xyzc1.W(), xyzc2.W(), xyzc3.W()}),
	}

	return tr, true
}

func HowManyBeaconsPart1(input string, overlapThreshold int) (_ int, err error) {
	scannerBeacons, err := parse(input)
	if err != nil {
		return 0, err
	}

	transforms := make([][]*transform, len(scannerBeacons))
	to0transforms := make([]bool, len(scannerBeacons))
	for i := 0; i < len(scannerBeacons); i++ {
		fillNeighbours(scannerBeacons[i])
		transforms[i] = make([]*transform, len(scannerBeacons))
	}

	// Find all transforms and hope that all in chained way overlap to 0 :crossed_fingers:
	// TODO(bwplotka): We could optimize by checking only half and then inverting vice versa transformation.
	for j := 0; j < len(scannerBeacons); j++ {
		for k := 0; k < len(scannerBeacons); k++ {
			if j == k {
				continue
			}

			tr, ok := findTransform(scannerBeacons[j], scannerBeacons[k], overlapThreshold)
			if !ok {
				continue
			}
			transforms[j][k] = tr
			to0transforms[k] = true // Potentially true, we will double-check later.
		}
	}

	// Assume 0 overlaps with others in direct or indirect way.
	allBeacons := map[vec3]struct{}{}
	for _, sXb := range scannerBeacons[0] {
		allBeacons[vec3{int64(sXb.X()), int64(sXb.Y()), int64(sXb.Z())}] = struct{}{}
	}

	cacheTransformStepsFromTo0 := make([][]int, len(transforms))
	for j := 1; j < len(scannerBeacons); j++ {
		if err := findTransformStepsTo0(transforms, j, to0transforms, cacheTransformStepsFromTo0); err != nil {
			return 0, err
		}

		prev := j
		for _, step := range cacheTransformStepsFromTo0[j] {
			transforms[step][prev].transformInPlace(scannerBeacons[j])
			prev = step
		}
		for _, sXb := range scannerBeacons[j] {
			allBeacons[vec3{int64(sXb.X()), int64(sXb.Y()), int64(sXb.Z())}] = struct{}{}
		}
	}
	return len(allBeacons), nil
}

type stack struct {
	s []int
}

func (s *stack) pop() int {
	i := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return i
}

func (s *stack) push(v int) {
	s.s = append(s.s, v)
}

func (s *stack) len() int {
	return len(s.s)
}

// TODO(bwplotka): Probably too complex, but it's 3am...
func findTransformStepsTo0(transforms [][]*transform, from int, to0transforms []bool, cacheSteps [][]int) error {
	if len(cacheSteps[from]) > 0 {
		return nil
	}

	// Use stack for DFS.
	var steps stack
	steps.push(from)

	defer func() {
		// On way back, always unwind cache steps and fill intermediate steps.
		// This allows us to hit cache before we hit preset to0transforms = false for that item.
		for i, step := range cacheSteps[from] {
			if len(cacheSteps[step]) > 0 && len(cacheSteps[step]) < len(cacheSteps[from][i+1:]) {
				continue
			}
			cacheSteps[step] = cacheSteps[from][i+1:]
		}

	}()

stackLoop:
	for steps.len() > 0 {
		currFrom := steps.s[steps.len()-1]
		for i, potentialTo := range transforms {
			if transforms[0][currFrom] != nil {
				// Got it, cache and return.
				cacheSteps[from] = append(cacheSteps[from], append(steps.s[1:], 0)...)
				return nil
			}

			if potentialTo[currFrom] == nil {
				continue
			}

			if len(cacheSteps[i]) > 0 {
				cacheSteps[from] = append(cacheSteps[from], append(append(steps.s[1:], i), cacheSteps[i]...)...)
				return nil
			}

			if !to0transforms[i] {
				// We learned previously, that this path is not chainable to 0.
				continue
			}

			to0transforms[currFrom] = false
			steps.push(i)
			continue stackLoop
		}

		// No chain from here, go back to parent.
		to0transforms[steps.pop()] = false
	}
	return errors.New("no chain, we cannot transform")
}

func ManhattanDistPart2(input string, overlapThreshold int) (_ int, err error) {
	scannerBeacons, err := parse(input)
	if err != nil {
		return 0, err
	}

	transforms := make([][]*transform, len(scannerBeacons))
	to0transforms := make([]bool, len(scannerBeacons))
	for i := 0; i < len(scannerBeacons); i++ {
		fillNeighbours(scannerBeacons[i])
		transforms[i] = make([]*transform, len(scannerBeacons))
	}

	// Find all transforms and hope that all in chained way overlap to 0 :crossed_fingers:
	// TODO(bwplotka): We could optimize by checking only half and then inverting vice versa transformation.
	for j := 0; j < len(scannerBeacons); j++ {
		for k := 0; k < len(scannerBeacons); k++ {
			if j == k {
				continue
			}
			tr, ok := findTransform(scannerBeacons[j], scannerBeacons[k], overlapThreshold)
			if !ok {
				continue
			}
			transforms[j][k] = tr
			to0transforms[k] = true // Potentially true, we will double-check later.
		}
	}

	cacheTransformStepsFromTo0 := make([][]int, len(transforms))

	scannerPositions := make([]mgl64.Vec3, len(scannerBeacons))
	for j := 1; j < len(scannerPositions); j++ {
		if err := findTransformStepsTo0(transforms, j, to0transforms, cacheTransformStepsFromTo0); err != nil {
			return 0, err
		}

		prev := j
		for _, step := range cacheTransformStepsFromTo0[j] {
			scannerPositions[j] = transforms[step][prev].transformVec3(scannerPositions[j])
			prev = step
		}
	}

	max := math.MinInt64
	for i := 0; i < len(scannerPositions); i++ {
		for j := 0; j < len(scannerPositions); j++ {
			if j == i {
				continue
			}

			d := int(manhattanDist(scannerPositions[i], scannerPositions[j]))
			if d > max {
				max = d
			}
		}
	}
	return max, nil
}
