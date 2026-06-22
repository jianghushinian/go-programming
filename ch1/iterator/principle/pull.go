//go:build ignore
// +build ignore

package main

func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	// 状态变量
	var (
		v         V    // 迭代值
		ok        bool // 值有效性标志
		done      bool // 迭代结束标志
		yieldNext bool // 同步标记（防止乱序调用）
	)

	// 创建迭代协程 G（未立即运行）
	c := newcoro(func(c *coro) {
		// yield 函数：G 协程逻辑
		yield := func(v1 V) bool {
			if done { // 已终止则不再继续
				return false
			}
			if !yieldNext { // 确保执行流程的正确性
				panic("iter.Pull: yield called again before next")
			}
			yieldNext = false
			v, ok = v1, true // 存储当前迭代值
			coroswitch(c)    // 让出（yield）给主协程 F
			return !done     // 返回是否继续迭代
		}

		seq(yield) // 执行原始迭代器逻辑
		var v0 V
		v, ok = v0, false // v、ok 置零
		done = true       // 标记迭代结束
	})

	// next 函数：主协程 F 恢复 G 的执行
	next = func() (v1 V, ok1 bool) {
		if done { // 迭代已结束
			return
		}
		yieldNext = true // 允许 G 执行 yield
		coroswitch(c)    // 恢复（resume）G 的执行
		return v, ok     // 返回 G 通过 yield 传递的值
	}

	// stop 函数：终止迭代
	stop = func() {
		if !done {
			done = true   // 标记终止
			coroswitch(c) // 恢复 G 执行清理
		}
	}
	return next, stop
}
