package abs_test

import (
	"testing"

	"github.com/jianghushinian/go-programming/ch6/abs"
)

// abs 黑盒测试
func TestAbs(t *testing.T) {
	got := abs.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
