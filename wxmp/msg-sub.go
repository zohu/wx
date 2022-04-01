package wxmp

import (
	"fmt"
	"github.com/hhcool/wx"
)

/**
订阅通知
*/

type ParamSubAddTemplate struct {
	Tid       string `json:"tid"`
	KidList   []int  `json:"kidList"`
	SceneDesc string `json:"sceneDesc"`
}
type ResSubAddTemplate struct {
	wx.Response
	PriTmplId string `json:"priTmplId"`
}

// SubAddTemplate
// @Description: 从公共模板库中选用模板，到私有模板库中
// @receiver ctx
// @param p
// @return *ResSubAddTemplate
// @return error
func (ctx *Context) SubAddTemplate(p *ParamSubAddTemplate) (*ResSubAddTemplate, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResSubAddTemplate
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiWxa + "/newtmpl/addtemplate").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("选用模板失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubAddTemplate(p)
		}
		return nil, fmt.Errorf("选用模板失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamSubDelTemplate struct {
	PriTmplId string `json:"priTmplId"`
}

// SubDelTemplate
// @Description: 删除模板
// @receiver ctx
// @param priTmplId
// @return *wx.Response
// @return error
func (ctx *Context) SubDelTemplate(priTmplId string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiWxa + "/newtmpl/deltemplate").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamSubDelTemplate{PriTmplId: priTmplId}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("删除模板失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubDelTemplate(priTmplId)
		}
		return nil, fmt.Errorf("删除模板失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResSubGetCategory struct {
	wx.Response
	Data []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

// SubGetCategory
// @Description: 获取公众号类目
// @receiver ctx
// @return *ResSubGetCategory
// @return error
func (ctx *Context) SubGetCategory() (*ResSubGetCategory, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResSubGetCategory
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiWxa + "/newtmpl/getcategory").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取公众号类目失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubGetCategory()
		}
		return nil, fmt.Errorf("获取公众号类目失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamSubGetTemplateKeywords struct {
	wx.ParamAccessToken
	Tid string `json:"tid" query:"tid"`
}
type ResSubGetTemplateKeywords struct {
	wx.Response
	Data []struct {
		Kid     int    `json:"kid"`
		Name    string `json:"name"`
		Example string `json:"example"`
		Rule    string `json:"rule"`
	} `json:"data"`
}

// SubGetTemplateKeywords
// @Description: 获取模板中的关键词
// @receiver ctx
// @param tid
// @return *ResSubGetTemplateKeywords
// @return error
func (ctx *Context) SubGetTemplateKeywords(tid string) (*ResSubGetTemplateKeywords, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	query := new(ParamSubGetTemplateKeywords)
	query.AccessToken = ctx.GetAccessToken()
	query.Tid = tid
	var res ResSubGetTemplateKeywords
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiWxa + "/newtmpl/getpubtemplatekeywords").
		SetQuery(query).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取模板中的关键词失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubGetTemplateKeywords(tid)
		}
		return nil, fmt.Errorf("获取模板中的关键词失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamSubGetTemplateTitle struct {
	wx.ParamAccessToken
	Ids   string `json:"ids" query:"ids"`
	Start int    `json:"start" query:"start"`
	Limit int    `json:"limit" query:"limit"`
}
type ResSubGetTemplateTitle struct {
	wx.Response
	Data []struct {
		Tid        int    `json:"tid"`
		Title      string `json:"title"`
		Type       int    `json:"type"`
		CategoryId string `json:"categoryId"`
	} `json:"data"`
}

// SubGetTemplateTitle
// @Description: 获取类目下的公共模板
// @receiver ctx
// @param ids
// @param start
// @param limit
// @return *ResSubGetTemplateTitle
// @return error
func (ctx *Context) SubGetTemplateTitle(ids string, start int, limit int) (*ResSubGetTemplateTitle, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	query := new(ParamSubGetTemplateTitle)
	query.AccessToken = ctx.GetAccessToken()
	query.Ids = ids
	query.Start = start
	query.Limit = limit
	var res ResSubGetTemplateTitle
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiWxa + "/newtmpl/getpubtemplatetitles").
		SetQuery(query).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取类目下的公共模板失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubGetTemplateTitle(ids, start, limit)
		}
		return nil, fmt.Errorf("获取类目下的公共模板失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResSubGetTemplates struct {
	wx.Response
	Data []struct {
		PriTmplId string `json:"priTmplId"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		Example   string `json:"example"`
		Type      int    `json:"type"`
	} `json:"data"`
}

// SubGetTemplates
// @Description: 获取私有模板列表
// @receiver ctx
// @return *ResSubGetTemplates
// @return error
func (ctx *Context) SubGetTemplates() (*ResSubGetTemplates, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResSubGetTemplates
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiWxa + "/newtmpl/gettemplate").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取私有模板列表失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubGetTemplates()
		}
		return nil, fmt.Errorf("获取私有模板列表失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamSubBizSend struct {
	Touser      string       `json:"touser"`
	TemplateId  string       `json:"template_id"`
	Page        string       `json:"page,omitempty"`
	Miniprogram *Miniprogram `json:"miniprogram,omitempty"`
	Data        map[string]struct {
		Value string `json:"value"`
	} `json:"data"`
}

// SubBizSend
// @Description: 发送订阅通知
// @receiver ctx
// @param p
// @return *wx.Response
// @return error
func (ctx *Context) SubBizSend(p *ParamSubBizSend) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/message/subscribe/bizsend").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("发送订阅通知失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubBizSend(p)
		}
		return nil, fmt.Errorf("发送订阅通知失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}
