package day4

import (
	"math"
	"strconv"
	"strings"
)

const boardSize = 5

// BingoPart1 as a bit of brute force solution.
func BingoPart1(input string) (_ int, err error) {
	firstRowLen := strings.IndexByte(input, '\n')
	numbersStr := strings.Split(input[:firstRowLen], ",")

	// Prepare single digits, so they are easier to match.
	// Put all to map, too, so it's easier to search.
	numbersScore := make(map[string]int, len(numbersStr))
	for i := range numbersStr {
		if len(numbersStr[i]) == 1 {
			numbersStr[i] = " " + numbersStr[i]
		}
		numbersScore[numbersStr[i]] = i
	}

	// Index of first char of board.
	var boards []int
	for i := firstRowLen + 2; i < len(input); i += boardSize*3*boardSize + 1 {
		boards = append(boards, i)
	}

	winningBoard := -1
	winningBoardScore := math.MaxInt
	for _, b := range boards {
		// Rows.
		for i := b; i < b+boardSize*3*boardSize; i += boardSize * 3 {
			score := -1
			for j := i; j < i+boardSize*3; j += 3 {
				order, ok := numbersScore[input[j:j+2]]
				if !ok {
					score = -1
					break
				}
				// We are looking for largest score, which indicates the "oldest" number that we will
				// need to wait in order to win with this row.
				if order > score {
					score = order
				}
			}
			// If we have match, check if this board scores maximum
			if score > -1 && score < winningBoardScore {
				winningBoard = b
				winningBoardScore = score
			}
		}

		// Same, but for columns.
		for i := b; i < b+boardSize*3; i += 3 {
			score := -1
			for j := i; j < i+boardSize*3*boardSize; j += boardSize * 3 {
				order, ok := numbersScore[input[j:j+2]]
				if !ok {
					score = -1
					break
				}
				if order > score {
					score = order
				}
			}
			if score > -1 && score < winningBoardScore {
				winningBoard = b
				winningBoardScore = score
			}
		}
	}

	// Calc unmarked numbers.
	var unmarked int
	for i := winningBoard; i < winningBoard+boardSize*3*boardSize; i += boardSize * 3 {
		for j := i; j < i+boardSize*3; j += 3 {
			val := input[j : j+2]
			if score, ok := numbersScore[val]; ok && score <= winningBoardScore {
				continue
			}

			v, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
			if err != nil {
				return 0, err
			}
			unmarked += int(v)
		}
	}

	lastNum, err := strconv.ParseInt(strings.TrimSpace(numbersStr[winningBoardScore]), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(lastNum) * unmarked, nil
}

func BingoPart2(input string) (_ int, err error) {
	firstRowLen := strings.IndexByte(input, '\n')
	numbersStr := strings.Split(input[:firstRowLen], ",")

	// Prepare single digits, so they are easier to match.
	// Put all to map, too, so it's easier to search.
	numbersScore := make(map[string]int, len(numbersStr))
	for i := range numbersStr {
		if len(numbersStr[i]) == 1 {
			numbersStr[i] = " " + numbersStr[i]
		}
		numbersScore[numbersStr[i]] = i
	}

	// Index of first char of board.
	var boards []int
	for i := firstRowLen + 2; i < len(input); i += boardSize*3*boardSize + 1 {
		boards = append(boards, i)
	}

	lastWinningBoard := -1
	lastWinningBoardScore := -1
	for _, b := range boards {
		boardScore := math.MaxInt
		// Rows.
		for i := b; i < b+boardSize*3*boardSize; i += boardSize * 3 {
			score := -1
			for j := i; j < i+boardSize*3; j += 3 {
				order, ok := numbersScore[input[j:j+2]]
				if !ok {
					score = -1
					break
				}
				// We are looking for largest score, which indicates the "oldest" number that we will
				// need to wait in order to win with this row.
				if order > score {
					score = order
				}
			}
			if score > -1 && score < boardScore {
				boardScore = score
			}
		}

		// Same, but for columns.
		for i := b; i < b+boardSize*3; i += 3 {
			score := -1
			for j := i; j < i+boardSize*3*boardSize; j += boardSize * 3 {
				order, ok := numbersScore[input[j:j+2]]
				if !ok {
					score = -1
					break
				}
				if order > score {
					score = order
				}
			}
			if score > -1 && score < boardScore {
				boardScore = score
			}
		}

		if boardScore > -1 && boardScore > lastWinningBoardScore {
			lastWinningBoard = b
			lastWinningBoardScore = boardScore
		}
	}

	// Calc unmarked numbers.
	var unmarked int
	for i := lastWinningBoard; i < lastWinningBoard+boardSize*3*boardSize; i += boardSize * 3 {
		for j := i; j < i+boardSize*3; j += 3 {
			val := input[j : j+2]
			if score, ok := numbersScore[val]; ok && score <= lastWinningBoardScore {
				continue
			}

			v, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
			if err != nil {
				return 0, err
			}
			unmarked += int(v)
		}
	}

	lastNum, err := strconv.ParseInt(strings.TrimSpace(numbersStr[lastWinningBoardScore]), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(lastNum) * unmarked, nil
}

//// ParseInt is a borrowed and modified version from https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/
//func ParseInt(val string) (intval int, _ error) {
//	factor := 1
//	for i := len(val) - 1; i >= 0; i-- {
//		c := val[i]
//
//		if c < '0' || c > '9' {
//			return 0, errors.Errorf("bad int: %q", val)
//		}
//		intval += int(c-'0') * factor
//		factor *= 10
//	}
//	return 0, nil
//}
