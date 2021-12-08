package day6

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	daysToFish = 7
	delayDays  = 2 // Days for new fish to be delayed for.
)

func SimLanternfish(input string, days int) (_ int, err error) {
	var lanternfishes []int64
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}
		age, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		lanternfishes = append(lanternfishes, age)
		i++
	}

	for d := 0; d < days; d++ {
		var fishesToAdd int
		for i := 0; i < len(lanternfishes); i++ {
			if lanternfishes[i] == 0 {
				lanternfishes[i] = daysToFish - 1
				fishesToAdd++
				continue
			}
			lanternfishes[i]--
		}
		for i := 0; i < fishesToAdd; i++ {
			lanternfishes = append(lanternfishes, daysToFish-1+delayDays)
		}
	}
	return len(lanternfishes), nil
}

func SimLanternfish_V2(input string, days int) (_ int64, err error) {
	fishesDuringDay := make([]int64, 7)
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}
		age, err := strconv.ParseInt(input[j:i], 10, 64)
		if err != nil {
			return 0, err
		}
		fishesDuringDay[age]++
		i++
	}

	index := 0
	var incubatedDay0, incubatedDay1 int64
	for d := 0; d < days; d++ {
		incubated := incubatedDay0
		incubatedDay0 = incubatedDay1
		// Let's create new fishes!
		incubatedDay1 = fishesDuringDay[index]
		// Incubated fishes can have 6 age.
		fishesDuringDay[index] += incubated

		if index == 6 {
			index = -1
		}
		index++
	}

	return incubatedDay0 + incubatedDay1 + fishesDuringDay[0] + fishesDuringDay[1] +
		fishesDuringDay[2] + fishesDuringDay[3] + fishesDuringDay[4] + fishesDuringDay[5] + fishesDuringDay[6], nil
}

// ParseInt is 3-4x times faster than strconv.ParseInt or Atoi.
func ParseInt(input string) (n int64, _ error) {
	factor := int64(1)
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] < '0' || input[i] > '9' {
			return 0, errors.Errorf("not a valid integerer: %v", input)
		}

		n += factor * int64(input[i]-'0')
		factor *= 10
	}

	return n, nil
}

func SimLanternfish_V3(input string, days int) (_ int64, err error) {
	fishesDuringDay := make([]int64, 7)
	for i := 0; i < len(input); {
		j := i

		for input[i] >= '0' && input[i] <= '9' {
			i++
		}

		age, err := ParseInt(input[j:i])
		if err != nil {
			return 0, err
		}

		fishesDuringDay[age]++
		i++
	}

	index := 0
	var incubatedDay0, incubatedDay1 int64
	for d := 0; d < days; d++ {
		incubated := incubatedDay0
		incubatedDay0 = incubatedDay1
		// Let's create new fishes!
		incubatedDay1 = fishesDuringDay[index]
		// Incubated fishes can have 6 age.
		fishesDuringDay[index] += incubated

		if index == 6 {
			index = -1
		}
		index++
	}

	return incubatedDay0 + incubatedDay1 + fishesDuringDay[0] + fishesDuringDay[1] +
		fishesDuringDay[2] + fishesDuringDay[3] + fishesDuringDay[4] + fishesDuringDay[5] + fishesDuringDay[6], nil
}

func SimLanternfish_V4(input string, days int) (_ int64, err error) {
	fishesDuringDay := make([]int64, 7)
	for i := 0; i < len(input); {
		switch input[i] {
		case '0':
			fishesDuringDay[0]++
		case '1':
			fishesDuringDay[1]++
		case '2':
			fishesDuringDay[2]++
		case '3':
			fishesDuringDay[3]++
		case '4':
			fishesDuringDay[4]++
		case '5':
			fishesDuringDay[5]++
		case '6':
			fishesDuringDay[6]++
		}
		i += 2
	}

	index := 0
	var incubatedDay0, incubatedDay1 int64
	for d := 0; d < days; d++ {
		incubated := incubatedDay0
		incubatedDay0 = incubatedDay1
		// Let's create new fishes!
		incubatedDay1 = fishesDuringDay[index]
		// Incubated fishes can have 6 age.
		fishesDuringDay[index] += incubated

		if index == 6 {
			index = -1
		}
		index++
	}

	return incubatedDay0 + incubatedDay1 + fishesDuringDay[0] + fishesDuringDay[1] +
		fishesDuringDay[2] + fishesDuringDay[3] + fishesDuringDay[4] + fishesDuringDay[5] + fishesDuringDay[6], nil
}

func SimLanternfish_V5(input string, days int) (_ int64, err error) {
	fishesDuringDay := make([]int64, 7)
	for i := 0; i < len(input); {
		fishesDuringDay[int(input[i]-'0')]++
		i += 2
	}

	index := 0
	var incubatedDay0, incubatedDay1 int64
	for d := 0; d < days; d++ {
		incubated := incubatedDay0
		incubatedDay0 = incubatedDay1
		// Let's create new fishes!
		incubatedDay1 = fishesDuringDay[index]
		// Incubated fishes can have 6 age.
		fishesDuringDay[index] += incubated

		if index == 6 {
			index = -1
		}
		index++
	}

	return incubatedDay0 + incubatedDay1 + fishesDuringDay[0] + fishesDuringDay[1] +
		fishesDuringDay[2] + fishesDuringDay[3] + fishesDuringDay[4] + fishesDuringDay[5] + fishesDuringDay[6], nil
}
