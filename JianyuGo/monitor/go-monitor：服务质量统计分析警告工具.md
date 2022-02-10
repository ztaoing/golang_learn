[go-monitor基于golang开发，是一个轻量的，用于服务质量监控并实现分析告警的工具。]

go-monitor目前并不是一个独立的服务，而是希望大多数基于golang开发的项目,如同引入一个日志组件一样使用。

[go-monitor能做什么]

通过上报接口、函数、或者是任意调用服务的耗时以及其状态(成功)，
go-monitor将按照设定的周期自动进行服务质量分析，统计，并输出详细的报告数据。

在服务质量达不到理想状态时，go-monitor将触发告警，并在服务质量回升时，触发恢复通知。

go-monitor提供非常多灵活的配置，以使其在大多数场景下都可以通过参数调整来胜任服务监控的职责。

go-monitor采用无锁队列的方式避免并发锁带来的性能问题，
MBP2012版本实测500万次上报数据(go test bench)仅花费1.6s即完成所有分析统计（此前并发锁方案为1.9s），
性能强大,允许你像记录日志一样来使用它，并且不需要担心IO压力
（大部分日志组件使用缓存写盘的方式提升性能，大并发下IO压力明显）。

[什么场景建议使用go-monitor]
例如我们开发了一个web应用以对外提供服务，我们可以嵌入go-monitor，
上报每一个访问的耗时、状态，以达到对我们整个web应用服务质量的监控，
也可以在服务质量下滑甚至不可用时及时作出告警，更详尽的，我们可以上报任何一个调用服务的状态，
例如我们所访问的数据库，所依赖的外部接口等，除了监控服务质量，

事实上也可以通过go-monitor提供的统计数据了解任何一个服务的平均时延，大到一个完整的接口，小到一个数据库查询语句。
而使用go-monitor的成本非常小，仅仅是在golang项目中引入go-monitor，像使用日志组件一样，毫无负担。

[使用方法]
    
    安装
    go get github.com/blurooo/go-monitor


go-monitor支持多实例，并鼓励使用多实例。实例之间互不影响，例如在同个应用下，
我们除了可以注册一个http服务监控之外，还可以注册一个函数耗时监控：

// 注册一个上报客户端,用于http服务质量监控
var httpReportClient = monitor.Register(monitor.ReportClientConfig {
    Name: "http服务监控",
})

// 注册一个上报客户端,用于函数耗时监控
var funcReportClient = monitor.Register(monitor.ReportClientConfig {
    Name: "函数耗时监控",
})

[告警]

go-monitor除了分析统计之外，还帮助实现告警策略，这依赖于服务异常的判定规则。
默认当上报code为200时，认为成功。
当然，在大多数应用中，如此简单的判定规则通常难以胜任各类复杂的场景。
所以go-monitor允许我们使用白名单的方式定制自己的一套规则：

    // 注册得到一个上报客户端用于http服务质量监控
    var httpReportClient = monitor.Register(monitor.ReportClientConfig {
    Name: "http服务监控",
    CodeFeatureMap: map[int]monitor.CodeFeature {
        0: {
        Success: true,
        Name: "成功",
    },
    10000: {
        Success: false,
        Name: "服务不可用",
    },
    }
    })

CodeFeatureMap中允许声明该状态码是否成功，并指定其名称(使用在统计报告中)，除此之外的code都将认为失败。
除了使用白名单机制来决断code之外，
go-monitor也提供了一个适应性更强的方式去判定（优先于CodeFeatureMap）：

    // 注册得到一个上报客户端用于http服务质量监控
    var httpReportClient = monitor.Register(monitor.ReportClientConfig {
        Name: "http服务监控",
        GetCodeFeature: func(code int) (success bool, name string) {
            if code == 0 {
                 return true, "成功"
            } else {
                 return false, "失败"
            }
        },
    })

在每个统计周期内，成功率达不到期望的值时，该条目将被标记，在连续标记若干个统计周期之后，
go-monitor便会触发成功率不达标告警，告警数据明确指明了具体的监控服务和告警条目，
并附带连续被标记为成功率不达标的几次统计数据，默认打印到控制台，但同样允许我们定制，
我们可以按照自己的意愿处理，例如发送邮件通知相关人等：

    // 注册一个上报客户端,用于http服务质量监控
    var httpReportClient = monitor.Register(monitor.ReportClientConfig {
        Name: "http服务监控",
        AlertCaller: func(clientName string, interfaceName string, alertType monitor.AlertType, recentOutputData []monitor.OutPutData) {
            // 处理相关告警
        }
    })

除了"成功率"不达标告警，go-monitor也提供了"耗时"不达标告警，精确到每个监控条目都允许定制耗时达标参数。

    // 一个上报客户端,全局的耗时达标值
    var httpReportClient = monitor.Register(monitor.ReportClientConfig {
        Name: "http服务监控",
        DefaultFastTime: 1000,   // 设定http上报客户端的默认耗时达标为1000ms内
    })

    // 具体到一个条目的耗时标准
        httpReportClient.AddEntryConfig("GET - /app/api/users", monitor.EntryConfig {
        FastLessThan: 100,   // 设定接口"GET - /app/api/users"的耗时达标值为100ms以内
    })
---

[恢复]

go-monitor同时也支持服务质量恢复通知，与告警的策略类似，当出现告警状态时，
后续若干次连续标记为服务达标的统计数据将触发恢复通知(和熔断器采用了相同的机制)，我们只需要定制RecoverCaller即可：

    var httpReportClient = monitor.Register(monitor.ReportClientConfig {
        Name: "http服务监控",
        RecoverCaller: func(clientName string, interfaceName string, alertType monitor.AlertType, recentOutputData []monitor.OutPutData) {
        // 处理恢复通知
        },
    })



还有更多灵活的配置在go-monitor中得到支持，欢迎大家在使用中发现它们，
更欢迎有意向的开发人参与到这份工作来，在设想中，希望go-monitor可以脱胎为一个完善的独立服务，
以支持任何系统接入（包括前后端上报），并提供尽可能多的现成方案，例如统计数据输出到数据库，邮箱告警，接口通知等。
在此抛砖引玉了：https://github.com/blurooo/go-monitor