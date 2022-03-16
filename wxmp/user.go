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
	Response
	Userinfo
}

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
		return nil, fmt.Errorf("查询用户信息失败，%s", err.Error())
	}
	if res.Errcode != 0 {
		return nil, fmt.Errorf("查询用户信息失败，%d-%s", res.Errcode, res.Errmsg)
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
