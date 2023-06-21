package wxprogram

import (
	"encoding/json"
	"github.com/zohu/wx"
)

/**
小程序
*/

type Context struct {
	wx.Context
}

// FindApp
// @Description: 返回小程序实体
// @param appid
// @return *Context
func FindApp(appid string) (*Context, *wx.Err) {
	wechat := wx.NewWechat()
	if m := wechat.HGetAll(wx.RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, &wx.Err{
			Appid: appid,
			Err:   "not exist",
			Desc:  "小程序不存在",
		}
	} else {
		app := new(wx.App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{wx.Context{App: app}}, nil
	}
}

type ParamCode2Session struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	JsCode    string `json:"js_code"`
	GrantType string `json:"grant_type"`
}
type ResCode2Session struct {
	wx.Response
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

// Code2Session
// @Description: 登录凭证
// @receiver ctx
// @param code
// @return *ResCode2Session
// @return *wx.Err
func (ctx *Context) Code2Session(code string) (*ResCode2Session, *wx.Err) {
	if !ctx.IsMiniApp() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not miniprogram",
			Desc:  "非小程序",
		}
	}

	wechat := wx.NewWechat()
	var res ResCode2Session
	if err := wechat.Get(wx.ApiSns + "/jscode2session").
		SetQuery(&ParamCode2Session{
			Appid:     ctx.Appid(),
			Secret:    ctx.AppSecret(),
			JsCode:    code,
			GrantType: "authorization_code",
		}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "login failure",
			Desc:    "登录失败",
		}
	}
	return &res, nil
}
