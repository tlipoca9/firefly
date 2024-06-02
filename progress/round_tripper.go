package progress

import (
	"net/http"

	"github.com/docker/go-units"
)

type RoundTripper struct {
	Description string
	Transport   http.RoundTripper
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.Transport == nil {
		rt.Transport = http.DefaultTransport
	}
	resp, err := rt.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return resp, err
	}

	resp.Body = TrackReader(
		resp.Body,
		WithDescription(rt.Description),
		WithTotal(resp.ContentLength),
		WithUnitFormatter(func(i int64) string {
			return units.BytesSize(float64(i))
		}),
	)

	return resp, err
}
