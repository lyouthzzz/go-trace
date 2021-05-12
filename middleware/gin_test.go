package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gotrace "github.com/lyouthzzz/go-trace"

	"github.com/stretchr/testify/require"
)

func TestGinServerInterceptor(t *testing.T) {
	engine := gin.New()
	engine.Use(ServerTracingInterceptor())
	engine.GET("/ping", func(c *gin.Context) {
		sc := gotrace.SpanContextFromContext(c.Request.Context())
		fmt.Println("global ticket id: " + sc.GetGlobalTicketId())
		fmt.Println("monitor id: " + sc.GetMonitorId())
		fmt.Println("parent rpc id: " + sc.GetParentRpcId())
		fmt.Println("rpc id: " + sc.GetRpcId())
		c.String(200, "pong")
	})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := engine.Run(":8000")
		require.NoError(t, err)
	}()
	time.Sleep(time.Second)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/ping", nil)
	req.Header.Set(gotrace.GloblaTicketIdKey, "theid")

	_, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	wg.Wait()
}
