package gotrace

import (
	"net/http"
	"strings"
)

type (
	Carrier interface {
		Get(key string) string
		Set(key, val string)
	}

	HTTPCarrier http.Header
	GRPCCarrier map[string][]string
)

func (h HTTPCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}

func (h HTTPCarrier) Set(key, val string) {
	http.Header(h).Set(key, val)
}

func (g GRPCCarrier) Get(key string) string {
	if vals, ok := g[strings.ToLower(key)]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (g GRPCCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	g[key] = append(g[key], val)
}
