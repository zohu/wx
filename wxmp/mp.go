package wxmp

import (
	"encoding/json"
	"github.com/zohu/wx"
)

type Context struct {
	wx.Context
}

// FindApp
// @Description: 返回微信实体
// @param appid
// @return *Context
func FindApp(appid string) (*Context, *wx.Err) {
	wechat := wx.NewWechat()
	if m := wechat.HGetAll(wx.RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, &wx.Err{
			Appid: appid,
			Err:   "not exist",
			Desc:  "公众号不存在",
		}
	} else {
		app := new(wx.App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{wx.Context{App: app}}, nil
	}
}
