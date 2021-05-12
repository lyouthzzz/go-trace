package kafka

import gotrace "github.com/lyouthzzz/go-trace"

var _ gotrace.Reporter = (*Reporter)(nil)

type Reporter struct {
}

func (r *Reporter) Send(span gotrace.Span) {
	// todo
}

func (r *Reporter) Close() error {
	return nil
}
