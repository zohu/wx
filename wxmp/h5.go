package wxmp

import (
	"fmt"
	"github.com/zohu/wx"
	"github.com/zohu/wx/utils"
	"net/url"
	"time"
)

type H5ScopeType string

const (
	H5ScopeTypeBase           H5ScopeType = "snsapi_base"
	H5ScopeTypeSnsapiUserinfo H5ScopeType = "snsapi_userinfo"
)

// H5GetOauth2URL
// @Description: 获取授权链接
// @receiver ctx
// @param redirectUri
// @param scope
// @param state
// @return string
// @return error
func (ctx *Context) H5GetOauth2URL(redirectUri string, scope H5ScopeType, state string) (string, error) {
	uri, err := url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("appid", ctx.Appid())
	params.Add("redirect_uri", redirectUri)
	params.Add("response_type", "code")
	params.Add("scope", string(scope))
	params.Add("state", state)
	uri.RawQuery = params.Encode()
	return uri.String() + "#wechat_redirect", nil
}

type H5Token struct {
	wx.Response
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}
type ResH5GetUserinfo struct {
	wx.Response
	Openid     string   `json:"openid"`     // 用户的唯一标识
	Nickname   string   `json:"nickname"`   // 用户昵称
	Sex        int      `json:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province"`   // 用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	Headimgurl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	Unionid    string   `json:"unionid"`
}

// H5GetUserinfo
// @Description: code换取用户信息
// @receiver ctx
// @param code
// @param openid
// @return *ResH5GetUserinfo
// @return error
func (ctx *Context) H5GetUserinfo(code string, scope H5ScopeType) (*ResH5GetUserinfo, error) {
	var token H5Token
	wechat := wx.NewWechat()
	if err := wechat.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		ctx.Appid(),
		ctx.App.AppSecret,
		code,
	)).BindJSON(&token).
		Do(); err != nil {
		return nil, fmt.Errorf("获取token失败，%s", err.Error())
	}
	if token.Errcode != 0 {
		return nil, fmt.Errorf("获取token失败，%d-%s", token.Errcode, token.Errmsg)
	}
	var userinfo ResH5GetUserinfo
	if scope == H5ScopeTypeBase {
		userinfo.Openid = token.Openid
		return &userinfo, nil
	}
	if err := wechat.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		token.AccessToken,
		token.Openid,
	)).BindJSON(&userinfo).
		Do(); err != nil {
		return nil, fmt.Errorf("获取用户信息失败，%s", err.Error())
	}
	if token.Errcode != 0 {
		return nil, fmt.Errorf("获取用户信息失败，%d-%s", token.Errcode, token.Errmsg)
	}
	return &userinfo, nil
}

type H5JsSdkConfig struct {
	NonceStr  string `json:"nonce_str"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

func (ctx *Context) H5GetJsSdkConfig(uri string) (*H5JsSdkConfig, error) {
	jssdk := new(H5JsSdkConfig)
	jssdk.Timestamp = time.Now().Unix()
	jssdk.NonceStr = utils.RandomStr(16)
	tk := ctx.GetTicket()
	if tk == "" {
		return jssdk, fmt.Errorf("获取ticket失败")
	}
	query := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		tk,
		jssdk.NonceStr,
		jssdk.Timestamp,
		uri,
	)
	jssdk.Signature = utils.Signature(query)
	return jssdk, nil
}
