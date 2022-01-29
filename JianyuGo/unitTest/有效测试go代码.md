[单元测试]
单元测试（UT）就是检验功能单元是否合格的工具。

发布的流程：
    写代码--》单元测试--》改bug--》集成测试--》预发布测试--》发布--》结束

编写UT代码存在难点：
    1、代码耦合度高，缺少必要的抽象与拆分，以至于不知道如何写UT
    2、存在第三方依赖，例如依赖数据库连接、HTTP请求、数据缓存等
课件，编写可测试代码的难点在于：解耦和依赖

对于难点1（高耦合）：我们需要面向接口编程。接口是对一类对象的抽象性描述，表明对象能够提供什么样的服务，它最主要的作用就是解耦调用者和实现着，
        这成为了可测试代码的关键
对于难点2（依赖）：可以通过mock测试来解决。mock测试就是在测试过程中，对于某些不容易构建或不容易获取的对象，用一个虚拟的对象来创建，以便测试。

【测试工具】
1、自带测试库testing test->math.go

2、断言库：testify:具有常见断言和mock工具链，最重要的是，它能与内置库testing很好地配合使用，：https://github.com/stretchr/testify
还要引入：github.com/stretchr/testify/assert

需要：     将if语句改为assert的形式
            if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
改为：
            assert.Equal(t,Add(tt.args.x,tt.args.y),tt.want,tt.name)

testify可以断言的类型非常丰富，例如断言equal、nil、type、断言两个指针是否指向同一个对象、断言包含、断言子集等

【接口mock框架：gomock】：https://github.com/golang/mock


[常见第三方mock依赖库]
1、go-sqlmock ：https://github.com/DATA-DOG/go-sqlmock :无需真正地数据库连接，就能够在测试中模拟sql驱动程序行为
2、httpmock：用于模拟外部资源的http响应，它使用模式匹配的方式匹配HTTP请求的URL，在匹配到特性的请求时就会返回预先设置好的响应。
    ：https://github.com/jarcoal/httpmock
3、gripmock：用于模拟grpc服务的服务器，通过使用.proto文件生成对grpc服务的实现，其项目地址为：https://github.com/tokopedia/gripmock
4、redismock：用于测试与redis服务器的交互：https://github.comelliotchance/redismock