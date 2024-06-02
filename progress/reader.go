package progress

import "io"

type Reader struct {
	r io.Reader
	p *Progress
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil {
		r.p.Finish()
	} else {
		r.p.Incr(n)
	}
	return
}

func (r *Reader) Close() error {
	r.p.Finish()
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
