package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var tr *http.Transport

// 初始化http配置
func init() {
	tr = &http.Transport{
		MaxIdleConns: 100,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, time.Second*2) // 设置建立连接超时
			if err != nil {
				return nil, err
			}
			//设置发送接收数据超时
			// 这里的3s超时，其实是在建立连接之后开始算的，而不是从单次调用开始算的超时！
			err = conn.SetDeadline(time.Now().Add(time.Second * 3))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

}
func main() {
	/**
	这段代码看上去没有问题，代码跑起来，也能正常收发消息，但是这段代码跑一段时间就会出现i/o timeout的报错

	问题：
	在生产中发生的现象是，golang服务在发起http调用时，虽然http.transport设置了3s的超时，但是会发生i/o timeout的报错
	但是查看下游服务的时候，发现下游服务其实100ms已经返回了。

	分析：
	明明服务端显示处理耗时才100ms，且客户端超时设置的是3s，怎么就出现超时报错i/o timeout呢？

	推测有两个可能：
	1、因为服务端打印的日志其实只是服务端应用层打印的日志。但客户端应用层发出数据后，中间还经过客户端的传输层，网络层，
	数据链路层和物理层，再经过服务端的物理层，数据链路层，网络层，传输层到服务端的应用层。服务端应用层处耗时100ms，再原路返回。
	那剩下的3s-100ms可能是耗在了整个流程里的各个层上。比如网络不好的情况下，传输层TCP使劲丢包重传之类的原因。

	2、网络没问题，客户端到服务端链路整个收发流程大概耗时就是100ms左右。客户端处理逻辑问题导致超时。

	一般遇到问题，大部分情况下都不会是底层网络的问题，大胆怀疑是自己的问题就对了，不死心就抓个包看下。

	*/
	/** 超时原因
	大家知道HTTP是应用层协议，传输层用的是TCP协议。

	HTTP协议从1.0以前，默认用的是短连接，每次发起请求都会建立TCP连接。收发数据。然后断开连接。

	TCP连接每次都是三次握手。每次断开都要四次挥手。

	其实没必要每次都建立新连接，建立的连接不断开就好了，每次发送数据都复用就好了。

	于是乎，HTTP协议从1.1之后就默认使用长连接。具体相关信息可以看之前的 这篇文章。

	那么golang标准库里也兼容这种实现。

	通过建立一个连接池，针对每个域名建立一个TCP长连接，比如http://baidu.com和http://golang.com 就是两个不同的域名。

	第一次访问http://baidu.com 域名的时候会建立一个连接，用完之后放到空闲连接池里，下次再要访问http://baidu.com 的时候会重新从连接池里把这个连接捞出来复用。

	插个题外话：这也解释了之前这篇文章里最后的疑问，为什么要强调是同一个域名：一个域名会建立一个连接，一个连接对应一个读goroutine和一个写goroutine。
	正因为是同一个域名，所以最后才会泄漏3个goroutine，如果不同域名的话，那就会泄漏 1+2*N 个协程，N就是域名数。

	分析：
	假设第一次请求要100ms，每次请求完http://baidu.com 后都放入连接池中，下次继续复用，重复29次，耗时2900ms。

	第30次请求的时候，连接从建立开始到服务返回前就已经用了3000ms，刚好到设置的3s超时阈值，那么此时客户端就会报超时 i/o timeout 。

	虽然这时候服务端其实才花了100ms，但耐不住前面29次加起来的耗时已经很长。

	也就是说只要通过 http.Transport 设置了 err = conn.SetDeadline(time.Now().Add(time.Second * 3))，
	并且你用了长连接，哪怕服务端处理再快，客户端设置的超时再长，总有一刻，你的程序会报超时错误。

	正确的姿势：
	1、原本是要给每次调用设置一个超时，而不是给整个连接设置超时，
	2、另外上面出现问题的原因是给长连接设置了超时，且长连接会复用。

	基于以上两点：
	var tr *http.Transport
	13
	14func init() {
	15    tr = &http.Transport{
	16        MaxIdleConns: 100,
	17        // 下面的代码被干掉了
	18        //Dial: func(netw, addr string) (net.Conn, error) {
	19        //  conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
	20        //  if err != nil {
	21        //      return nil, err
	22        //  }
	23        //  err = conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送接受数据超时
	24        //  if err != nil {
	25        //      return nil, err
	26        //  }
	27        //  return conn, nil
	28        //},
	29    }
	30}
	31

	func Get(url string) ([]byte, error) {
	34    m := make(map[string]interface{})
	35    data, err := json.Marshal(m)
	36    if err != nil {
	37        return nil, err
	38    }
	39    body := bytes.NewReader(data)
	40    req, _ := http.NewRequest("Get", url, body)
	41    req.Header.Add("content-type", "application/json")
	42
	43    client := &http.Client{
	44        Transport: tr,
	45        Timeout: 3*time.Second,  // 注意：超时加在这里，是每次调用的超时
	46    }
	47    res, err := client.Do(req)
	48    if res != nil {
	49        defer res.Body.Close()
	50    }
	51    if err != nil {
	52        return nil, err
	53    }
	54    resBody, err := ioutil.ReadAll(res.Body)
	55    if err != nil {
	56        return nil, err
	57    }
	58    return resBody, nil
	59}

	看注释会发现，改动的点有两个：
	1、http.Transport里的建立连接时的一些超时设置干掉了。
	2、在发起http请求的时候会场景http.Client，此时加入超时设置，这里的超时就可以理解为单次请求的超时了。同样可以看下注释
	Timeout specifies a time limit for requests "made by this Client".

	到这里，代码就改好了，实际生产中问题也就解决了。

	实例代码里，如果拿去跑的话，其实还会下面的错：
	1Get http://www.baidu.com/: EOF
	这个是因为调用得太猛了，http://www.baidu.com 那边主动断开的连接，可以理解为一个限流措施，目的是为了保护服务器，毕竟每个人都像这么搞，服务器是会炸的。。。

	解决方案很简单，每次HTTP调用中间加个sleep间隔时间就好。

	*/
	for {
		_, err := get("http://www.baidu.com")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func get(url string) ([]byte, error) {
	m := make(map[string]interface{})
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// build request
	body := bytes.NewReader(data)
	r, _ := http.NewRequest("get", url, body)
	r.Header.Add("context-type", "application/json")

	// bulid client
	client := &http.Client{
		Transport: tr,
	}

	// do request
	res, err := client.Do(r)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	// get response
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}
