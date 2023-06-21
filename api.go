package wx

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/zohu/zlog"
)

/**
接口定义
*/

const (
	ApiCgiBin     = "https://api.weixin.qq.com/cgi-bin"
	ApiMpCgiBin   = "https://mp.weixin.qq.com/cgi-bin"
	ApiWorkCgiBin = "https://qyapi.weixin.qq.com/cgi-bin"

	ApiWxa    = "https://api.weixin.qq.com/wxa"
	ApiWxaapi = "https://api.weixin.qq.com/wxaapi"
	ApiSns    = "https://api.weixin.qq.com/sns"
)

// HTTP CLIENT
func debug() gout.DebugFunc {
	return func(o *gout.DebugOption) {
		o.Debug = true
		o.Write = zlog.SafeWriter()
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
