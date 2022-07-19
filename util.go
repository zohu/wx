package wx

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/gtls/structs"
	"time"
)

func StructToMap(data interface{}) map[string]interface{} {
	return structs.Map(data)
}

// HTTP CLIENT
func debug() gout.DebugFunc {
	return func(o *gout.DebugOption) {
		o.Debug = true
		o.Write = log.SafeWriter()
	}
}
func (w *Wechat) Post(url string) *dataflow.DataFlow {
	if w.debug {
		return gout.POST(url).Debug(debug())
	}
	return gout.POST(url)
}
func (w *Wechat) Get(url string) *dataflow.DataFlow {
	if w.debug {
		return gout.GET(url).Debug(debug())
	}
	return gout.GET(url)
}

// Time 复制 time.Time 对象，并返回复制体的指针
func Time(t time.Time) *time.Time {
	return &t
}

// String 复制 string 对象，并返回复制体的指针
func String(s string) *string {
	return &s
}

// Bool 复制 bool 对象，并返回复制体的指针
func Bool(b bool) *bool {
	return &b
}

// Float64 复制 float64 对象，并返回复制体的指针
func Float64(f float64) *float64 {
	return &f
}

// Float32 复制 float32 对象，并返回复制体的指针
func Float32(f float32) *float32 {
	return &f
}

// Int64 复制 int64 对象，并返回复制体的指针
func Int64(i int64) *int64 {
	return &i
}

// Int32 复制 int64 对象，并返回复制体的指针
func Int32(i int32) *int32 {
	return &i
}

func RandStr() {

}
