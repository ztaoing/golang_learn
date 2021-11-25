
# grpc中是如何处理异常的：grpc的异常码
* 0：Ok：返回成功
* 1：Canceled：操作已取消
* 2：Unknown：未知错误。如果从另一个地址空间接收到的状态值属 于在该地址空间中未知的错误空间，则可以返回此错误的示例。 没有返回足够的错误信息的API引发的错误也可能会转换为此错误
* 3：InvalidArgument：表示客户端指定了无效的参数。 请注意，这与FailedPrecondition不同。 它表示无论系统状态如何（例如格式错误的文件名）都有问题的参数
* 4：DeadlineExceeded：意味着操作在完成之前过期。 对于更改系统状态的操作，即使操作成功完成，也可能会返回此错误。 例如，服务器的成功响应可能会延迟足够的时间以使截止日期到期
* 5：NotFound：表示找不到某个请求的实体（例如文件或目录）
* 6：AlreadyExists：表示尝试创建实体失败，因为已经存在
* 7：PermissionDenied：表示调用者没有执行指定操作的权限。它不能用于因耗尽某些资源而引起的拒绝（使用ResourceExhausted代替这些错误）。如果调用者无法识别，则不能使用它（使用Unauthenticated代替这些错误）
* 8：ResourceExhausted：表示某些资源已耗尽，可能是每个用户的配额，或者整个文件系统空间不足
* 9：FailedPrecondition：表示操作被拒绝，因为系统不处于操作执行所需的状态。
* 10：Aborted：表示操作被中止，通常是由于并发问题（如序列器检查失败，事务异常终止等）造成的。请参阅上面的试金石测试以确定FailedPrecondition，Aborted和Unavailable之间的差异
* 11：OutOfRange：表示操作尝试超过有效范围。
* 12：Unimplemented：该方法未实现
* 13：Internal： 意味着底层系统预期的一些不变量已被打破。 如果你看到其中的一个错误，那么事情就会非常糟糕
* 14：Unavailable：内部Grpc服务不可用，请求不到
* 15：DataLoss：指示不可恢复的数据丢失或损坏
* 16：Unauthenticated：表示请求没有有效的操作认证凭证

go语言推荐的是一个error和一个正常的信息：
如果有错误：就返回一个error和nil
如果没有错误：就返回一个逆流和正常信息

grpc处理了不同语言之间的错误差异性：go的grpc客户端调用pathoy的grpc服务端
