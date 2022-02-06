
Go官方提供的文件操作标准库分散在os、ioutil等多个包中，里面有非常多的方法涵盖了文件操作的所有场景，
不过因为我平时开发过程中需要直接操作文件的场景其实并不多，
在加上Go标准库的文档太难搜索，每次遇到要使用文件函数时都是直接Google查对应的函数
。偶然查到国外一个人在2015年写的博客，他用常用的文件函数汇总了30个文件操作场景，
包括四大类：基本操作、读写操作、文件压缩、其他操作。
每一个文件操作都给了代码示例。写的非常好，强烈推荐你阅读一下，浏览一下它的目录，
然后放到收藏夹里吃灰图片图片，万一哪天用到了还能拿来参考一下。

    
    原文链接：https://www.devdungeon.com/content/working-files-go

[一切皆文件]
UNIX 的一个基础设计就是"万物皆文件"(everything is a file)。
我们不必知道操作系统的设备驱动把什么映射给了一个文件描述符，操作系统为设备提供了文件格式的接口。

Go语言中的reader和writer接口也类似。我们只需简单的读写字节，不必知道reader的数据来自哪里，
也不必知道writer将数据发送到哪里。你可以在/dev下查看可用的设备，有些可能需要较高的权限才能访问。

[文件基本操作]
* 创建空文件 -->createFile
* Truncate裁剪文件-->Truncate
* 获取文件信息-->fileInfo
* 重命名和移动-->rename
* 删除文件-->dele
* 打开和关闭文件-->openClose
* 检查文件是否存在-->exist
* 检查读写权限-->permission
* 改变权限、拥有者、时间戳-->change
* 创建硬连接和软连接-->createLink
* 复制文件-->copy
* 跳转到文件指定位置-->jump
* 写文件-->write 写入的模式
* 快写文件-->fastWrite
* 使用缓存写-->writeBuf
* 读取最多N个字节-->readN
* 读取正好N个字节-->readN
* 读取最少N个字节-->readN
* 读取全部字节-->readN
* 快读到内存-->readIntoMemo
* 使用缓存读-->useBuff
* 使用scanner-->useScanner

【文件压缩】
* 打包（zip）文件-->zipPack
* 解压(unzip)文件-->zipUnpack
* 压缩文件-->pack
* 解压缩文件-->unpack

[文件其他操作]
* 临时文件和目录-->tempFile
* 通过HTTP下载文件-->downloadHttp
* 哈希和摘要-->hash