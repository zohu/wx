package wxmp

import (
	"fmt"
	"github.com/hhcool/wx"
)

type Userinfo struct {
	Subscribe      int    `json:"subscribe"`
	Openid         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int    `json:"subscribe_time"`
	Unionid        string `json:"unionid"`
	Remark         string `json:"remark"`
	Groupid        int    `json:"groupid"`
	TagidList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}
type ParamUserFromOpenid struct {
	wx.ParamAccessToken
	Openid string `query:"openid"`
}
type ResUserFromOpenid struct {
	wx.Response
	Userinfo
}

// UserFromOpenid
// @Description: 通过openID获取用户信息
// @receiver ctx
// @param openid
// @return *Userinfo
// @return error
func (ctx *Context) UserFromOpenid(openid string) (*Userinfo, error) {
	wechat := wx.NewWechat()
	p := new(ParamUserFromOpenid)
	p.AccessToken = ctx.GetAccessToken()
	p.Openid = openid
	var res ResUserFromOpenid
	if err := wechat.Post(wx.ApiMp + "/user/info").
		SetQuery(&p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询用户信息失败，%s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserFromOpenid(openid)
		}
		return nil, fmt.Errorf("%s 查询用户信息失败，%s，%d-%s", ctx.Appid(), openid, res.Errcode, res.Errmsg)
	}
	return &Userinfo{
		Subscribe:      res.Subscribe,
		Openid:         res.Openid,
		Language:       res.Language,
		SubscribeTime:  res.SubscribeTime,
		Unionid:        res.Unionid,
		Remark:         res.Remark,
		Groupid:        res.Groupid,
		TagidList:      res.TagidList,
		SubscribeScene: res.SubscribeScene,
		QrScene:        res.QrScene,
		QrSceneStr:     res.QrSceneStr,
	}, nil
}

type ParamQueryUserList struct {
	wx.ParamAccessToken
	NextOpenid string `query:"next_openid,omitempty"`
}
type ResQueryUserList struct {
	wx.Response
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

// UserList
// @Description: 查询用户列表
// @receiver ctx
// @param nextOpenID
// @return *ResQueryUserList
// @return error
func (ctx *Context) UserList(nextOpenID string) (*ResQueryUserList, error) {
	var res ResQueryUserList
	param := new(ParamQueryUserList)
	param.AccessToken = ctx.GetAccessToken()
	if nextOpenID != "" {
		param.NextOpenid = nextOpenID
	}
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiMp + "/user/get").
		SetQuery(&param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询用户列表失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserList(nextOpenID)
		}
		return nil, fmt.Errorf("%s 查询用户列表失败，%d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}
