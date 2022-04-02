package wxmp

import (
	"fmt"
	"github.com/hhcool/wx"
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
func (ctx *Context) Qrcode(p *ParamNewQrcode) (*ResQrcode, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
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
	err := wechat.Post(wx.ApiMp + "/qrcode/create").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&param).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%s 获取二维码失败（%s）", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.Qrcode(p)
		}
		return nil, fmt.Errorf("%s 获取二维码失败（%d-%s）", ctx.Appid(), res.Errcode, res.Errmsg)
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
func (ctx *Context) GenShorten(data string, ex ...int) (*ResGenShorten, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	param := new(ParamGenShorten)
	param.LongData = data
	if len(ex) > 0 {
		param.ExpireSeconds = ex[0]
	}
	var res ResGenShorten
	if err := wechat.Post(wx.ApiMp + "/shorten/gen").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取短KEY %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.GenShorten(data, ex...)
		}
		return nil, fmt.Errorf("获取短KEY %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
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
func (ctx *Context) FetchGenShorten(shortKey string) (*ResFetchGenShorten, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	param := new(ParamFetchGenShorten)
	param.ShortKey = shortKey
	var res ResFetchGenShorten
	if err := wechat.Post(wx.ApiMp + "/shorten/fetch").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("还原短key %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.FetchGenShorten(shortKey)
		}
		return nil, fmt.Errorf("还原短key %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}
