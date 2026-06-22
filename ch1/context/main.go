package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ========== ① root: backgroundCtx ==========
	ctx1 := context.Background()

	// ========== ② cancelCtx ==========
	ctx2, cancel2 := context.WithCancel(ctx1)

	// ========== ③ valueCtx (k = v1) ==========
	ctx3 := context.WithValue(ctx1, "k", "v1")

	// ========== ④ cancelCtx (child of ②) ==========
	ctx4, cancel4 := context.WithCancel(ctx2)

	// ========== ⑤ withoutCancelCtx (child of ②) ==========
	ctx5 := context.WithoutCancel(ctx2)

	// ========== ⑥ timerCtx (child of ③) ==========
	ctx6, cancel6 := context.WithTimeout(ctx3, 10*time.Second)

	// ========== ⑦ cancelCtx (child of ③) ==========
	ctx7, cancel7 := context.WithCancel(ctx3)

	// ========== ⑧ cancelCtx (child of ⑤) ==========
	ctx8, cancel8 := context.WithCancel(ctx5)

	// ========== ⑨ valueCtx (k = v2, child of ⑥) ==========
	ctx9 := context.WithValue(ctx6, "k", "v2")

	// ========== ⑩ withoutCancelCtx (child of ⑨) ==========
	ctx10 := context.WithoutCancel(ctx9)

	// ================== 验证：数据流（Value 查找） ==================
	fmt.Println("Value lookup:")
	fmt.Println("ctx6 Value(k): ", ctx6.Value("k"))  // v1
	fmt.Println("ctx10 Value(k):", ctx10.Value("k")) // v2

	// ================== 验证：控制流（取消传播） ==================
	cancel2() // 取消 ② cancelCtx

	time.Sleep(10 * time.Millisecond)

	fmt.Println("\nCancellation status after cancel2():")

	fmt.Println("ctx4 Done:", isCanceled(ctx4)) // true（级联取消）
	fmt.Println("ctx8 Done:", isCanceled(ctx8)) // false（被 withoutCancelCtx 隔离）

	// 手动取消
	cancel4()
	cancel6()
	cancel7()
	cancel8()

	// 防止编译器未使用变量优化
	_ = ctx7
}

// 辅助函数：判断 context 是否被取消
func isCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
