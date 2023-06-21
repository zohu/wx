package wxmp

import (
	"github.com/zohu/wx"
)

type QrType int

const (
	QrTypeForever   QrType = 1 // 永久
	QrTypeTemporary QrType = 2 // 临时
)

type ParamNewQrcode struct {
	Expire int64  `json:"expire,omitempty"`
	Type   QrType `json:"type"` // 1永久2临时
	Scene  string `json:"scene"`
}
type ParamQrcode struct {
	ExpireSeconds int64  `json:"expire_seconds,omitempty"` // 二维码有效时间，以秒为单位。 最大不超过2592000（即30天）
	ActionName    string `json:"action_name"`              // 二维码类型，QR_SCENE为临时的整型参数值，QR_STR_SCENE为临时的字符串参数值，QR_LIMIT_SCENE为永久的整型参数值，QR_LIMIT_STR_SCENE为永久的字符串参数值
	ActionInfo    struct {
		Scene struct {
			SceneId  int    `json:"scene_id,omitempty"`  // 场景值ID，临时二维码时为32位非0整型，永久二维码时最大值为100000
			SceneStr string `json:"scene_str,omitempty"` // 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
		} `json:"scene"`
	} `json:"action_info"` // 二维码详情
}

type ResQrcode struct {
	wx.Response
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	Url           string `json:"url"`
}

// Qrcode
// @Description: 生成带参数的二维码
// @receiver ctx
// @param p
// @return string
// @return error
func (ctx *Context) Qrcode(p *ParamNewQrcode) (*ResQrcode, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	var param ParamQrcode
	param.ActionInfo.Scene.SceneStr = p.Scene
	param.ExpireSeconds = p.Expire
	if p.Type == 1 {
		param.ActionName = "QR_LIMIT_STR_SCENE"
	} else {
		param.ActionName = "QR_STR_SCENE"
	}
	wechat := wx.NewWechat()
	var res ResQrcode
	err := wechat.Post(wx.ApiCgiBin + "/qrcode/create").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&param).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.Qrcode(p)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get qrcode",
			Desc:    "获取二维码失败",
		}
	}
	return &res, nil
}

type ParamGenShorten struct {
	LongData      string `json:"long_data"`
	ExpireSeconds int    `json:"expire_seconds,omitempty"`
}
type ResGenShorten struct {
	wx.Response
	ShortKey string `json:"short_key"`
}

// GenShorten
// @Description: 获取短key
// @receiver ctx
// @param data
// @param ex
// @return *ResShortKey
// @return error
func (ctx *Context) GenShorten(data string, ex ...int) (*ResGenShorten, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	param := new(ParamGenShorten)
	param.LongData = data
	if len(ex) > 0 {
		param.ExpireSeconds = ex[0]
	}
	var res ResGenShorten
	if err := wechat.Post(wx.ApiCgiBin + "/shorten/gen").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.GenShorten(data, ex...)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get short key",
			Desc:    "获取短key失败",
		}
	}
	return &res, nil
}

type ParamFetchGenShorten struct {
	ShortKey string `json:"short_key"`
}
type ResFetchGenShorten struct {
	wx.Response
	LongData      string `json:"long_data"`
	CreateTime    int    `json:"create_time"`
	ExpireSeconds int    `json:"expire_seconds"`
}

// FetchGenShorten
// @Description: 还原短key
// @receiver ctx
// @param shortKey
// @return *ResFetchGenShorten
// @return error
func (ctx *Context) FetchGenShorten(shortKey string) (*ResFetchGenShorten, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	param := new(ParamFetchGenShorten)
	param.ShortKey = shortKey
	var res ResFetchGenShorten
	if err := wechat.Post(wx.ApiCgiBin + "/shorten/fetch").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.FetchGenShorten(shortKey)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to fetch short key",
			Desc:    "还原短key失败",
		}
	}
	return &res, nil
}
