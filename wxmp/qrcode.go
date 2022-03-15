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
	ExpireSeconds int64  `json:"expire_seconds"` // 二维码有效时间，以秒为单位。 最大不超过2592000（即30天）
	ActionName    string `json:"action_name"`    // 二维码类型，QR_SCENE为临时的整型参数值，QR_STR_SCENE为临时的字符串参数值，QR_LIMIT_SCENE为永久的整型参数值，QR_LIMIT_STR_SCENE为永久的字符串参数值
	ActionInfo    struct {
		Scene struct {
			SceneId  int    `json:"scene_id"`  // 场景值ID，临时二维码时为32位非0整型，永久二维码时最大值为100000
			SceneStr string `json:"scene_str"` // 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
		} `json:"scene"`
	} `json:"action_info"` // 二维码详情
}

type ResQrcode struct {
	Response
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	Url           string `json:"url"`
}

func (ctx *Context) Qrcode(p ParamNewQrcode) (string, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return "", fmt.Errorf("企业微信：应用 %s 非企业号", ctx.App.Appid)
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
		return "", fmt.Errorf("微信：获取二维码失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return "", fmt.Errorf("微信：获取二维码失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return fmt.Sprintf("%s/showqrcode?ticket=%s", wx.ApiMp, res.Ticket), nil
}
