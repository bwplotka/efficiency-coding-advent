package day6

import (
	"strings"
	"testing"

	"github.com/bwplotka/efficiency-advent-2021/day2"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestSimLanternfish(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		days     int
		expected int64
	}{
		{
			input:    `3,4,3,1,2`,
			days:     80,
			expected: 5934,
		},
		{
			input:    day2.ReadTestInput(t), // 300 initial
			days:     80,
			expected: 362639,
		},
		{
			input:    `3,4,3,1,2`,
			days:     256,
			expected: 26984457539,
		},
		{
			input:    day2.ReadTestInput(t), // 300 initial
			days:     256,
			expected: 1639854996917,
		},
	} {
		t.Run("", func(t *testing.T) {
			input := strings.TrimSpace(tcase.input) + "\n"

			t.Run("SimLanternfish", func(t *testing.T) {
				if tcase.days > 80 {
					t.Skip("NOPE, too slow")
				}

				ans, err := SimLanternfish(input, tcase.days)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, int64(ans))
			})
			t.Run("SimLanternfish_V2", func(t *testing.T) {
				ans, err := SimLanternfish_V2(input, tcase.days)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("SimLanternfish_V3", func(t *testing.T) {
				ans, err := SimLanternfish_V3(input, tcase.days)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("SimLanternfish_V4", func(t *testing.T) {
				ans, err := SimLanternfish_V4(input, tcase.days)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
			t.Run("SimLanternfish_V5", func(t *testing.T) {
				ans, err := SimLanternfish_V4(input, tcase.days)
				testutil.Ok(t, err)
				testutil.Equals(t, tcase.expected, ans)
			})
		})
	}
}

var Answer int64

// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
func BenchmarkSimLanternfish(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = SimLanternfish_V5(input, 256)
	}
}

//// go test -count 5 -run '^$' -bench . -memprofile=v1.mem.pprof -cpuprofile=v1.cpu.pprof > v1.txt
//func BenchmarkParseInt(b *testing.B) {
//	b.ReportAllocs()
//
//	b.Run("std", func(b *testing.B) {
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			Answer, _ = strconv.ParseInt("1234", 10, 64)
//		}
//	})
//	b.Run("custom", func(b *testing.B) {
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			Answer, _ = ParseInt("1234")
//		}
//	})
//}
