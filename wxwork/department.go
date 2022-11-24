package wxwork

import (
	"fmt"
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

func (ctx *Context) DepartmentCreat(p Department) (int64, error) {
	if !ctx.IsWork() {
		return -1, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	var res ResponseDepartmentCreat
	err := wechat.Post(wx.ApiWork + "/department/create").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return -1, fmt.Errorf("企业微信：创建部门失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentCreat(p)
		}
		return -1, fmt.Errorf("企业微信：创建部门失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res.Id, nil
}

func (ctx *Context) DepartmentUpdate(p Department) error {
	if !ctx.IsWork() {
		return fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	var res wx.Response
	err := wechat.Post(wx.ApiWork + "/department/update").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return fmt.Errorf("企业微信：更新部门失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentUpdate(p)
		}
		return fmt.Errorf("企业微信：更新部门失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return nil
}

type ParamDepartmentDelete struct {
	wx.ParamAccessToken
	Id int64 `json:"id" query:"id"`
}

func (ctx *Context) DepartmentDelete(p ParamDepartmentDelete) error {
	if !ctx.IsWork() {
		return fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res wx.Response
	err := wechat.Get(wx.ApiWork + "/department/delete").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return fmt.Errorf("企业微信：删除部门失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentDelete(p)
		}
		return fmt.Errorf("企业微信：删除部门失败（%d-%s）", res.Errcode, res.Errmsg)
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

func (ctx *Context) DepartmentList(p ParamDepartmentList) ([]Department, error) {
	if !ctx.IsWork() {
		return nil, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res ResponseDepartmentList
	err := wechat.Get(wx.ApiWork + "/department/list").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, fmt.Errorf("企业微信：查询部门失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.DepartmentList(p)
		}
		return nil, fmt.Errorf("企业微信：查询部门失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res.Department, nil
}
