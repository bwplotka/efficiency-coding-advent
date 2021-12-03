package day2

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"text/scanner"

	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/pkg/errors"
)

func ReadTestInput(t testing.TB) string {
	f, err := ioutil.ReadFile("input.txt")
	testutil.Ok(t, err)
	return string(f)
}

func DivePart1(input string) (_ int, err error) {
	s := scanner.Scanner{
		Error: func(s *scanner.Scanner, msg string) {
			err = errors.New(msg)
		},
	}
	s.Init(strings.NewReader(input))

	var forward, depth int
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		dir := s.TokenText()

		tok = s.Scan()
		if tok == scanner.EOF {
			return 0, errors.Errorf("expected digit after %v", dir)
		}

		digit, err := strconv.ParseInt(s.TokenText(), 10, 64)
		if err != nil {
			return 0, err
		}

		switch dir {
		case "forward":
			forward += int(digit)
		case "down":
			depth += int(digit)
		case "up":
			depth -= int(digit)
		default:
			return 0, errors.Errorf("unknown direction %v", dir)
		}
	}

	return depth * forward, nil
}

func DivePart2(input string) (_ int, err error) {
	s := scanner.Scanner{
		Error: func(s *scanner.Scanner, msg string) {
			err = errors.New(msg)
		},
	}
	s.Init(strings.NewReader(input))

	var aim, forward, depth int
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		dir := s.TokenText()

		tok = s.Scan()
		if tok == scanner.EOF {
			return 0, errors.Errorf("expected digit after %v", dir)
		}

		digit, err := strconv.ParseInt(s.TokenText(), 10, 64)
		if err != nil {
			return 0, err
		}

		switch dir {
		case "forward":
			forward += int(digit)
			depth += aim * int(digit)
		case "down":
			aim += int(digit)
		case "up":
			aim -= int(digit)
		default:
			return 0, errors.Errorf("unknown direction %v", dir)
		}
	}

	return depth * forward, nil
}

// DivePart2_V2 is optimized version of DivePart2_V1 for latency.
// Scanner is quite bad, let's reimplement using IndexByte as in https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/.
func DivePart2_V2(input string) (_ int, err error) {
	var aim, forward, depth int

	var line string
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]

		// Last char is digit.
		digit, err := strconv.ParseInt(string(line[len(line)-1]), 10, 64)
		if err != nil {
			return 0, err
		}

		switch line[0] {
		case 'f':
			forward += int(digit)
			depth += aim * int(digit)
		case 'd':
			aim += int(digit)
		case 'u':
			aim -= int(digit)
		default:
			return 0, errors.Errorf("unknown direction %v", line)
		}
	}
	return depth * forward, nil
}

// Borrowed and modified from https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/
func parseInt(val string) (intval int64, _ error) {
	factor := int64(1)
	for i := len(val) - 1; i >= 0; i-- {
		c := val[i]

		if c < '0' || c > '9' {
			return 0, errors.Errorf("bad int: %q", val)
		}
		intval += int64(c-'0') * factor
		factor *= 10
	}
	return 0, nil
}

func parseDigit(val uint8) (intval int, _ error) {
	if val < '0' || val > '9' {
		return 0, errors.Errorf("bad int: %q", val)
	}
	return int(val - '0'), nil
}

// DivePart2_V3 is optimized version of DivePart2_V2 for latency.
// mallogc GC and int parsing are main offendant. Let's reduce allocs and optimize int parsing.
// This solution has 0 heap allocs.
func DivePart2_V3(input string) (_ int, err error) {
	var aim, forward, depth int
	var line string
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]

		// Last char is digit.
		digit, err := parseDigit(line[len(line)-1])
		if err != nil {
			return 0, err
		}

		switch line[0] {
		case 'f':
			forward += digit
			depth += aim * digit
		case 'd':
			aim += digit
		case 'u':
			aim -= digit
		default:
			return 0, errors.Errorf("unknown direction %v", line)
		}
	}
	return depth * forward, nil
}

// DivePart2_V4 is optimized version of DivePart2_V3 for latency.
// IndexByte were main offendants.
// This is still SLOWER than index Byte. V3 is the fastest.
func DivePart2_V4(input string) (_ int, err error) {
	var aim, forward, depth, lastNewLineOffset int
	var line string
	for i, c := range input {
		if c != '\n' {
			continue
		}

		line = input[lastNewLineOffset:i]
		lastNewLineOffset = i + 1

		// Last char is digit.
		digit, err := parseDigit(line[len(line)-1])
		if err != nil {
			return 0, err
		}

		switch line[0] {
		case 'f':
			forward += digit
			depth += aim * digit
		case 'd':
			aim += digit
		case 'u':
			aim -= digit
		default:
			return 0, errors.Errorf("unknown direction %v", line)
		}
	}
	return depth * forward, nil
}

// DivePart2_V5 is optimized version of DivePart2_V3 for latency.
// We can optimize the digit parsing further.
func DivePart2_V5(input string) (_ int, err error) {
	var aim, forward, depth int
	var line string
	for len(input) > 0 {
		i := strings.IndexByte(input, '\n')
		if i == -1 {
			break
		}
		line = input[0:i]
		input = input[i+1:]

		// Last char is digit.
		digit := int(line[len(line)-1] - '0')

		switch line[0] {
		case 'f':
			forward += digit
			depth += aim * digit
		case 'd':
			aim += digit
		case 'u':
			aim -= digit
		default:
			return 0, errors.Errorf("unknown direction %v", line)
		}
	}
	return depth * forward, nil
}
