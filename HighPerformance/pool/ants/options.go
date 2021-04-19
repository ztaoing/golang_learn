package ants

import "time"

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// pool配置
type Options struct {

	// 用来清理过期的goroutine的过期时间，他会每一个ExpiryDuration扫描一遍所有的workers，然后清理那些超过ExpiryDuration，没有使用的worker
	ExpiryDuration time.Duration

	// 是否在初始化pool的时候，预先申请内存
	PreAlloc bool

	// 允许阻塞在pool.Submit上的最大goroutine的数量
	// 如果是0，则没有这个限制
	MaxBlockingTasks int

	// 当为true的时候，Pool.Submit永远不会被阻塞
	// 当Pool.Submit无法立刻完成的时候会返回ErrPoolOverload的错误
	// 当为true的时候MaxBlockingTasks是不起作用的
	Nonblocking bool

	// PanicHandler是用来处理每一个goroutine的的panic的
	// if nil, panics will be thrown out again from worker goroutines.
	// 当为nil的时候，panic会从goroutine中被再次抛出来
	PanicHandler func(interface{})

	// Logger是一个用来记录日志信息的定制组件，如果没有设置就会使用log包中的默认的日志组件
	Logger Logger
}

// WithOptions 入参是Options结构体
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

// WithExpiryDuration 设置清理goroutine的间隔时间
func WithExpiryDuration(expiryDuration time.Duration) Option {
	return func(opts *Options) {
		opts.ExpiryDuration = expiryDuration
	}
}

// WithPreAlloc 表示是否需要预分配内存
func WithPreAlloc(preAlloc bool) Option {
	return func(opts *Options) {
		opts.PreAlloc = preAlloc
	}
}

// WithMaxBlockingTasks 当达到pool的容量的时候，设置允许阻塞的最大值
func WithMaxBlockingTasks(maxBlockingTasks int) Option {
	return func(opts *Options) {
		opts.MaxBlockingTasks = maxBlockingTasks
	}
}

// WithNonblocking 非阻塞：当没有可用的worker的时候就直接返回nil
func WithNonblocking(nonblocking bool) Option {
	return func(opts *Options) {
		opts.Nonblocking = nonblocking
	}
}

// WithPanicHandler 设置用户处理panic的handler
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

// WithLogger 设置定制的日志组件
func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}
