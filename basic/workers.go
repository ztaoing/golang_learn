package main

import "fmt"

//工作池的goroutine数目
const Numbers = 10

//工作任务
type task struct {
	begin  int
	end    int
	result chan<- int
}

func main() {
	workers := Numbers

	//工作通道
	taskChan := make(chan task, 10)

	//结果通道
	resultChan := make(chan int, 10)

	//worker信号通道
	doneSignal := make(chan struct{}, 10)

	//初始化task的goroutine，计算101个自然数之和
	go InitTask(taskChan, resultChan, 101)

	//分发任务到NUMBER个goroutine池中
	go DistributeTask(taskChan, workers, doneSignal)

	//获取各个goroutine处理完任务的通知，并关闭结果通道
	go CloseResult(doneSignal, resultChan, workers)

	//通过结果通道获取结果并汇总

	sum := ProcessResult(resultChan)
	fmt.Println(sum)
}

//任务处理:计算begin到end的和
//执行结构写入结果chan

func (t *task) do() {
	sum := 0
	for i := t.begin; i <= t.end; i++ {
		sum += i
	}
	//将结果放入chan中
	t.result <- sum
}

//初始化待处理task chan
func InitTask(taskchan chan<- task, r chan<- int, num int) {

	qu := num / 10
	mod := num % 10
	high := qu * 10

	for j := 0; j < qu; j++ {
		b := 10*j + 1
		e := 10 * (j + 1)
		//fmt.Println(b,e)
		tsk := task{
			begin:  b,
			end:    e,
			result: r,
		}
		//fmt.Println(b, e)
		//将单个分组的范围任务发送给taskchan
		taskchan <- tsk
	}

	if mod != 0 {
		tsk := task{
			begin:  high + 1,
			end:    num,
			result: r,
		}
		//fmt.Println(high+1, num)
		taskchan <- tsk
	}

	// 关闭channel
	close(taskchan)
}

//读取task chan ，并分发到worker goroutine处理，总的数量是workers
func DistributeTask(taskchan <-chan task, workers int, done chan<- struct{}) {

	for i := 0; i < workers; i++ {
		go ProcessTask(taskchan, done)
	}

}

//工作goroutine处理具体任务工作，并将处理结果发送到结果chan中
func ProcessTask(taskchan <-chan task, done chan<- struct{}) {

	for t := range taskchan {
		t.do()
	}

	//发送计算完成信号
	done <- struct{}{}
}

//通过done channel同步等待所有工作goroutine的结束，然后关闭结果channel
func CloseResult(done chan struct{}, r chan int, workers int) {
	for i := 0; i < workers; i++ {
		<-done
	}
	close(done)
	close(r)
}

//读取结果，汇总结果
func ProcessResult(resultChan chan int) int {
	sum := 0
	for r := range resultChan {
		//fmt.Println(r)
		sum += r
	}
	return sum
}
