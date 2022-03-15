package wxmp

import (
	"encoding/json"
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/wx"
)

type Context struct {
	wx.Context
}

type Response struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// FindApp
// @Description: 返回微信实体
// @param appid
// @return *Context
func FindApp(appid string) *Context {
	wechat := wx.NewWechat()
	if m := wechat.HGetAll(wx.RdsAppPrefix + appid).Val(); len(m) == 0 {
		log.Errorf("FindApp 微信不存在：%s", appid)
		return nil
	} else {
		app := new(wx.App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{wx.Context{App: app}}
	}
}
