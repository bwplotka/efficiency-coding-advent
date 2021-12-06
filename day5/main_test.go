package day5

import (
	"strings"
	"testing"

	"github.com/bwplotka/efficiency-advent-2021/day2"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestVentsOverlapPart1(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2
2,2 -> 2,0`, // Extra test case over simple case in advent to check vertical line overlaps.
			expected: 7,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 4421, // Not 1786 and not 2714 (previous runs with bugs)(:
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("VentsOverlapPart1", func(t *testing.T) {
				ans, err := VentsOverlapPart1(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

func TestVentsOverlapPart2(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2,
2,2 -> 2,0`, // Extra test case over simple case in advent to check vertical line overlaps.
			expected: 12 + 2,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 2, // Not 19720
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("VentsOverlapPart2", func(t *testing.T) {
				ans, err := VentsOverlapPart2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

var Answer int

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkVentsOverlapPart2(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = VentsOverlapPart2(input)
	}
}
