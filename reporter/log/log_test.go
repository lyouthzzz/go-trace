package log

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lyouthzzz/go-trace/reporter/log/rotate"
	"go.uber.org/zap"
)

func TestRotateLog(t *testing.T) {
	now := time.Now()
	rw := rotate.NewWriter(fmt.Sprintf("/Users/y.liu/go/src/github.com/lyouthzzz/go-trace/log/%04d%02d%02d%02d%02d%02d.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()), rotate.MaxSizeOption(1))

	logger := NewMultiWriter(DebugLevel, rw, os.Stdout)
	catLog := logger.With(zap.String("cate", "A"))
	count := 0
	for {
		catLog.Info("example", zap.Int("count", count), zap.Error(errors.New("I m Error")))
		count++
		if count > 100000 {
			break
		}
	}

}
