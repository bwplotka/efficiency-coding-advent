package day7

import (
	"math"
	"sort"
	"strconv"
)

func MinFuelPart1(input string) (_ int, err error) {
	positions := make([]int, 0, 1000)
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}
		pos, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		i++
		positions = append(positions, int(pos))

	}
	sort.Ints(positions)

	median := positions[len(positions)/2]

	var fuelUsed int
	for _, p := range positions {
		if median > p {
			fuelUsed += median - p
			continue
		}

		fuelUsed += p - median
	}
	return fuelUsed, nil
}

func MinFuelPart2(input string) (_ int, err error) {
	positions := make([]int, 0, 1000)

	max := math.MinInt
	min := math.MaxInt
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}
		pos64, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		i++
		pos := int(pos64)
		positions = append(positions, pos)

		if pos > max {
			max = pos
		}
		if pos < min {
			min = pos
		}
	}

	lowestFuel := math.MaxInt
	for i := min; i <= max; i++ {
		bestPos := i

		var fuelUsed int
		for _, p := range positions {
			distance := p - bestPos
			if bestPos > p {
				distance = bestPos - p
			}
			// Arithmetic progression:
			// Sn = ((a1 + an) / 2) * n
			fuelUsed += int(float64(1+distance) / 2.0 * float64(distance))
		}

		if fuelUsed < lowestFuel {
			lowestFuel = fuelUsed
		}
	}

	return lowestFuel, nil
}

func MinFuelPart2_V2(input string) (_ int, err error) {
	positions := make([]int, 0, 1000)

	sum := 0
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}
		pos, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		i++
		positions = append(positions, int(pos))
		sum += int(pos)

	}
	mean := int(math.Ceil(float64(sum) / float64(len(positions))))

	lowestFuel := math.MaxInt
	// Somewhere between those two should be lowest fuel number.
	for i := mean - 1; i <= mean; i++ {
		bestPos := i

		var fuelUsed int
		for _, p := range positions {
			distance := p - bestPos
			if bestPos > p {
				distance = bestPos - p
			}
			// Arithmetic progression:
			// Sn = ((a1 + an) / 2) * n
			fuelUsed += int(float64(1+distance) / 2.0 * float64(distance))
		}

		if fuelUsed < lowestFuel {
			lowestFuel = fuelUsed
		}
	}

	return lowestFuel, nil
}
