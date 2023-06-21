package wxwork

import (
	"github.com/zohu/wx"
)

type Department struct {
	Id       int64  `json:"id" desc:"部门ID"`
	Name     string `json:"name" desc:"部门名称"`
	NameEn   string `json:"name_en" desc:"部门名称英文"`
	Parentid int64  `json:"parentid" desc:"父部门ID"`
	Order    int64  `json:"order" desc:"在父部门中的排序，越大越靠前"`
}

type ResponseDepartmentCreat struct {
	wx.Response
	Id int64 `json:"id" desc:"创建的部门ID"`
}

func (ctx *Context) DepartmentCreat(p Department) (int64, *wx.Err) {
	if !ctx.IsWork() {
		return -1, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	var res ResponseDepartmentCreat
	err := wechat.Post(wx.ApiWorkCgiBin + "/department/create").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return -1, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentCreat(p)
		}
		return -1, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to create department",
			Desc:    "创建部门失败",
		}
	}
	return res.Id, nil
}

func (ctx *Context) DepartmentUpdate(p Department) *wx.Err {
	if !ctx.IsWork() {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	var res wx.Response
	err := wechat.Post(wx.ApiWorkCgiBin + "/department/update").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentUpdate(p)
		}
		return &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to update department",
			Desc:    "更新部门失败",
		}
	}
	return nil
}

type ParamDepartmentDelete struct {
	wx.ParamAccessToken
	Id int64 `json:"id" query:"id"`
}

func (ctx *Context) DepartmentDelete(p ParamDepartmentDelete) *wx.Err {
	if !ctx.IsWork() {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res wx.Response
	err := wechat.Get(wx.ApiWorkCgiBin + "/department/delete").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentDelete(p)
		}
		return &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to delete department",
			Desc:    "删除部门失败",
		}
	}
	return nil
}

type ParamDepartmentList struct {
	wx.ParamAccessToken
	Id int64 `json:"id" query:"id"`
}
type ResponseDepartmentList struct {
	wx.Response
	Department []Department `json:"department"`
}

func (ctx *Context) DepartmentList(p ParamDepartmentList) ([]Department, *wx.Err) {
	if !ctx.IsWork() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res ResponseDepartmentList
	err := wechat.Get(wx.ApiWorkCgiBin + "/department/list").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentList(p)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to query department list",
			Desc:    "查询部门失败",
		}
	}
	return res.Department, nil
}
