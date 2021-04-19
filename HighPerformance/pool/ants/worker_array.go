package ants

import (
	"errors"
	"time"
)

var (
	// errQueueIsFull 当前队列已满
	errQueueIsFull = errors.New("the queue is full")

	// 向已经释放的worker queue中插入元素的时候
	errQueueIsReleased = errors.New("the queue length is zero")
)

type workerArray interface {
	len() int
	isEmpty() bool
	insert(worker *goWorker) error
	detach() *goWorker                                 //取出一个任务
	retrieveExpiry(duration time.Duration) []*goWorker //取回过期
	reset()
}

type arrayType int

const (
	stackType     arrayType = 1 << iota //栈
	loopQueueType                       //队列
)

func newWorkerArray(aType arrayType, size int) workerArray {
	switch aType {
	case stackType:
		return newWorkerStack(size)
	case loopQueueType:
		return newWorkerLoopQueue(size)
	default:
		return newWorkerStack(size)
	}
}
