package wxmp

import (
	"fmt"
	"github.com/zohu/wx"
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
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	p := new(ParamUserFromOpenid)
	p.AccessToken = ctx.GetAccessToken()
	p.Openid = openid
	var res ResUserFromOpenid
	if err := wechat.Post(wx.ApiMp + "/user/info").
		SetQuery(&p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询用户信息失败，%s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserFromOpenid(openid)
		}
		return nil, fmt.Errorf("%s 查询用户信息失败，%s，%d-%s", ctx.App.Appid, openid, res.Errcode, res.Errmsg)
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
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
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
		return nil, fmt.Errorf("%s 查询用户列表失败 %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserList(nextOpenID)
		}
		return nil, fmt.Errorf("%s 查询用户列表失败，%d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type Tag struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type ParamUserTagCreate struct {
	Tag *Tag `json:"tag"`
}
type ResUserTagCreate struct {
	wx.Response
	Tag *Tag `json:"tag"`
}

// UserTagCreate
// @Description: 创建标签
// @receiver ctx
// @param name
// @return *ResUserTagCreate
// @return error
func (ctx *Context) UserTagCreate(name string) (*ResUserTagCreate, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res ResUserTagCreate
	param := new(ParamUserTagCreate)
	param.Tag.Name = name
	if err := wechat.Post(wx.ApiMp + "/tags/create").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("创建标签失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagCreate(name)
		}
		return nil, fmt.Errorf("创建标签失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResUserTagQuery struct {
	wx.Response
	Tags []struct {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"tags"`
}

// UserTagQuery
// @Description: 获取公众号已创建的标签
// @receiver ctx
// @return *ResUserTagQuery
// @return error
func (ctx *Context) UserTagQuery() (*ResUserTagQuery, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res ResUserTagQuery
	if err := wechat.Get(wx.ApiMp + "/tags/get").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("查询标签失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagQuery()
		}
		return nil, fmt.Errorf("查询标签失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagEdit struct {
	Tag *Tag `json:"tag"`
}

// UserTagEdit
// @Description: 编辑标签
// @receiver ctx
// @param id
// @param name
// @return *wx.Response
// @return error
func (ctx *Context) UserTagEdit(id int64, name string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res wx.Response
	param := new(ParamUserTagEdit)
	param.Tag.ID = id
	param.Tag.Name = name
	if err := wechat.Post(wx.ApiMp + "/tags/update").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("更新标签失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagEdit(id, name)
		}
		return nil, fmt.Errorf("更新标签失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagDel struct {
	Tag *Tag `json:"tag"`
}

// UserTagDel
// @Description: 删除标签
// @receiver ctx
// @param id
// @return *wx.Response
// @return error
func (ctx *Context) UserTagDel(id int64) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res wx.Response
	param := new(ParamUserTagEdit)
	param.Tag.ID = id
	if err := wechat.Post(wx.ApiMp + "/tags/delete").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("删除标签失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagDel(id)
		}
		return nil, fmt.Errorf("删除标签失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagGetUser struct {
	Tagid      int64  `json:"tagid"`
	NextOpenid string `json:"next_openid,omitempty"`
}
type ResUserTagGetUser struct {
	wx.Response
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

// UserTagGetUser
// @Description: 获取标签下粉丝列表
// @receiver ctx
// @param id
// @param nextOpenid
// @return *ResUserTagGetUser
// @return error
func (ctx *Context) UserTagGetUser(id int64, nextOpenid string) (*ResUserTagGetUser, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res ResUserTagGetUser
	param := new(ParamUserTagGetUser)
	param.Tagid = id
	param.NextOpenid = nextOpenid
	if err := wechat.Post(wx.ApiMp + "/user/tag/get").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取标签下粉丝列表 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagGetUser(id, nextOpenid)
		}
		return nil, fmt.Errorf("获取标签下粉丝列表 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagBatch struct {
	OpenidList []string `json:"openid_list"`
	Tagid      int64    `json:"tagid"`
}

// UserTagBatch
// @Description: 批量为用户打标签
// @receiver ctx
// @param openid
// @param tagid
// @return *wx.Response
// @return error
func (ctx *Context) UserTagBatch(openid []string, tagid int64) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res wx.Response
	param := new(ParamUserTagBatch)
	param.OpenidList = openid
	param.Tagid = tagid
	if err := wechat.Post(wx.ApiMp + "/tags/members/batchtagging").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("批量为用户打标签 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagBatch(openid, tagid)
		}
		return nil, fmt.Errorf("批量为用户打标签 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagUnBatch struct {
	OpenidList []string `json:"openid_list"`
	Tagid      int64    `json:"tagid"`
}

// UserTagUnBatch
// @Description: 批量为用户取消标签
// @receiver ctx
// @param openid
// @param tagid
// @return *wx.Response
// @return error
func (ctx *Context) UserTagUnBatch(openid []string, tagid int64) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res wx.Response
	param := new(ParamUserTagUnBatch)
	param.OpenidList = openid
	param.Tagid = tagid
	if err := wechat.Post(wx.ApiMp + "/tags/members/batchuntagging").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("批量为用户取消标签 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagUnBatch(openid, tagid)
		}
		return nil, fmt.Errorf("批量为用户取消标签 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserTagGetFromUser struct {
	Openid string `json:"openid"`
}
type ResUserTagGetFromUser struct {
	wx.Response
	TagidList []int64 `json:"tagid_list"`
}

// UserTagGetFromUser
// @Description: 获取用户身上的标签列表
// @receiver ctx
// @param openid
// @return *ResUserTagGetFromUser
// @return error
func (ctx *Context) UserTagGetFromUser(openid string) (*ResUserTagGetFromUser, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res ResUserTagGetFromUser
	param := new(ParamUserTagGetFromUser)
	param.Openid = openid
	if err := wechat.Post(wx.ApiMp + "/tags/getidlist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取用户身上的标签列表 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserTagGetFromUser(openid)
		}
		return nil, fmt.Errorf("获取用户身上的标签列表 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserRemarkUpdate struct {
	Openid string `json:"openid"`
	Remark string `json:"remark"`
}

// UserRemarkUpdate
// @Description: 设置用户备注名
// @receiver ctx
// @param openid
// @param remark
// @return *wx.Response
// @return error
func (ctx *Context) UserRemarkUpdate(openid string, remark string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var res wx.Response
	param := new(ParamUserRemarkUpdate)
	param.Openid = openid
	param.Remark = remark
	if err := wechat.Post(wx.ApiMp + "/user/info/updateremark").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("设置用户备注名 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserRemarkUpdate(openid, remark)
		}
		return nil, fmt.Errorf("设置用户备注名 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserGetBlackList struct {
	BeginOpenid string `json:"begin_openid"`
}
type ResUserGetBlackList struct {
	wx.Response
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

// UserGetBlackList
// @Description: 获取公众号的黑名单列表
// @receiver ctx
// @param beginOpenid
// @return *ResUserGetBlackList
// @return error
func (ctx *Context) UserGetBlackList(beginOpenid string) (*ResUserGetBlackList, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	param := new(ParamUserGetBlackList)
	param.BeginOpenid = beginOpenid
	var res ResUserGetBlackList
	if err := wechat.Post(wx.ApiMp + "/tags/members/getblacklist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取公众号的黑名单列表 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserGetBlackList(beginOpenid)
		}
		return nil, fmt.Errorf("获取公众号的黑名单列表 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserBlackListPush struct {
	OpenidList []string `json:"openid_list"`
}

// UserBlackListPush
// @Description: 拉黑用户
// @receiver ctx
// @param openidList
// @return *wx.Response
// @return error
func (ctx *Context) UserBlackListPush(openidList []string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	param := new(ParamUserBlackListPush)
	param.OpenidList = openidList
	var res wx.Response
	if err := wechat.Post(wx.ApiMp + "/tags/members/batchblacklist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("拉黑用户 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserBlackListPush(openidList)
		}
		return nil, fmt.Errorf("拉黑用户 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamUserBlackListUnPush struct {
	OpenidList []string `json:"openid_list"`
}

// UserBlackListUnPush
// @Description: 取消拉黑
// @receiver ctx
// @param openidList
// @return *wx.Response
// @return error
func (ctx *Context) UserBlackListUnPush(openidList []string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	param := new(ParamUserBlackListUnPush)
	param.OpenidList = openidList
	var res wx.Response
	if err := wechat.Post(wx.ApiMp + "/tags/members/batchunblacklist").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(param).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("取消拉黑 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.UserBlackListUnPush(openidList)
		}
		return nil, fmt.Errorf("取消拉黑 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}
