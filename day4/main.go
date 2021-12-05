package day4

import (
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

// BingoPart2_V2 is optimized version of BingoPart2. Main offendant was map access.
// Before optimizing map I did loop fusion, with obviously no effect, since map was main problem (:
// So v2 is almost no faster than v1.
func BingoPart2_V2(input string) (_ int, err error) {
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
		// Rows and columns.
		for i := 0; i < boardSize; i++ {
			rowScore := -1
			colScore := -1
			for j := 0; j < boardSize; j++ {
				// Rows.
				if rowScore > -2 {
					off := b + i*boardSize*3 + j*3
					order, ok := numbersScore[input[off:off+2]]
					if !ok {
						rowScore = -2
					} else if order > rowScore {
						rowScore = order
					}
				}

				// Column.
				if colScore > -2 {
					off := b + j*boardSize*3 + i*3
					order, ok := numbersScore[input[off:off+2]]
					if !ok {
						colScore = -2
					} else if order > colScore {
						colScore = order
					}
				}

				if rowScore+colScore == -4 {
					break
				}
			}

			if rowScore > -1 && rowScore < boardScore {
				boardScore = rowScore
			}
			if colScore > -1 && colScore < boardScore {
				boardScore = colScore
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

const numberOffset = uint8(48)

type lookupArr_V1 []int

func (a *lookupArr_V1) set(number string, order int) {
	if order == 0 {
		order = -1 // This trick allow us to treat 0 as not-found.
	}
	if len(number) == 1 {
		(*a)[number[0]-numberOffset] = order
		return
	}

	(*a)[(number[0]-numberOffset)*10+number[1]-numberOffset] = order
}
func (a *lookupArr_V1) lookup(number string) (int, bool) {
	var val int
	if number[0] == ' ' {
		val = (*a)[number[1]-numberOffset]
	} else {
		val = (*a)[(number[0]-numberOffset)*10+number[1]-numberOffset]
	}

	if val == 0 {
		return 0, false
	}
	if val == -1 {
		return 0, true
	}
	return val, true
}

// BingoPart2_V3 is optimized version of BingoPart2_V2. Main offendant was still map access, so we switched
// from map to array and gained 80% latency improvement and 50% space.
func BingoPart2_V3(input string) (_ int, err error) {
	firstRowLen := strings.IndexByte(input, '\n')
	numbersStr := strings.Split(input[:firstRowLen], ",")

	// Prepare single digits, so they are easier to match.

	// Instead of putting our numbers to map with generic hashing, we could use 2 dimensional array,
	// since we know the digits will be between 48 (0 digit) to 57 (9 digit). By offsetting 48 this would mean
	// a single array of 10*10=100 elements. Since we know how we will do lookup, we can make this array flat further, by
	// multiplying offset for second digit search by 10.
	numbersScore := make(lookupArr_V1, 100)
	for i := range numbersStr {
		numbersScore.set(numbersStr[i], i)
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
		// Rows and columns.
		for i := 0; i < boardSize; i++ {
			rowScore := -1
			colScore := -1
			for j := 0; j < boardSize; j++ {
				// Rows.
				if rowScore > -2 {
					off := b + i*boardSize*3 + j*3
					order, ok := numbersScore.lookup(input[off : off+2])
					if !ok {
						rowScore = -2
					} else if order > rowScore {
						rowScore = order
					}
				}

				// Column.
				if colScore > -2 {
					off := b + j*boardSize*3 + i*3
					order, ok := numbersScore.lookup(input[off : off+2])
					if !ok {
						colScore = -2
					} else if order > colScore {
						colScore = order
					}
				}

				if rowScore+colScore == -4 {
					break
				}
			}

			if rowScore > -1 && rowScore < boardScore {
				boardScore = rowScore
			}
			if colScore > -1 && colScore < boardScore {
				boardScore = colScore
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
			if score, ok := numbersScore.lookup(val); ok && score <= lastWinningBoardScore {
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

type boardStats struct {
	firstCharPos int
	score        int
}

// BingoPart2_V4 is optimized version of BingoPart2_V3.
// Main offendant is number of iterations we have to do, so what about... efficient concurrency?
// NOPE: It's much worse in this implementation.
func BingoPart2_V4(input string) (_ int, err error) {
	runtime.GOMAXPROCS(6) // My machine has 12 cores, let's use at least half of them (when set to 12 it's slower).

	firstRowLen := strings.IndexByte(input, '\n')
	numbersStr := strings.Split(input[:firstRowLen], ",")

	// Prepare single digits, so they are easier to match.

	// Instead of putting our numbers to map with generic hashing, we could use 2 dimensional array,
	// since we know the digits will be between 48 (0 digit) to 57 (9 digit). By offsetting 48 this would mean
	// a single array of 10*10=100 elements. Since we know how we will do lookup, we can make this array flat further, by
	// multiplying offset for second digit search by 10.
	numbersScore := make(lookupArr_V1, 100)
	for i := range numbersStr {
		numbersScore.set(numbersStr[i], i)
	}

	boards := make([]boardStats, (len(input)-firstRowLen+2)/(boardSize*3*boardSize+1)) // Number of boards.

	wg := sync.WaitGroup{}
	for i := 0; i < len(boards); i++ {
		pos := firstRowLen + 2 + i*(boardSize*3*boardSize+1)
		boards[i].firstCharPos = pos
		boards[i].score = math.MaxInt

		wg.Add(1)

		// Go routine calculating the score for each board.
		go func(boardNum int) {
			defer wg.Done()

			// Rows and columns.
			for i := 0; i < boardSize; i++ {
				rowScore := -1
				colScore := -1
				for j := 0; j < boardSize; j++ {
					// Rows.
					if rowScore > -2 {
						off := boards[boardNum].firstCharPos + i*boardSize*3 + j*3
						order, ok := numbersScore.lookup(input[off : off+2])
						if !ok {
							rowScore = -2
						} else if order > rowScore {
							rowScore = order
						}
					}

					// Column.
					if colScore > -2 {
						off := boards[boardNum].firstCharPos + j*boardSize*3 + i*3
						order, ok := numbersScore.lookup(input[off : off+2])
						if !ok {
							colScore = -2
						} else if order > colScore {
							colScore = order
						}
					}

					if rowScore+colScore == -4 {
						break
					}
				}

				if rowScore > -1 && rowScore < boards[boardNum].score {
					boards[boardNum].score = rowScore
				}
				if colScore > -1 && colScore < boards[boardNum].score {
					boards[boardNum].score = colScore
				}
			}
		}(i)
		wg.Wait()
	}

	lastWinningBoardFirstCharPos := -1
	lastWinningBoardScore := -1
	for _, b := range boards {
		if b.score > -1 && b.score > lastWinningBoardScore {
			lastWinningBoardFirstCharPos = b.firstCharPos
			lastWinningBoardScore = b.score
		}
	}

	// Calc unmarked numbers.
	var unmarked int
	for i := lastWinningBoardFirstCharPos; i < lastWinningBoardFirstCharPos+boardSize*3*boardSize; i += boardSize * 3 {
		for j := i; j < i+boardSize*3; j += 3 {
			val := input[j : j+2]
			if score, ok := numbersScore.lookup(val); ok && score <= lastWinningBoardScore {
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

const workers = 5

// BingoPart2_V5 is optimized version of BingoPart2_V4.
// Let's try workers approach.
// NOTE: Seems faster than sequential! (-42.74%)
func BingoPart2_V5(input string) (_ int, err error) {
	runtime.GOMAXPROCS(3) // My machine has 12 cores, let's use at least half of them (when set to 12 it's slower).

	firstRowLen := strings.IndexByte(input, '\n')
	numbersStr := strings.Split(input[:firstRowLen], ",")

	// Prepare single digits, so they are easier to match.

	// Instead of putting our numbers to map with generic hashing, we could use 2 dimensional array,
	// since we know the digits will be between 48 (0 digit) to 57 (9 digit). By offsetting 48 this would mean
	// a single array of 10*10=100 elements. Since we know how we will do lookup, we can make this array flat further, by
	// multiplying offset for second digit search by 10.
	numbersScore := make(lookupArr_V1, 100)
	for i := range numbersStr {
		numbersScore.set(numbersStr[i], i)
	}

	boards := make([]boardStats, (len(input)-firstRowLen+2)/(boardSize*3*boardSize+1)) // Number of boards.

	wg := sync.WaitGroup{}
	for i := 0; i < len(boards); i++ {
		pos := firstRowLen + 2 + i*(boardSize*3*boardSize+1)
		boards[i].firstCharPos = pos
		boards[i].score = math.MaxInt
	}

	shardElems := len(boards) / workers
	if shardElems < 1 {
		shardElems = 1
	}

	for i := 0; i < len(boards); i += shardElems {
		if i+shardElems > len(boards) {
			shardElems = len(boards) - i
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for k := start; k < end; k++ {
				// Rows and columns.
				for i := 0; i < boardSize; i++ {
					rowScore := -1
					colScore := -1
					for j := 0; j < boardSize; j++ {
						// Rows.
						if rowScore > -2 {
							off := boards[k].firstCharPos + i*boardSize*3 + j*3
							order, ok := numbersScore.lookup(input[off : off+2])
							if !ok {
								rowScore = -2
							} else if order > rowScore {
								rowScore = order
							}
						}

						// Column.
						if colScore > -2 {
							off := boards[k].firstCharPos + j*boardSize*3 + i*3
							order, ok := numbersScore.lookup(input[off : off+2])
							if !ok {
								colScore = -2
							} else if order > colScore {
								colScore = order
							}
						}

						if rowScore+colScore == -4 {
							break
						}
					}

					if rowScore > -1 && rowScore < boards[k].score {
						boards[k].score = rowScore
					}
					if colScore > -1 && colScore < boards[k].score {
						boards[k].score = colScore
					}
				}
			}
		}(i, i+shardElems)
	}
	wg.Wait()

	lastWinningBoardFirstCharPos := -1
	lastWinningBoardScore := -1
	for _, b := range boards {
		if b.score > -1 && b.score > lastWinningBoardScore {
			lastWinningBoardFirstCharPos = b.firstCharPos
			lastWinningBoardScore = b.score
		}
	}

	// Calc unmarked numbers.
	var unmarked int
	for i := lastWinningBoardFirstCharPos; i < lastWinningBoardFirstCharPos+boardSize*3*boardSize; i += boardSize * 3 {
		for j := i; j < i+boardSize*3; j += 3 {
			val := input[j : j+2]
			if score, ok := numbersScore.lookup(val); ok && score <= lastWinningBoardScore {
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
