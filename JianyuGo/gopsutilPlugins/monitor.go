package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/process"
)

func main() {
	// 创建进程对象
	p, _ := process.NewProcess(int32(os.Getpid()))
	// 进程的CPU使用率
	cpuPercent, _ := p.Percent(time.Second)
	cp := cpuPercent / float64(runtime.NumCPU())
	fmt.Println(cp)

}

/**
在容器下获取进程的指标：cgroup限制了容器可以使用资源的情况

*/
