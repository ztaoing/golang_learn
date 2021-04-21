package ants

import "time"

// 环形队列：是worker_array接口的一个实现
type loopQueue struct {
	items  []*goWorker //
	expiry []*goWorker //到期的任务
	head   int         //头
	tail   int         //尾
	size   int         //环形队列的容量,当容量为0时，即已经被释放
	isFull bool        //环形队列是否已满
}

func newWorkerLoopQueue(size int) *loopQueue {
	return &loopQueue{
		items: make([]*goWorker, size),
		size:  size,
	}
}

func (wq *loopQueue) len() int {
	if wq.size == 0 {
		return 0
	}

	if wq.head == wq.tail {
		if wq.isFull {
			return wq.size
		}
		return 0
	}

	if wq.tail > wq.head {
		return wq.tail - wq.head
	}

	return wq.size - wq.head + wq.tail
}

func (wq *loopQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}

func (wq *loopQueue) insert(worker *goWorker) error {
	// 当被释放后，环形队列的容量为0
	if wq.size == 0 {
		return errQueueIsReleased
	}

	if wq.isFull {
		return errQueueIsFull
	}
	// 将worker的指针保存起来
	wq.items[wq.tail] = worker
	wq.tail++
	// 已满，将tail的索引更新为第一个索引
	if wq.tail == wq.size {
		wq.tail = 0
	}
	// 已满
	if wq.tail == wq.head {
		wq.isFull = true
	}

	return nil
}

// 从环形队列的第一个位置取出任务
func (wq *loopQueue) detach() *goWorker {
	if wq.isEmpty() {
		return nil
	}

	w := wq.items[wq.head]
	wq.items[wq.head] = nil
	wq.head++
	// 回到环的头部
	if wq.head == wq.size {
		wq.head = 0
	}
	// 更新isFull
	wq.isFull = false

	return w
}

// 回收过期任务
func (wq *loopQueue) retrieveExpiry(duration time.Duration) []*goWorker {
	if wq.isEmpty() {
		return nil
	}
	// 清空过期队列
	wq.expiry = wq.expiry[:0]
	// 过期时间
	expiryTime := time.Now().Add(-duration)
	// 环形队列不为空
	for !wq.isEmpty() {
		// 此任务的recycleTime
		if expiryTime.Before(wq.items[wq.head].recycleTime) {
			break
		}
		// 加入到过期队列中
		wq.expiry = append(wq.expiry, wq.items[wq.head])
		wq.items[wq.head] = nil
		// 后移head
		wq.head++
		if wq.head == wq.size {
			wq.head = 0
		}
		wq.isFull = false
	}

	return wq.expiry
}

func (wq *loopQueue) reset() {
	if wq.isEmpty() {
		return
	}

Releasing:
	if w := wq.detach(); w != nil {
		// 逐个将队列中的任务都设置为nil
		w.task <- nil
		goto Releasing
	}
	// 重置队列的状态
	wq.items = wq.items[:0]
	wq.size = 0
	wq.head = 0
	wq.tail = 0
}
