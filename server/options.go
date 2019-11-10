package server

type Options struct {
	Addr  string
	SrvID uint32
}

type Option func(o *Options)

func newOptions() *Options {
	return &Options{
		Addr: ":8080",
	}
}

func WithLocalAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithSrvID(id uint32) Option {
	return func(o *Options) {
		o.SrvID = id
	}
}
