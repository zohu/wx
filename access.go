package wx

import (
	"github.com/zohu/zlog"
	"go.uber.org/zap"
	"time"
)

type ParamAccessToken struct {
	AccessToken string `json:"access_token" query:"access_token"`
}
type ResponseAccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

/**
-------------------------------------------------------
微信AccessToken
*/

type ParamMpAccessToken struct {
	GrantType string `json:"grant_type" query:"grant_type"`
	Appid     string `json:"appid" query:"appid"`
	Secret    string `json:"secret" query:"secret"`
}
type ParamMpTicket struct {
	AccessToken string `json:"access_token" query:"access_token"`
	Type        string `json:"type" query:"type"`
}
type ResponseMpTicket struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func (c *Context) newAccessTokenForMp() (string, string) {
	var res ResponseAccessToken
	err := wechat.Get(ApiMp + "/token").
		SetQuery(ParamMpAccessToken{
			GrantType: "client_credential",
			Appid:     c.App.Appid,
			Secret:    c.App.AppSecret,
		}).
		BindJSON(&res).
		Do()
	if err != nil {
		zlog.Warn("获取token失败", zap.Error(err))
		return "", ""
	}
	if res.Errcode != 0 {
		zlog.Warn("获取token失败", zap.Int("errcode", res.Errcode), zap.String("errmsg", res.Errmsg))
		return "", ""
	}
	var tk ResponseMpTicket
	err = wechat.Get(ApiMp + "/ticket/getticket").
		SetQuery(&ParamMpTicket{
			AccessToken: res.AccessToken,
			Type:        "jsapi",
		}).
		BindJSON(&tk).
		Do()
	return res.AccessToken, tk.Ticket
}

/**
-------------------------------------------------------
企业微信AccessToken
*/

type ParamWorkAccessToken struct {
	Corpid     string `json:"corpid"`
	Corpsecret string `json:"corpsecret"`
}

// newAccessTokenForWork
// @Description: 企业微信获取新token
// @receiver c
// @return string
func (c *Context) newAccessTokenForWork() (string, string) {
	var res ResponseAccessToken
	err := wechat.Get(ApiWork + "/gettoken").
		SetQuery(ParamWorkAccessToken{
			Corpid:     c.App.Appid,
			Corpsecret: c.App.AppSecret,
		}).
		BindJSON(&res).
		Do()
	if err != nil {
		zlog.Warn("获取token失败", zap.Error(err))
		return "", ""
	}
	if res.Errcode != 0 {
		zlog.Warn("获取token失败", zap.Int("errcode", res.Errcode), zap.String("errmsg", res.Errmsg))
	}
	var tk ResponseMpTicket
	err = wechat.Get(ApiWork + "/get_jsapi_ticket").
		SetQuery(ParamAccessToken{AccessToken: res.AccessToken}).
		BindJSON(&tk).
		Do()
	return res.AccessToken, tk.Ticket
}

/*
*
-------------------------------------------------------
刷新全局的AccessToken
*/
func (w *Wechat) refreshAccessToken() {
	defer func() {
		if r := recover(); r != nil {
			zlog.Warn("刷新token失败", zap.Any("err", r))
			w.refreshAccessToken()
		}
	}()
	expTicker := time.NewTicker(time.Minute * 55)
	for {
		select {
		case <-expTicker.C:
			apps := wechat.SMembers(RdsAppListPrefix).Val()
			zlog.Info("刷新token", zap.Strings("apps", apps))
			for i := range apps {
				appid := apps[i]
				if ctx, err := FindApp(appid); err == nil {
					ctx.NewAccessToken()
				}
			}
		}
	}
}
