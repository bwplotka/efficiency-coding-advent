package day6

import (
	"strconv"
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
