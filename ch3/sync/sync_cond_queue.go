//go:build ignore
// +build ignore

package main

import (
	"sync"
)

type Interface interface {
	Add(item any)                   // 元素入队
	Get() (item any, shutdown bool) // 元素出队
	Len() int                       // 获取队列长度
	ShutDown()                      // 关闭队列
	ShuttingDown() bool             // 队列是否已经关闭
}

type Queue struct {
	// 条件变量
	cond *sync.Cond

	// 队列
	queue []any

	// 队列是否关闭的标识
	shuttingDown bool
}

func New() *Queue {
	return &Queue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue) Add(item any) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.shuttingDown { // 如果队列已经关闭，则直接返回，不再入队
		return
	}

	q.queue = append(q.queue, item) // 入队
	q.cond.Signal()                 // 唤醒一个等待者，通知队列中有数据了
}

func (q *Queue) Get() (item any, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait() // 如果队列为空且未关闭，阻塞等待队列中有数据时被唤醒
	}
	if len(q.queue) == 0 { // 如果此时队列为空，说明队列已经被关闭了
		return nil, true
	}

	// 出队逻辑
	item = q.queue[0]
	q.queue[0] = nil // 主动清除引用，帮助 GC 回收
	q.queue = q.queue[1:]

	return item, false
}

func (q *Queue) ShutDown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.shuttingDown = true // 标记队列关闭
	q.cond.Broadcast()
}

func (q *Queue) ShuttingDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return q.shuttingDown // 返回队列是否关闭标识
}

func (q *Queue) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return len(q.queue) // 返回队列当前长度
}
