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
