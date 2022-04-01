package wx

import (
	"github.com/hhcool/gtls/log"
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

func (c *Context) newAccessTokenForMp() string {
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
		log.Warn("获取token失败", zap.Error(err))
		return ""
	}
	if res.Errcode != 0 {
		log.Warn("获取token失败", zap.Int("errcode", res.Errcode), zap.String("errmsg", res.Errmsg))
	}
	return res.AccessToken
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
func (c *Context) newAccessTokenForWork() string {
	var res ResponseAccessToken
	err := wechat.Get(ApiWork + "/gettoken").
		SetQuery(ParamWorkAccessToken{
			Corpid:     c.App.Appid,
			Corpsecret: c.App.AppSecret,
		}).
		BindJSON(&res).
		Do()
	if err != nil {
		log.Warn("获取token失败", zap.Error(err))
		return ""
	}
	if res.Errcode != 0 {
		log.Warn("获取token失败", zap.Int("errcode", res.Errcode), zap.String("errmsg", res.Errmsg))
	}
	return res.AccessToken
}

/**
-------------------------------------------------------
刷新全局的AccessToken
*/
func (w *Wechat) refreshAccessToken() {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("刷新token失败", zap.Any("err", r))
			w.refreshAccessToken()
		}
	}()
	expTicker := time.NewTicker(time.Minute * 55)
	for {
		select {
		case <-expTicker.C:
			apps := wechat.SMembers(RdsAppListPrefix).Val()
			log.Info("刷新token", zap.Strings("apps", apps))
			for i := range apps {
				appid := apps[i]
				if ctx, err := FindApp(appid); err == nil {
					ctx.NewAccessToken()
				}
			}
		}
	}
}
