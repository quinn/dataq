package queue

type Options struct {
	Path string
}

type Option func(*Options)

func WithPath(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}
