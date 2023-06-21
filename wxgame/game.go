package wxgame

import (
	"encoding/json"
	"github.com/zohu/wx"
)

/**
小游戏
*/

type Context struct {
	wx.Context
}

type Response struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// FindApp
// @Description: 返回小游戏实体
// @param appid
// @return *Context
func FindApp(appid string) (*Context, *wx.Err) {
	wechat := wx.NewWechat()
	if m := wechat.HGetAll(wx.RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, &wx.Err{
			Appid: appid,
			Err:   "not exist",
			Desc:  "小游戏不存在",
		}
	} else {
		app := new(wx.App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{wx.Context{App: app}}, nil
	}
}
