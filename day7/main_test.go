package day7

import (
	"strings"
	"testing"

	"github.com/bwplotka/efficiency-advent-2021/day2"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestMinFuelPart1(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int64
	}{
		{
			input:    `16,1,2,0,4,2,7,1,2,14`,
			expected: 37,
		},
		{
			input:    day2.ReadTestInput(t), // 300 initial
			expected: 340987,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("MinFuelPart1", func(t *testing.T) {
				ans, err := MinFuelPart1(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, int64(ans))
			})
		})
	}
}

func TestMinFuelPart2(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected int64
	}{
		{
			input:    `16,1,2,0,4,2,7,1,2,14`,
			expected: 168,
		},
		{
			input:    day2.ReadTestInput(t), // 300 initial
			expected: 96987874,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("MinFuelPart2", func(t *testing.T) {
				ans, err := MinFuelPart2(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, int64(ans))
			})
			t.Run("MinFuelPart3", func(t *testing.T) {
				ans, err := MinFuelPart3(input)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, int64(ans))
			})
		})
	}
}

/*
var Answer int64

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkSimLanternfish(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = SimLanternfish_V2(input, 80)
	}
}
*/
