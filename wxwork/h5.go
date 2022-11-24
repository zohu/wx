package wxwork

import (
	"fmt"
	"github.com/zohu/wx"
	"github.com/zohu/wx/utils"
	"github.com/zohu/wx/wxmp"
	"time"
)

type WxWorkUserInfo struct {
	wx.Response
	Userid     string `json:"userid"`
	UserTicket string `json:"user_ticket"`
}

// H5GetUserinfo code换取用户信息
func (ctx *Context) H5GetUserinfo(code string) (*WxWorkUserInfo, error) {
	wechat := wx.NewWechat()
	userinfo := WxWorkUserInfo{}
	if err := wechat.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=%s&code=%s",
		ctx.GetAccessToken(), code)).BindJSON(&userinfo).
		Do(); err != nil {
		return nil, fmt.Errorf("获取用户信息失败，%s", err.Error())
	}
	if userinfo.Errcode != 0 {
		return nil, fmt.Errorf("获取用户信息失败，%d-%s", userinfo.Errcode, userinfo.Errmsg)
	}
	return &userinfo, nil
}

func (ctx *Context) H5GetJsSdkConfig(uri string) (*wxmp.H5JsSdkConfig, error) {
	workJsSdk := new(wxmp.H5JsSdkConfig)
	workJsSdk.Timestamp = time.Now().Unix()
	workJsSdk.NonceStr = utils.RandomStr(16)
	tk := ctx.GetTicket()
	if tk == "" {
		return workJsSdk, fmt.Errorf("获取ticket失败")
	}
	query := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		tk,
		workJsSdk.NonceStr,
		workJsSdk.Timestamp,
		uri,
	)
	workJsSdk.Signature = utils.Signature(query)
	return workJsSdk, nil
}
