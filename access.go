package wx

import (
	"github.com/hhcool/log"
	"strings"
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
		log.Errorf("newAccessTokenForMp %s", err.Error())
		return ""
	}
	if res.Errcode != 0 {
		log.Errorf("newAccessTokenForMp (%d)-%s", res.Errcode, res.Errmsg)
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
		log.Errorf("newAccessTokenForWork %s", err.Error())
		return ""
	}
	if res.Errcode != 0 {
		log.Errorf("newAccessTokenForWork (%d)-%s", res.Errcode, res.Errmsg)
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
			log.Errorf("refreshAccessToken %v", r)
			w.refreshAccessToken()
		}
	}()
	expTicker := time.NewTicker(time.Minute * 55)
	for {
		select {
		case <-expTicker.C:
			apps := wechat.SMembers(RdsAppListPrefix).Val()
			log.Infof("refreshAccessToken 刷新Token [%s]", strings.Join(apps, ","))
			for i := range apps {
				appid := apps[i]
				ctx := FindApp(appid)
				ctx.NewAccessToken()
			}
		}
	}
}
