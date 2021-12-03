package day3

import (
	"math"
	"strings"
)

func BinDiagnosePart1(input string) (_ int, err error) {
	var lines int
	var bits []int
	var line string
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]
		lines++

		if len(bits) == 0 {
			bits = make([]int, len(line))
		}

		for j, b := range line {
			if b == '1' {
				bits[j]++
			}
		}
	}

	var epsilon, gamma int
	for i, b := range bits {
		if b > lines/2 {
			gamma += int(math.Pow(2, float64(len(bits)-i-1)))
		} else {
			epsilon += int(math.Pow(2, float64(len(bits)-i-1)))
		}
	}
	return gamma * epsilon, nil
}

func BinDiagnosePart2(input string) (_ int, err error) {
	var ox, co2 []string
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		co2 = append(co2, input[0:i])
		ox = append(ox, input[0:i])
		input = input[i+1:]
	}

	return findRating(ox, true) *
		findRating(co2, false), nil
}

func findRating(numbers []string, common bool) int {
	var pos int
	nLen := float64(len(numbers))
	for nLen > 1 {
		var bits float64
		for _, n := range numbers {
			if n[pos] == '1' {
				bits++
			}
		}

		var filter uint8
		if common {
			filter = uint8('0')
			if bits >= nLen/2 {
				filter = '1'
			}
		} else {
			filter = uint8('0')
			if bits < nLen/2 {
				filter = '1'
			}
		}
		nLen = 0
		for _, n := range numbers {
			if n[pos] == filter {
				numbers[int(nLen)] = n
				nLen++
			}
		}
		numbers = numbers[:int(nLen)]
		pos++
	}

	return binToDecimal(numbers[0])
}

func binToDecimal(binary string) (ret int) {
	for i, b := range binary {
		if b == '1' {
			ret += int(math.Pow(2, float64(len(binary)-i-1)))
		}
	}
	return ret
}

// BinDiagnosePart2_V2 is an optimized version of BinDiagnosePart2.
// Pprof showed appends (grow slice) and looping numbers in findRating as main offenders.
// Let's avoid appends with prealloc, and operate mostly on input.
func BinDiagnosePart2_V2(input string) (_ int, err error) {
	var (
		numLen   = strings.IndexByte(input, '\n')
		numCount = len(input) / (numLen + 1)
	)

	ox := make([]int, 0, numCount)
	co2 := make([]int, 0, numCount)

	// Shared loop for first iteration.
	var bits int
	for i := 0; i < numCount; i++ {
		// Add up?
		// Optimize if we see majority?
		if input[i*(numLen+1)] == '1' {
			bits++
		}
	}

	filter := uint8('0')
	if bits > numCount/2 {
		filter = '1'
	}

	for i := 0; i < numCount; i++ {
		j := i * (numLen + 1)
		if input[j] == filter {
			ox = append(ox, j)
		} else {
			co2 = append(co2, j)
		}
	}

	return findRating_V2(input, numLen, ox, true) * findRating_V2(input, numLen, co2, false), nil
}

func findRating_V2(input string, numLen int, filtered []int, common bool) int {
	pos := 1
	fLen := float64(len(filtered))
	for fLen > 1 {
		var bits float64
		for _, i := range filtered {
			if input[i+pos] == '1' {
				bits++
			}
		}

		var filter uint8
		if common {
			filter = uint8('0')
			if bits >= fLen/2 {
				filter = '1'
			}
		} else {
			filter = uint8('0')
			if bits < fLen/2 {
				filter = '1'
			}
		}
		fLen = 0
		for _, i := range filtered {
			if input[i+pos] == filter {
				filtered[int(fLen)] = i
				fLen++
			}
		}
		filtered = filtered[:int(fLen)]
		pos++
	}

	return binToDecimal(input[filtered[0] : filtered[0]+numLen])
}
