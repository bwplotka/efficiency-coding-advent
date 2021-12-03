package day3

import (
	"strings"
	"testing"

	"github.com/bwplotka/efficiency-advent-2021/day2"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestBinDiagnosePart1(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010`,
			expected: 198,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 3429254,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("BinDiagnosePart1", func(t *testing.T) {
				ans, err := BinDiagnosePart1(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

func TestDivePart2(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010`,
			expected: 230,
		},
		{
			input:    day2.ReadTestInput(t),
			expected: 5410338,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("BinDiagnosePart2", func(t *testing.T) {
				ans, err := BinDiagnosePart2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("BinDiagnosePart2_V2", func(t *testing.T) {
				ans, err := BinDiagnosePart2_V2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

var Answer int

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkBinDiagnosePart2(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = BinDiagnosePart2(input)
	}
}
