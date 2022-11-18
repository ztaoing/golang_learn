/**
* @Author:zhoutao
* @Date:2022/3/10 16:32
* @Desc:
 */

package main

import (
	"log"
	"net/http"
	"time"
)

var addr = "8080"

func mains() {
	// 创建路由
	/**
	func NewServeMux() *ServeMux { return new(ServeMux) }

	type ServeMux struct {
		mu    sync.RWMutex
		m     map[string]muxEntry // 回调的map
		es    []muxEntry // slice of entries sorted from longest to shortest.
		hosts bool       // whether any patterns contain hostnames
	}

	type muxEntry struct {
		h       Handler // 需要注册进入的方法,他是一个接口,所以需要注册进来的handler实现了ServeHTTP(ResponseWriter, *Request)
		pattern string
	}

	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
	*/
	// mux是实现了Handler接口
	mux := http.NewServeMux()

	// 设置路由规则，sayBye是一个函数，作为参数
	/**
	func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
		if handler == nil {
			panic("http: nil handler")
		}
		// 把handler注入到ServeMux的m中
		mux.Handle(pattern, HandlerFunc(handler))
	}

	func (mux *ServeMux) Handle(pattern string, handler Handler) {
		mux.mu.Lock()
		defer mux.mu.Unlock()

		if pattern == "" {
			panic("http: invalid pattern")
		}
		if handler == nil {
			panic("http: nil handler")
		}
		if _, exist := mux.m[pattern]; exist {
			panic("http: multiple registrations for " + pattern)
		}
		// 先验证map是否为空
		if mux.m == nil {
			mux.m = make(map[string]muxEntry)
		}
		// 将传入的pattern和handler作为一个muxEntry
		e := muxEntry{h: handler, pattern: pattern}
		// 保存到map中
		mux.m[pattern] = e
		if pattern[len(pattern)-1] == '/' {
			mux.es = appendSorted(mux.es, e)
		}

		if pattern[0] != '/' {
			mux.hosts = true
		}
	}
	*/
	mux.HandleFunc("/bye", sayBye)
	// 创建服务器
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux, // 路由
	}
	// 监听端口并提供服务
	log.Println("starting http server at " + addr)
	// 开启服务
	/**
	for {
		// 获取到socket
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		tempDelay = 0
		// new一个http的conn
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew, runHooks) // before Serve can return
		//  处理链接
		go c.serve(connCtx)
	}

	// 处理链接实际调用的是ServeHTTP
	func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
		if r.RequestURI == "*" {
			if r.ProtoAtLeast(1, 1) {
				w.Header().Set("Connection", "close")
			}
			w.WriteHeader(StatusBadRequest)
			return
		}
		// 从注册的路由中使用patten取出handler,先匹配是否存在
		h, _ := mux.Handler(r)
		// 然后处理
		h.ServeHTTP(w, r)
	}


	*/
	log.Fatal(server.ListenAndServe())
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 1)
	w.Write([]byte("bye, this is http server"))
}
