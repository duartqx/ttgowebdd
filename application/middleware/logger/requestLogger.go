package logger

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/duartqx/ttgowebdd/application/middleware/logger/colors"
)

type RequestLogger struct {
	method string
	status int
	since  time.Duration
	path   string
}

func NewRequestLoggerBuilder() *RequestLogger {
	return &RequestLogger{}
}

func (rl *RequestLogger) SetMethod(method string) *RequestLogger {
	rl.method = method
	return rl
}

func (rl *RequestLogger) SetPath(path string) *RequestLogger {
	rl.path = path
	return rl
}

func (rl *RequestLogger) SetSince(since time.Duration) *RequestLogger {
	rl.since = since
	return rl
}

func (rl *RequestLogger) SetStatus(status int) *RequestLogger {
	rl.status = status
	return rl
}

func (rl RequestLogger) GetMethod() string {
	return rl.method
}

func (rl RequestLogger) GetStatus() int {
	return rl.status
}

func (rl RequestLogger) GetSince() time.Duration {
	return rl.since
}

func (rl RequestLogger) GetPath() string {
	return rl.path
}

func (rl RequestLogger) Log() {
	log.Println(
		fmt.Sprintf(
			"| %s | %s | %s | %s",
			rl.padAndColor(7, rl.GetMethod()),
			rl.padAndColor(0, rl.GetStatus()),
			rl.pad(12, rl.GetSince()),
			rl.GetPath(),
		),
	)
}

func (rl RequestLogger) LogErr(err interface{}) {
	log.Println(
		fmt.Sprintf(
			"| %s | %s | %s | %s %s",
			rl.padAndColor(7, rl.GetMethod()),
			rl.padAndColor(0, rl.GetStatus()),
			rl.pad(12, rl.GetSince()),
			rl.GetPath(),
			rl.getColored(err),
		),
	)
}

func (rl RequestLogger) pad(padding int, value interface{}) string {
	var (
		v string = fmt.Sprint(value)
		r int    = int(math.Max(float64(padding-len(v)), 0))
	)
	return v + strings.Repeat(" ", r)
}

func (rl RequestLogger) padAndColor(padding int, value interface{}) string {
	if padding > 0 {
		return rl.getColored(rl.pad(padding, value))
	}
	return rl.getColored(value)
}

func (rl RequestLogger) getColored(value interface{}) string {
	var colored string
	switch {
	case rl.status >= 100 && rl.status < 200:
		colored, _ = colors.ColorIt(colors.Cyan, value)
	case rl.status >= 200 && rl.status < 300:
		colored, _ = colors.ColorIt(colors.Green, value)
	case rl.status >= 300 && rl.status < 400:
		colored, _ = colors.ColorIt(colors.Magenta, value)
	case rl.status >= 400 && rl.status < 500:
		colored, _ = colors.ColorIt(colors.Yellow, value)
	default:
		colored, _ = colors.ColorIt(colors.Red, value)
	}
	return colored
}
