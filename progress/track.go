package progress

import (
	"io"
	"net/http"
)

func TrackSlice[T any](s []T, fn func(i int, item T) bool, options ...TaskModelOption) {
	options = append(options, WithTotal(len(s)))
	p := New(NewTaskModel("", options...))
	p.Start()
	defer p.Finish()
	for i, item := range s {
		if !fn(i, item) {
			break
		}
		p.Incr(1)
	}
}

func TrackMap[K comparable, V any](m map[K]V, fn func(k K, v V) bool, options ...TaskModelOption) {
	options = append(options, WithTotal(len(m)))
	p := New(NewTaskModel("", options...))
	p.Start()
	defer p.Finish()
	for k, v := range m {
		if !fn(k, v) {
			break
		}
		p.Incr(1)
	}
}

func TrackReader(r io.Reader, options ...TaskModelOption) *Reader {
	p := New(NewTaskModel("", options...))
	p.Start()
	return &Reader{
		r: r, p: p,
	}
}

func TrackHTTP(description string, t http.RoundTripper) *RoundTripper {
	return &RoundTripper{
		Description: description,
		Transport:   t,
	}
}
