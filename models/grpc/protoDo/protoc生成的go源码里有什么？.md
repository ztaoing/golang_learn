
protoc将message解析为struct对象；为客户端创建client；为服务端生成server；

为客户端创建client：通过调用invoke方法向服务端发起请求。call id 就是传递的方法地址：如/Greeter/SayHello

为服务端创建server：为服务端生成接口，并将生成的接口注册到stub中。