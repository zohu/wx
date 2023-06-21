package wxmp

import (
	"github.com/zohu/wx"
)

type Kf struct {
	KfId             string `json:"kf_id"`
	KfAccount        string `json:"kf_account"`
	KfHeadimgurl     string `json:"kf_headimgurl,omitempty"`
	KfNick           string `json:"kf_nick,omitempty"`
	KfWx             string `json:"kf_wx,omitempty"`
	InviteWx         string `json:"invite_wx,omitempty"`
	InviteStatus     string `json:"invite_status,omitempty"`
	InviteExpireTime int    `json:"invite_expire_time,omitempty"`
	Status           int    `json:"status,omitempty"`
	AcceptedCase     int    `json:"accepted_case,omitempty"`
}

type ResKfList struct {
	wx.Response
	KfList []Kf `json:"kf_list"`
}

// KfList
// @Description: 获取客服基本信息
// @receiver ctx
// @param shortKey
// @return *ResKfList
// @return error
func (ctx *Context) KfList() (*ResKfList, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	var res ResKfList
	if err := wechat.Get(wx.ApiCgiBin + "/customservice/getkflist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
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
			return ctx.KfList()
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get kf list",
			Desc:    "获取客服基本信息失败",
		}
	}
	return &res, nil
}

// KfOnlineList
// @Description: 查询在线客服基本信息
// @receiver ctx
// @return *ResKfList
// @return error
func (ctx *Context) KfOnlineList() (*ResKfList, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	var res ResKfList
	if err := wechat.Get(wx.ApiCgiBin + "/customservice/getonlinekflist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
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
			return ctx.KfOnlineList()
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get online kf list",
			Desc:    "查询在线客服基本信息失败",
		}
	}
	return &res, nil
}

type ParamKfAdd struct {
	KfAccount string `json:"kf_account"`
	Nickname  string `json:"nickname"`
}

// KfAdd
// @Description: 添加客服帐号
// @receiver ctx
// @param account
// @param name
// @return *wx.Response
// @return error
func (ctx *Context) KfAdd(account string, name string) (*wx.Response, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	param := new(ParamKfAdd)
	param.KfAccount = account
	param.Nickname = name
	var res wx.Response
	if err := wechat.Post(wx.ApiCgiBin + "/customservice/kfaccount/add").
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
			return ctx.KfAdd(account, name)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to add kf account",
			Desc:    "添加客服帐号失败",
		}
	}
	return &res, nil
}

type ParamKfInvite struct {
	KfAccount string `json:"kf_account"`
	InviteWx  string `json:"invite_wx"`
}

// KfInvite
// @Description: 邀请绑定客服帐号
// @receiver ctx
// @param account
// @param inviteWx
// @return *wx.Response
// @return error
func (ctx *Context) KfInvite(account string, inviteWx string) (*wx.Response, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	param := new(ParamKfInvite)
	param.KfAccount = account
	param.InviteWx = inviteWx
	var res wx.Response
	if err := wechat.Post(wx.ApiCgiBin + "/customservice/kfaccount/inviteworker").
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
			return ctx.KfInvite(account, inviteWx)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to invite kf account",
			Desc:    "邀请绑定客服账号失败",
		}
	}
	return &res, nil
}

type ParamKfUpdate struct {
	KfAccount string `json:"kf_account"`
	Nickname  string `json:"nickname"`
}

// KfUpdate
// @Description: 设置客服信息
// @receiver ctx
// @param account
// @param name
// @return *wx.Response
// @return error
func (ctx *Context) KfUpdate(account string, name string) (*wx.Response, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	param := new(ParamKfUpdate)
	param.KfAccount = account
	param.Nickname = name
	var res wx.Response
	if err := wechat.Post(wx.ApiCgiBin + "/customservice/kfaccount/update").
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
			return ctx.KfUpdate(account, name)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to update kf account",
			Desc:    "设置客服信息失败",
		}
	}
	return &res, nil
}

//type ParamKfUploadHeadImg struct {
//	wx.ParamAccessToken
//	KfAccount string `json:"kf_account" query:"kf_account"`
//}
//
