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
	//pool中元素的个数
	len() int
	//pool是否为空
	isEmpty() bool
	//插入一个goWorker的指针到pool中
	insert(worker *goWorker) error
	//取出一个任务
	detach() *goWorker
	//取回过期
	retrieveExpiry(duration time.Duration) []*goWorker
	//重置整个pool
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
