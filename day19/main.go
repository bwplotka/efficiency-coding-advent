package day19

import (
	"github.com/bwplotka/efficiency-advent-2021/day6"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/pkg/errors"
)

type vec struct {
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

func (b *beacon) manhattanDist(o *beacon) int64 {
	return abs(b.X()-o.X()) + abs(b.Y()-o.Y()) + abs(b.Z()-o.Z())
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
		for input[i] == '-' || (input[i] >= '0' && input[i] <= '9') {
			i++
		}

		x, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		j = i
		for input[i] == '-' || (input[i] >= '0' && input[i] <= '9') {
			i++
		}
		y, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		j = i
		for input[i] == '-' || (input[i] >= '0' && input[i] <= '9') {
			i++
		}
		z, err := day6.ParseInt(input[j:i])
		if err != nil {
			return nil, err
		}
		i++

		b := beacon{
			Vec3:       mgl64.Vec3{float64(x), float64(y), float64(z)},
			neighbours: map[int64]int{},
		}

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

		i          int
		normalized = map[struct{ x, y, z float64 }]struct{}{}
	)

outer:
	for sAb := range scanA {
		for sBb := range scanB {
			bA := &(scanA[sAb])
			bB := &(scanB[sBb])

			if !bA.isPossiblySame(bB, overlapThreshold) {
				continue
			}

			x, y, z := round(bA.Normalize()).Elem()
			k := struct{ x, y, z float64 }{x, y, z}

			if _, ok := normalized[k]; ok {
				continue
			}
			normalized[k] = struct{}{}
			bsScanA[i] = bA.Vec3
			bsScanB[i] = bB.Vec3
			if i < 3 {
				i++
				continue
			}

			break outer

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

	return &transform{
		// Inverse of change CB will get us from scanner B to A.
		rotation: mgl64.Mat3FromCols(xyzc1.Vec3(), xyzc2.Vec3(), xyzc3.Vec3()).Inv(),
		move:     round(mgl64.Vec3{xyzc1.W(), xyzc2.W(), xyzc3.W()}),
	}, true
}

func HowManyBeaconsPart1(input string, overlapThreshold int) (_ int, err error) {
	scannerBeacons, err := parse(input)
	if err != nil {
		return 0, err
	}

	transforms := make([][]*transform, len(scannerBeacons))
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
		}
	}

	// First find the first scanner that overlaps with every other scanner, using brute force.
	// TODO(bwplotka): We could use some graph / tree alg for this (:
	var i = 0
outerTransformationFind:
	for ; i < len(scannerBeacons); i++ {
		for j := 0; j < len(scannerBeacons); j++ {
			if j == i {
				continue
			}
			err = transformBeacons(transforms, i, j, nil)
			if err == nil {
				break outerTransformationFind
			}
		}
	}
	if err != nil {
		return 0, errors.New("no scanner overlaps in direct or indirect manner with every other scanner")
	}

	allBeacons := map[vec]struct{}{}
	for j := 0; j < len(scannerBeacons); j++ {
		if j != i {
			_ = transformBeacons(transforms, i, j, scannerBeacons[j])
		}
		for _, sXb := range scannerBeacons[j] {
			allBeacons[vec{int64(sXb.X()), int64(sXb.Y()), int64(sXb.Z())}] = struct{}{}
		}
	}
	return len(allBeacons), nil
}

func transformBeacons(transforms [][]*transform, to int, from int, bs []beacon) error {
	was := map[int]struct{}{from: {}}
outer:
	for {
		if transforms[to][from] != nil {
			if bs != nil {
				transforms[to][from].transformInPlace(bs)
			}
			return nil
		}

		for i, potentialTo := range transforms {
			if potentialTo[from] == nil {
				continue
			}
			if _, ok := was[i]; ok {
				continue
			}

			was[from] = struct{}{}
			if bs != nil {
				potentialTo[from].transformInPlace(bs)
			}

			from = i
			continue outer
		}
		return errors.New("no chain, we cannot transform")
	}
}
