/**
* @Author:zhoutao
* @Date:2022/2/10 14:01
* @Desc:
 */

package main

import (
	"time"

	"github.com/blurooo/go-monitor"
)

// 注册一个上报客户端,用于http服务质量监控

var httpReportClient = monitor.Register(monitor.ReportClientConfig{
	Name:             "http服务监控",
	StatisticalCycle: 1000, // 每100ms统计一次服务质量
	/*	OutputCaller: func(o *monitor.OutPutData) {
		// 写入数据库等逻辑

	},*/
})

func main() {
	t := time.NewTicker(10 * time.Millisecond)
	for curTime := range t.C {
		// 每10ms向http监控客户端上报一条http服务数据，耗时0-100ms，状态为200
		httpReportClient.Report("GET - /app/api/users", uint32(curTime.Nanosecond()%100), 200)
	}
}

/**
go-monitor将每个统计周期(100ms，默认1min)输出一条服务质量分析报告，例如：
{"timestamp":"2018-01-24T09:10:55.190503145Z","clientName":"http服务监控","interfaceName":"GET - /app/api/users","count":10,"successCount":10,"successRate":1,"successMsAver":48,"maxMs":98,"minMs":9,"fastCount":10,"fastRate":1,"failCount":0,"failDistribution":{},"timeConsumingDistribution":{"100~150":0,"150~200":0,"200~250":0,"250~300":0,"300~350":0,"350~400":0,"400~450":0,"450~500":0,"<100":10,">500":0}}

默认的报告数据将输出在控制台，但允许我们定制，例如打印到日志文件或写入数据库等，只需传入我们自己的OutputCaller即可：
var httpReportClient = monitor.Register(monitor.ReportClientConfig {
    Name: "http服务监控",
    StatisticalCycle: 100,  // 每100ms统计一次服务质量
    OutputCaller: func(o *monitor.OutPutData) {
        // 写入数据库等逻辑
        ...
    },
})

*/
