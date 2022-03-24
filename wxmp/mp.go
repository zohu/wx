package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/hhcool/wx"
)

type Context struct {
	wx.Context
}

// FindApp
// @Description: 返回微信实体
// @param appid
// @return *Context
func FindApp(appid string) (*Context, error) {
	wechat := wx.NewWechat()
	if m := wechat.HGetAll(wx.RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, fmt.Errorf("FindApp 微信不存在：%s", appid)
	} else {
		app := new(wx.App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{wx.Context{App: app}}, nil
	}
}
