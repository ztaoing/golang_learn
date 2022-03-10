/**
* @Author:zhoutao
* @Date:2022/3/10 19:34
* @Desc:
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// 有助于反向代理知识的理解

func main() {
	// 创建连接池对象

	/**

	 */
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 30, // 链接超时
			KeepAlive: time.Second * 30, // 长连接超时时间
		}).DialContext, // 拨号的上下文
		MaxIdleConns:          100,              //最大空闲连接数
		IdleConnTimeout:       time.Second * 90, // 空闲超时时间，如果超时MaxIdleConns就会减1
		TLSHandshakeTimeout:   time.Second * 10, //tls握手超时时间
		ExpectContinueTimeout: time.Second * 1,  //100-continue状态码超时时间
	}
	// 创建客户端
	client := &http.Client{
		Timeout:   time.Second * 30, //请求超时时间
		Transport: transport,
	}
	// 请求数据
	/**
	func (c *Client) Get(url string) (resp *Response, err error) {
		// new一个request
		req, err := NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		return c.Do(req)
	}

	func (c *Client) Do(req *Request) (*Response, error) {
		return c.do(req)
	}

	func (c *Client) do(req *Request) (retres *Response, reterr error) {
		if testHookClientDoResult != nil {
			defer func() { testHookClientDoResult(retres, reterr) }()
		}
		if req.URL == nil {
			req.closeBody()
			return nil, &url.Error{
				Op:  urlErrorOp(req.Method),
				Err: errors.New("http: nil Request.URL"),
			}
		}

		var (
			deadline      = c.deadline()
			reqs          []*Request
			resp          *Response
			copyHeaders   = c.makeHeadersCopier(req)
			reqBodyClosed = false // have we closed the current req.Body?

			// Redirect behavior:
			redirectMethod string
			includeBody    bool
		)
		uerr := func(err error) error {
			// the body may have been closed already by c.send()
			if !reqBodyClosed {
				req.closeBody()
			}
			var urlStr string
			if resp != nil && resp.Request != nil {
				urlStr = stripPassword(resp.Request.URL)
			} else {
				urlStr = stripPassword(req.URL)
			}
			return &url.Error{
				Op:  urlErrorOp(reqs[0].Method),
				URL: urlStr,
				Err: err,
			}
		}
		for {
			// For all but the first request, create the next
			// request hop and replace req.
			if len(reqs) > 0 {
				loc := resp.Header.Get("Location")
				if loc == "" {
					resp.closeBody()
					return nil, uerr(fmt.Errorf("%d response missing Location header", resp.StatusCode))
				}
				u, err := req.URL.Parse(loc)
				if err != nil {
					resp.closeBody()
					return nil, uerr(fmt.Errorf("failed to parse Location header %q: %v", loc, err))
				}
				host := ""
				if req.Host != "" && req.Host != req.URL.Host {
					// If the caller specified a custom Host header and the
					// redirect location is relative, preserve the Host header
					// through the redirect. See issue #22233.
					if u, _ := url.Parse(loc); u != nil && !u.IsAbs() {
						host = req.Host
					}
				}
				ireq := reqs[0]

				// 对发过来的请求进行校验，然后封装一个更为完整的request
				req = &Request{
					Method:   redirectMethod,
					Response: resp,
					URL:      u,
					Header:   make(Header),
					Host:     host,
					Cancel:   ireq.Cancel,
					ctx:      ireq.ctx,
				}
				if includeBody && ireq.GetBody != nil {
					req.Body, err = ireq.GetBody()
					if err != nil {
						resp.closeBody()
						return nil, uerr(err)
					}
					req.ContentLength = ireq.ContentLength
				}

				// Copy original headers before setting the Referer,
				// in case the user set Referer on their first request.
				// If they really want to override, they can do it in
				// their CheckRedirect func.
				copyHeaders(req)

				// Add the Referer header from the most recent
				// request URL to the new one, if it's not https->http:
				if ref := refererForURL(reqs[len(reqs)-1].URL, req.URL); ref != "" {
					req.Header.Set("Referer", ref)
				}
				err = c.checkRedirect(req, reqs)

				// Sentinel error to let users select the
				// previous response, without closing its
				// body. See Issue 10069.
				if err == ErrUseLastResponse {
					return resp, nil
				}

				// Close the previous response's body. But
				// read at least some of the body so if it's
				// small the underlying TCP connection will be
				// re-used. No need to check for errors: if it
				// fails, the Transport won't reuse it anyway.
				const maxBodySlurpSize = 2 << 10
				if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
					io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
				}
				resp.Body.Close()

				if err != nil {
					// Special case for Go 1 compatibility: return both the response
					// and an error if the CheckRedirect function failed.
					// See https://golang.org/issue/3795
					// The resp.Body has already been closed.
					ue := uerr(err)
					ue.(*url.Error).URL = loc
					return resp, ue
				}
			}

			reqs = append(reqs, req)
			var err error
			var didTimeout func() bool
			// 向下游发起实际的请求，
			if resp, didTimeout, err = c.send(req, deadline); err != nil {
				// c.send() always closes req.Body
				reqBodyClosed = true
				if !deadline.IsZero() && didTimeout() {
					err = &httpError{
						err:     err.Error() + " (Client.Timeout exceeded while awaiting headers)",
						timeout: true,
					}
				}
				return nil, uerr(err)
			}

			var shouldRedirect bool
			redirectMethod, shouldRedirect, includeBody = redirectBehavior(req.Method, resp, reqs[0])
			if !shouldRedirect {
				return resp, nil
			}

			req.closeBody()
		}
	}

	// didTimeout is non-nil only if err != nil.
	func (c *Client) send(req *Request, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
		if c.Jar != nil {
			for _, cookie := range c.Jar.Cookies(req.URL) {
				req.AddCookie(cookie)
			}
		}
		// 发送请求，获得response响应
		resp, didTimeout, err = send(req, c.transport(), deadline)
		if err != nil {
			return nil, didTimeout, err
		}
		if c.Jar != nil {
			if rc := resp.Cookies(); len(rc) > 0 {
				c.Jar.SetCookies(req.URL, rc)
			}
		}
		// 返回response
		return resp, nil, nil
	}


	// send issues an HTTP request.
	// Caller should close resp.Body when done reading from it.
	func send(ireq *Request, rt RoundTripper, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
		req := ireq // req is either the original request, or a modified fork

		if rt == nil {
			req.closeBody()
			return nil, alwaysFalse, errors.New("http: no Client.Transport or DefaultTransport")
		}

		if req.URL == nil {
			req.closeBody()
			return nil, alwaysFalse, errors.New("http: nil Request.URL")
		}

		if req.RequestURI != "" {
			req.closeBody()
			return nil, alwaysFalse, errors.New("http: Request.RequestURI can't be set in client requests")
		}

		// forkReq forks req into a shallow clone of ireq the first
		// time it's called.
		forkReq := func() {
			if ireq == req {
				req = new(Request)
				*req = *ireq // shallow clone
			}
		}

		// Most the callers of send (Get, Post, et al) don't need
		// Headers, leaving it uninitialized. We guarantee to the
		// Transport that this has been initialized, though.
		if req.Header == nil {
			forkReq()
			req.Header = make(Header)
		}

		if u := req.URL.User; u != nil && req.Header.Get("Authorization") == "" {
			username := u.Username()
			password, _ := u.Password()
			forkReq()
			req.Header = cloneOrMakeHeader(ireq.Header)
			req.Header.Set("Authorization", "Basic "+basicAuth(username, password))
		}

		if !deadline.IsZero() {
			forkReq()
		}
		stopTimer, didTimeout := setRequestCancel(req, rt, deadline)


		// 向下游请求RoundTrip方法,RoundTrip是RoundTripper接口中的方法；Transport实现了这个方法，即实现了RoundTripper接口
		resp, err = rt.RoundTrip(req)
		if err != nil {
			stopTimer()
			if resp != nil {
				log.Printf("RoundTripper returned a response & error; ignoring response")
			}
			if tlsErr, ok := err.(tls.RecordHeaderError); ok {
				// If we get a bad TLS record header, check to see if the
				// response looks like HTTP and give a more helpful error.
				// See golang.org/issue/11111.
				if string(tlsErr.RecordHeader[:]) == "HTTP/" {
					err = errors.New("http: server gave HTTP response to HTTPS client")
				}
			}
			return nil, didTimeout, err
		}
		if resp == nil {
			return nil, didTimeout, fmt.Errorf("http: RoundTripper implementation (%T) returned a nil *Response with a nil error", rt)
		}
		if resp.Body == nil {
			// The documentation on the Body field says “The http Client and Transport
			// guarantee that Body is always non-nil, even on responses without a body
			// or responses with a zero-length body.” Unfortunately, we didn't document
			// that same constraint for arbitrary RoundTripper implementations, and
			// RoundTripper implementations in the wild (mostly in tests) assume that
			// they can use a nil Body to mean an empty one (similar to Request.Body).
			// (See https://golang.org/issue/38095.)
			//
			// If the ContentLength allows the Body to be empty, fill in an empty one
			// here to ensure that it is non-nil.
			if resp.ContentLength > 0 && req.Method != "HEAD" {
				return nil, didTimeout, fmt.Errorf("http: RoundTripper implementation (%T) returned a *Response with content length %d but a nil Body", rt, resp.ContentLength)
			}
			resp.Body = io.NopCloser(strings.NewReader(""))
		}
		if !deadline.IsZero() {
			resp.Body = &cancelTimerBody{
				stop:          stopTimer,
				rc:            resp.Body,
				reqDidTimeout: didTimeout,
			}
		}
		return resp, nil, nil
	}


	// Transport实现了这个方法，即实现了RoundTripper接口
	func (t *Transport) roundTrip(req *Request) (*Response, error) {
		t.nextProtoOnce.Do(t.onceSetNextProtoDefaults)
		ctx := req.Context()
		trace := httptrace.ContextClientTrace(ctx)

		if req.URL == nil {
			req.closeBody()
			return nil, errors.New("http: nil Request.URL")
		}
		if req.Header == nil {
			req.closeBody()
			return nil, errors.New("http: nil Request.Header")
		}
		scheme := req.URL.Scheme
		isHTTP := scheme == "http" || scheme == "https"
		if isHTTP {
			for k, vv := range req.Header {
				if !httpguts.ValidHeaderFieldName(k) {
					req.closeBody()
					return nil, fmt.Errorf("net/http: invalid header field name %q", k)
				}
				for _, v := range vv {
					if !httpguts.ValidHeaderFieldValue(v) {
						req.closeBody()
						return nil, fmt.Errorf("net/http: invalid header field value %q for key %v", v, k)
					}
				}
			}
		}

		origReq := req
		cancelKey := cancelKey{origReq}
		req = setupRewindBody(req)

		if altRT := t.alternateRoundTripper(req); altRT != nil {
			if resp, err := altRT.RoundTrip(req); err != ErrSkipAltProtocol {
				return resp, err
			}
			var err error
			req, err = rewindBody(req)
			if err != nil {
				return nil, err
			}
		}
		if !isHTTP {
			req.closeBody()
			return nil, badStringError("unsupported protocol scheme", scheme)
		}
		if req.Method != "" && !validMethod(req.Method) {
			req.closeBody()
			return nil, fmt.Errorf("net/http: invalid method %q", req.Method)
		}
		if req.URL.Host == "" {
			req.closeBody()
			return nil, errors.New("http: no Host in request URL")
		}

		for {
			select {
			case <-ctx.Done():
				req.closeBody()
				return nil, ctx.Err()
			default:
			}

			// treq gets modified by roundTrip, so we need to recreate for each retry.
			treq := &transportRequest{Request: req, trace: trace, cancelKey: cancelKey}
			cm, err := t.connectMethodForRequest(treq)
			if err != nil {
				req.closeBody()
				return nil, err
			}

			// Get the cached or newly-created connection to either the
			// host (for http or https), the http proxy, or the http proxy
			// pre-CONNECTed to https server. In any case, we'll be ready
			// to send it requests.

		// 第一步
			pconn, err := t.getConn(treq, cm)
			if err != nil {
				t.setReqCanceler(cancelKey, nil)
				req.closeBody()
				return nil, err
			}

			var resp *Response
			if pconn.alt != nil {
				// HTTP/2 path.
				t.setReqCanceler(cancelKey, nil) // not cancelable with CancelRequest
				resp, err = pconn.alt.RoundTrip(req)
			} else {
				// 第二步：
				resp, err = pconn.roundTrip(treq)
			}
			if err == nil {
				resp.Request = origReq
				return resp, nil
			}

			// Failed. Clean up and determine whether to retry.
			if http2isNoCachedConnError(err) {
				if t.removeIdleConn(pconn) {
					t.decConnsPerHost(pconn.cacheKey)
				}
			} else if !pconn.shouldRetryRequest(req, err) {
				// Issue 16465: return underlying net.Conn.Read error from peek,
				// as we've historically done.
				if e, ok := err.(transportReadFromServerError); ok {
					err = e.err
				}
				return nil, err
			}
			testHookRoundTripRetried()

			// Rewind the body if we're able to.
			req, err = rewindBody(req)
			if err != nil {
				return nil, err
			}
		}
	}

	// 第一步：
	func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (pc *persistConn, err error) {
		req := treq.Request
		trace := treq.trace
		ctx := req.Context()
		if trace != nil && trace.GetConn != nil {
			trace.GetConn(cm.addr())
		}

		w := &wantConn{
			cm:         cm,
			key:        cm.key(),
			ctx:        ctx,
			ready:      make(chan struct{}, 1),
			beforeDial: testHookPrePendingDial,
			afterDial:  testHookPostPendingDial,
		}
		defer func() {
			if err != nil {
				w.cancel(t, err)
			}
		}()

		// Queue for idle connection.  获取空闲连接
		if delivered := t.queueForIdleConn(w); delivered {
			pc := w.pc
			// Trace only for HTTP/1.
			// HTTP/2 calls trace.GotConn itself.
			if pc.alt == nil && trace != nil && trace.GotConn != nil {
				trace.GotConn(pc.gotIdleConnTrace(pc.idleAt))
			}
			// set request canceler to some non-nil function so we
			// can detect whether it was cleared between now and when
			// we enter roundTrip
			t.setReqCanceler(treq.cancelKey, func(error) {})
			return pc, nil
		}

		cancelc := make(chan error, 1)
		t.setReqCanceler(treq.cancelKey, func(err error) { cancelc <- err })

		// Queue for permission to dial.
		t.queueForDial(w)

		// Wait for completion or cancellation.
		select {
		case <-w.ready:
			// Trace success but only for HTTP/1.
			// HTTP/2 calls trace.GotConn itself.
			if w.pc != nil && w.pc.alt == nil && trace != nil && trace.GotConn != nil {
				trace.GotConn(httptrace.GotConnInfo{Conn: w.pc.conn, Reused: w.pc.isReused()})
			}
			if w.err != nil {
				// If the request has been canceled, that's probably
				// what caused w.err; if so, prefer to return the
				// cancellation error (see golang.org/issue/16049).
				select {
				case <-req.Cancel:
					return nil, errRequestCanceledConn
				case <-req.Context().Done():
					return nil, req.Context().Err()
				case err := <-cancelc:
					if err == errRequestCanceled {
						err = errRequestCanceledConn
					}
					return nil, err
				default:
					// return below
				}
			}
			return w.pc, w.err
		case <-req.Cancel:
			return nil, errRequestCanceledConn
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case err := <-cancelc:
			if err == errRequestCanceled {
				err = errRequestCanceledConn
			}
			return nil, err
		}
	}

	// 第二步：持久化连接
	func (pc *persistConn) roundTrip(req *transportRequest) (resp *Response, err error) {
		testHookEnterRoundTrip()
		if !pc.t.replaceReqCanceler(req.cancelKey, pc.cancelRequest) {
			pc.t.putOrCloseIdleConn(pc)
			return nil, errRequestCanceled
		}
		pc.mu.Lock()
		pc.numExpectedResponses++
		headerFn := pc.mutateHeaderFunc
		pc.mu.Unlock()

		if headerFn != nil {
			headerFn(req.extraHeaders())
		}

		// Ask for a compressed version if the caller didn't set their
		// own value for Accept-Encoding. We only attempt to
		// uncompress the gzip stream if we were the layer that
		// requested it.
		requestedGzip := false
		if !pc.t.DisableCompression &&
			req.Header.Get("Accept-Encoding") == "" &&
			req.Header.Get("Range") == "" &&
			req.Method != "HEAD" {
			// Request gzip only, not deflate. Deflate is ambiguous and
			// not as universally supported anyway.
			// See: https://zlib.net/zlib_faq.html#faq39
			//
			// Note that we don't request this for HEAD requests,
			// due to a bug in nginx:
			//   https://trac.nginx.org/nginx/ticket/358
			//   https://golang.org/issue/5522
			//
			// We don't request gzip if the request is for a range, since
			// auto-decoding a portion of a gzipped document will just fail
			// anyway. See https://golang.org/issue/8923
			requestedGzip = true
			req.extraHeaders().Set("Accept-Encoding", "gzip")
		}

		var continueCh chan struct{}
		if req.ProtoAtLeast(1, 1) && req.Body != nil && req.expectsContinue() {
			continueCh = make(chan struct{}, 1)
		}

		if pc.t.DisableKeepAlives &&
			!req.wantsClose() &&
			!isProtocolSwitchHeader(req.Header) {
			req.extraHeaders().Set("Connection", "close")
		}

		gone := make(chan struct{})
		defer close(gone)

		defer func() {
			if err != nil {
				pc.t.setReqCanceler(req.cancelKey, nil)
			}
		}()

		const debugRoundTrip = false

		// Write the request concurrently with waiting for a response,
		// in case the server decides to reply before reading our full
		// request body.
		startBytesWritten := pc.nwrite
		writeErrCh := make(chan error, 1)

		// 向writech channel中写入writeRequest结构体,由谁来读呢？writeLoop监听writech channel
		pc.writech <- writeRequest{req, writeErrCh, continueCh}

		resc := make(chan responseAndError)
		pc.reqch <- requestAndChan{
			req:        req.Request,
			cancelKey:  req.cancelKey,
			ch:         resc,
			addedGzip:  requestedGzip,
			continueCh: continueCh,
			callerGone: gone,
		}

		var respHeaderTimer <-chan time.Time
		cancelChan := req.Request.Cancel
		ctxDoneChan := req.Context().Done()
		pcClosed := pc.closech
		canceled := false
		for {
			testHookWaitResLoop()
			select {
			case err := <-writeErrCh:
				if debugRoundTrip {
					req.logf("writeErrCh resv: %T/%#v", err, err)
				}
				if err != nil {
					pc.close(fmt.Errorf("write error: %v", err))
					return nil, pc.mapRoundTripError(req, startBytesWritten, err)
				}
				if d := pc.t.ResponseHeaderTimeout; d > 0 {
					if debugRoundTrip {
						req.logf("starting timer for %v", d)
					}
					timer := time.NewTimer(d)
					defer timer.Stop() // prevent leaks
					respHeaderTimer = timer.C
				}
			case <-pcClosed:
				pcClosed = nil
				if canceled || pc.t.replaceReqCanceler(req.cancelKey, nil) {
					if debugRoundTrip {
						req.logf("closech recv: %T %#v", pc.closed, pc.closed)
					}
					return nil, pc.mapRoundTripError(req, startBytesWritten, pc.closed)
				}
			case <-respHeaderTimer:
				if debugRoundTrip {
					req.logf("timeout waiting for response headers.")
				}
				pc.close(errTimeout)
				return nil, errTimeout
			case re := <-resc:
				if (re.res == nil) == (re.err == nil) {
					panic(fmt.Sprintf("internal error: exactly one of res or err should be set; nil=%v", re.res == nil))
				}
				if debugRoundTrip {
					req.logf("resc recv: %p, %T/%#v", re.res, re.err, re.err)
				}
				if re.err != nil {
					return nil, pc.mapRoundTripError(req, startBytesWritten, re.err)
				}
				return re.res, nil
			case <-cancelChan:
				canceled = pc.t.cancelRequest(req.cancelKey, errRequestCanceled)
				cancelChan = nil
			case <-ctxDoneChan:
				canceled = pc.t.cancelRequest(req.cancelKey, req.Context().Err())
				cancelChan = nil
				ctxDoneChan = nil
			}
		}
	}


	// 从writech channel中读取writeRequest结构体
	func (pc *persistConn) writeLoop() {
		defer close(pc.writeLoopDone)
		for {
			select {
			case wr := <-pc.writech:
				startBytesWritten := pc.nwrite

				// 如果从writech读取到writeRequest，就像pc.bw写入数据
				err := wr.req.Request.write(pc.bw, pc.isProxy, wr.req.extra, pc.waitForContinue(wr.continueCh))
				if bre, ok := err.(requestBodyReadError); ok {
					err = bre.error
					// Errors reading from the user's
					// Request.Body are high priority.
					// Set it here before sending on the
					// channels below or calling
					// pc.close() which tears down
					// connections and causes other
					// errors.
					wr.req.setError(err)
				}
				if err == nil {
					err = pc.bw.Flush()
				}
				if err != nil {
					if pc.nwrite == startBytesWritten {
						err = nothingWrittenError{err}
					}
				}
				pc.writeErrCh <- err // to the body reader, which might recycle us
				wr.ch <- err         // to the roundTrip function
				if err != nil {
					pc.close(err)
					return
				}
			case <-pc.closech:
				return
			}
		}
	}


	*/

	/**
	连接池：transport流程
	1、客户端发起请求
	2、roundtrip接收到客户端的连接
	3、从缓存池中获取连接，没有空闲连接则新建连接
	4、客户端的请求发送到服务端，并等待响应
	5、从服务端获得response之后发送给客户端
	6、将持久连接放回连接池中


	type Transport struct {
		idleMu       sync.Mutex
		closeIdle    bool                                // user has requested to close all idle conns
		idleConn     map[connectMethodKey][]*persistConn // most recently used at end
		idleConnWait map[connectMethodKey]wantConnQueue  // waiting getConns
		idleLRU      connLRU

		reqMu       sync.Mutex
		reqCanceler map[cancelKey]func(error)

		altMu    sync.Mutex   // guards changing altProto only
		altProto atomic.Value // of nil or map[string]RoundTripper, key is URI scheme

		connsPerHostMu   sync.Mutex
		connsPerHost     map[connectMethodKey]int
		connsPerHostWait map[connectMethodKey]wantConnQueue // waiting getConns

		// Proxy specifies a function to return a proxy for a given
		// Request. If the function returns a non-nil error, the
		// request is aborted with the provided error.
		//
		// The proxy type is determined by the URL scheme. "http",
		// "https", and "socks5" are supported. If the scheme is empty,
		// "http" is assumed.
		//
		// If Proxy is nil or returns a nil *URL, no proxy is used.
		Proxy func(*Request) (*url.URL, error)

		// DialContext specifies the dial function for creating unencrypted TCP connections.
		// If DialContext is nil (and the deprecated Dial below is also nil),
		// then the transport dials using package net.
		//
		// DialContext runs concurrently with calls to RoundTrip.
		// A RoundTrip call that initiates a dial may end up using
		// a connection dialed previously when the earlier connection
		// becomes idle before the later DialContext completes.
		DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

		// Dial specifies the dial function for creating unencrypted TCP connections.
		//
		// Dial runs concurrently with calls to RoundTrip.
		// A RoundTrip call that initiates a dial may end up using
		// a connection dialed previously when the earlier connection
		// becomes idle before the later Dial completes.
		//
		// Deprecated: Use DialContext instead, which allows the transport
		// to cancel dials as soon as they are no longer needed.
		// If both are set, DialContext takes priority.
		Dial func(network, addr string) (net.Conn, error)

		// DialTLSContext specifies an optional dial function for creating
		// TLS connections for non-proxied HTTPS requests.
		//
		// If DialTLSContext is nil (and the deprecated DialTLS below is also nil),
		// DialContext and TLSClientConfig are used.
		//
		// If DialTLSContext is set, the Dial and DialContext hooks are not used for HTTPS
		// requests and the TLSClientConfig and TLSHandshakeTimeout
		// are ignored. The returned net.Conn is assumed to already be
		// past the TLS handshake.
		DialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)

		// DialTLS specifies an optional dial function for creating
		// TLS connections for non-proxied HTTPS requests.
		//
		// Deprecated: Use DialTLSContext instead, which allows the transport
		// to cancel dials as soon as they are no longer needed.
		// If both are set, DialTLSContext takes priority.
		DialTLS func(network, addr string) (net.Conn, error)

		// TLSClientConfig specifies the TLS configuration to use with
		// tls.Client.
		// If nil, the default configuration is used.
		// If non-nil, HTTP/2 support may not be enabled by default.
		TLSClientConfig *tls.Config

		// TLSHandshakeTimeout specifies the maximum amount of time waiting to
		// wait for a TLS handshake. Zero means no timeout.
		TLSHandshakeTimeout time.Duration

		// DisableKeepAlives, if true, disables HTTP keep-alives and
		// will only use the connection to the server for a single
		// HTTP request.
		//
		// This is unrelated to the similarly named TCP keep-alives.
		DisableKeepAlives bool

		// DisableCompression, if true, prevents the Transport from
		// requesting compression with an "Accept-Encoding: gzip"
		// request header when the Request contains no existing
		// Accept-Encoding value. If the Transport requests gzip on
		// its own and gets a gzipped response, it's transparently
		// decoded in the Response.Body. However, if the user
		// explicitly requested gzip it is not automatically
		// uncompressed.
		DisableCompression bool

		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		MaxIdleConns int

		// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
		// (keep-alive) connections to keep per-host. If zero,
		// DefaultMaxIdleConnsPerHost is used.
		MaxIdleConnsPerHost int

		// MaxConnsPerHost optionally limits the total number of
		// connections per host, including connections in the dialing,
		// active, and idle states. On limit violation, dials will block.
		//
		// Zero means no limit.
		MaxConnsPerHost int

		// IdleConnTimeout is the maximum amount of time an idle
		// (keep-alive) connection will remain idle before closing
		// itself.
		// Zero means no limit.
		IdleConnTimeout time.Duration

		// ResponseHeaderTimeout, if non-zero, specifies the amount of
		// time to wait for a server's response headers after fully
		// writing the request (including its body, if any). This
		// time does not include the time to read the response body.
		ResponseHeaderTimeout time.Duration

		// ExpectContinueTimeout, if non-zero, specifies the amount of
		// time to wait for a server's first response headers after fully
		// writing the request headers if the request has an
		// "Expect: 100-continue" header. Zero means no timeout and
		// causes the body to be sent immediately, without
		// waiting for the server to approve.
		// This time does not include the time to send the request header.
		ExpectContinueTimeout time.Duration

		// TLSNextProto specifies how the Transport switches to an
		// alternate protocol (such as HTTP/2) after a TLS ALPN
		// protocol negotiation. If Transport dials an TLS connection
		// with a non-empty protocol name and TLSNextProto contains a
		// map entry for that key (such as "h2"), then the func is
		// called with the request's authority (such as "example.com"
		// or "example.com:1234") and the TLS connection. The function
		// must return a RoundTripper that then handles the request.
		// If TLSNextProto is not nil, HTTP/2 support is not enabled
		// automatically.
		TLSNextProto map[string]func(authority string, c *tls.Conn) RoundTripper

		// ProxyConnectHeader optionally specifies headers to send to
		// proxies during CONNECT requests.
		// To set the header dynamically, see GetProxyConnectHeader.
		ProxyConnectHeader Header

		// GetProxyConnectHeader optionally specifies a func to return
		// headers to send to proxyURL during a CONNECT request to the
		// ip:port target.
		// If it returns an error, the Transport's RoundTrip fails with
		// that error. It can return (nil, nil) to not add headers.
		// If GetProxyConnectHeader is non-nil, ProxyConnectHeader is
		// ignored.
		GetProxyConnectHeader func(ctx context.Context, proxyURL *url.URL, target string) (Header, error)

		// MaxResponseHeaderBytes specifies a limit on how many
		// response bytes are allowed in the server's response
		// header.
		//
		// Zero means to use a default limit.
		MaxResponseHeaderBytes int64

		// WriteBufferSize specifies the size of the write buffer used
		// when writing to the transport.
		// If zero, a default (currently 4KB) is used.
		WriteBufferSize int

		// ReadBufferSize specifies the size of the read buffer used
		// when reading from the transport.
		// If zero, a default (currently 4KB) is used.
		ReadBufferSize int

		// nextProtoOnce guards initialization of TLSNextProto and
		// h2transport (via onceSetNextProtoDefaults)
		nextProtoOnce      sync.Once
		h2transport        h2Transport // non-nil if http2 wired up
		tlsNextProtoWasNil bool        // whether TLSNextProto was nil when the Once fired

		// ForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero
		// Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
		// By default, use of any those fields conservatively disables HTTP/2.
		// To use a custom dialer or TLS config and still attempt HTTP/2
		// upgrades, set this to true.
		ForceAttemptHTTP2 bool
	}
	*/
	resp, err := client.Get("http://127.0.0.1:1210/bye")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 读取内容
	bds, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bds))

}
