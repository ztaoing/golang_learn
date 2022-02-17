文章来源于[Go招聘]

---
周末邮件收到一封https://dev.to/ 关于Go文章的推送：7 GitHub projects to make you a better Go Developer💥[1]，

简单看了一下文章内容，主要介绍了作者参与开源的SigNoz的介绍和Gopher在开发过程中可能会用到7个Go库，下面顺便给大家介绍一下：

[关于SigNoz]

SigNoz[2] 帮助开发人员监控应用程序并解决他们部署的应用程序中的问题，它是 DataDog、NewRelic 等的开源替代品。

下面是SigNoz的Web页面展示图。感兴趣的同学可以到该项目仓库查看帮助并尝试使用。主要特性包括：

* 可以查看 p99 延迟、服务错误率、外部 API 调用等指标。
* 也可以查看各个请求链路的详细火焰图来找到问题的根本原因。

SigNoz vs Promethus:
监控相关metrics的话，Prometheus 是不错的选择。但SigNoz 主要是在metrics和trace中间提供一个集成的用户界面——类似于 Datadog 等 SaaS 供应商提供的服务——并提供对trace的高级过滤和聚合，这也是 Jaeger 目前所缺乏的。

SigNoz VS Jaeger:
Jaeger 只做分布式链路追踪。SigNoz 指标和跟踪都会做，另外还有日志管理。
此外，SigNoz 与 Jaeger 相比还有一些更高级的功能：
* Jaeger UI 不显示任何关于trace或过滤trace的metrics
* Jaeger 无法聚合在过滤后的trace。

上面主要介绍了作者参与的SigNoz项目的介绍。接下来主要介绍开发过程中有利于Gopher的7个仓库

---

[Awesome-go]
awesome-go[3] 就不多介绍了，可以称为 Go 编程语言的百科全书。

主要罗列了大量Go 框架、一些库的精选列表。awesome系列的大部分语言栈都有。感兴趣的同学也可以去Github进行搜索。

[Project-layout]
project-layout[4] 相信Gopher同学们也是众所周知了，这个仓库主要包含 Go 应用程序项目的基本布局。

目前为止据我了解很多公司也纷纷采用这种项目布局模式。

前段时间还被Russ Cox吐槽了：this is not a standard Go project layout，想看热闹的同学还可以戳戳Russ Cox 看不下去了：golang-standards/project-layout 不是 Go 标准布局。

[Go-Kit]
Go-kit[5] 是一个用于在 Go 中构建微服务的编程工具库。主要解决了分布式系统和应用程序架构中的常见问题，所以主要让开发者专注于业务交付。

[Go-patterns]
go-patterns[6] 这个 仓库包含了 Go 语言的惯用设计和应用程序模式的精选集合。您可以找到以下模式：创建模式、结构模式、行为模式、并发模式、消息传递模式。

[ Learn-go-with-tests]
learn-go-with-tests[7] Go 是一种学习测试驱动开发的好语言，因为 Go 的标准库提供了一个内置的测试包。这个 仓库有一个 Go 基础知识列表和测试驱动代码实现的例子。

中文版本的可以阅读Go语言中文网出品的通过测试驱动开发来学习Go[8]

[The Ultimate Go Study Guide]
ultimate-go[9] 已转向由gotraining[10] 维护，Ultimate Go 学习指南主要是为参加 Ultimate Go 课程的学生准备的笔记集。

带有逐行注释的示例程序以帮助学生更好地理解代码。另外关于此仓库目前已经收集成册，开卖了电子版，售价9.9美元。更多详情可戳 ultimate-go-notebook[11] 。

[Learngo]
learngo[12] 这个仓库中，你会发现数有以千计的 Go 示例、练习。

---

[参考资料]

[1]
7 GitHub projects to make you a better Go Developer💥: https://dev.to/ankit01oss/7-github-projects-to-make-you-a-better-go-developer-2nmh

[2]
SigNoz: https://github.com/SigNoz/signoz

[3]
awesome-go: https://github.com/avelino/awesome-go

[4]
project-layout: https://github.com/golang-standards/project-layout

[5]
Go-kit: https://github.com/go-kit/kit

[6]
go-patterns: https://github.com/tmrts/go-patterns

[7]
learn-go-with-tests: https://github.com/quii/learn-go-with-tests

[8]
通过测试驱动开发来学习Go: https://studygolang.gitbook.io/learn-go-with-tests/

[9]
ultimate-go: https://github.com/hoanhan101/ultimate-go

[10]
gotraining: https://github.com/ardanlabs/gotraining

[11]
ultimate-go-notebook: https://education.ardanlabs.com/courses/ultimate-go-notebook

[12]
learngo: https://github.com/inancgumus/learngo