package day4

import (
	"strings"
	"testing"

	"github.com/bwplotka/efficiency-advent-2021/day2"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestBingoPart1(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7`,
			expected: 4512,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 4662,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := tcase.input + "\n"

			t.Run("BingoPart1", func(t *testing.T) {
				ans, err := BingoPart1(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

func TestBignoPart2(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7`,
			expected: 1924,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 12080,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("BingoPart2", func(t *testing.T) {
				ans, err := BingoPart2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("BingoPart2_V2", func(t *testing.T) {
				ans, err := BingoPart2_V2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("BingoPart2_V3", func(t *testing.T) {
				ans, err := BingoPart2_V3(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("BingoPart2_V4", func(t *testing.T) {
				ans, err := BingoPart2_V4(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

var Answer int

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkBingoPart2(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = BingoPart2_V4(input)
	}
}
