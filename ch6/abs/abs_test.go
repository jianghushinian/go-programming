package abs

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"
)

func TestAbs(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}

	got = Abs(2)
	if got != 2 {
		t.Errorf("Abs(2) = %f; want 2", got)
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	tb.Log("testing setup")
	return func(tb testing.TB) {
		tb.Log("testing teardown")
	}
}

func TestAbs_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{
			name: "positive",
			x:    2,
			want: 2,
		},
		{
			name: "negative",
			x:    -3,
			want: 3,
		},
	}
	for _, tt := range tests {
		if got := Abs(tt.x); got != tt.want {
			t.Errorf("Abs(%f) = %v, want %v", tt.x, got, tt.want)
		}
	}
}

func TestAbs_TableDrivenWithSubtests(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{
			name: "positive",
			x:    2,
			want: 2,
		},
		{
			name: "negative",
			x:    -3,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%f) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func setup() {
	fmt.Println("TestMain setup")
}

func teardown() {
	fmt.Println("TestMain teardown")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

func BenchmarkAbsResetTimer(b *testing.B) {
	time.Sleep(100 * time.Millisecond) // 耗时的准备工作
	b.ResetTimer()                     // 重置计时器，之前的耗时不计入结果
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

func BenchmarkAbsStopTimerStartTimer(b *testing.B) {
	b.StopTimer()                      // 暂停计时
	time.Sleep(100 * time.Millisecond) // 耗时的准备工作
	b.StartTimer()                     // 恢复计时
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}

func BenchmarkAbsParallel(b *testing.B) {
	b.SetParallelism(2) // 设置并发度为 2 * GOMAXPROCS
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Abs(-1)
		}
	})
}

func ExampleAbs() {
	fmt.Println(Abs(-1))
	fmt.Println(Abs(2))
	// Output:
	// 1
	// 2
}

func FuzzAbs(f *testing.F) {
	// 添加种子语料（初始测试用例）
	cases := []float64{0, 1, -1, math.MaxFloat64, -math.MaxFloat64}
	for _, tc := range cases {
		f.Add(tc)
	}

	// 定义模糊测试逻辑
	f.Fuzz(func(t *testing.T, x float64) {
		result := Abs(x)

		// 基本属性验证：结果应该始终为非负数
		if result < 0 {
			t.Errorf("Abs(%f) = %f, expected non-negative result", x, result)
		}
		// 验证对称性：Abs(x) 应该等于 Abs(-x)
		if result != Abs(-x) {
			t.Errorf("Abs(%f) = %f != Abs(-%f) = %f", x, result, x, Abs(-x))
		}
	})
}
