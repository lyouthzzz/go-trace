package gotrace

import "io"

type Reporter interface {
	io.Closer
	Send(span Span)
}
