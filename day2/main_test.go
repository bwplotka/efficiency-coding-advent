package day2

import (
	"strings"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestDivePart1(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int
	}{
		{
			input: `forward 5
down 5
forward 8
up 3
down 8
forward 2`,
			expected: 150,
		},
		{
			input:    ReadTestInput(t),
			expected: 1868935,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("DivePart1", func(t *testing.T) {
				ans, err := DivePart1(input)
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
			input: `forward 5
down 5
forward 8
up 3
down 8
forward 2`,
			expected: 900,
		},
		{
			input:    ReadTestInput(t),
			expected: 1965970888,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("DivePart2", func(t *testing.T) {
				ans, err := DivePart2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("DivePart2_V2", func(t *testing.T) {
				ans, err := DivePart2_V2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("DivePart2_V3", func(t *testing.T) {
				ans, err := DivePart2_V3(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("DivePart2_V4", func(t *testing.T) {
				ans, err := DivePart2_V4(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("DivePart2_V5", func(t *testing.T) {
				ans, err := DivePart2_V5(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

var Answer int

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkDivePart2(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = DivePart2_V5(input)
	}
}
