package wxmp

/**
基础消息能力
*/

import (
	"fmt"
	"github.com/zohu/wx"
)

type ParamMsgSetIndustry struct {
	IndustryId1 string `json:"industry_id1"`
	IndustryId2 string `json:"industry_id2"`
}

// MsgSetIndustry
// @Description: 设置所属行业
// @receiver ctx
// @param p
// @return *wx.Response
// @return error
func (ctx *Context) MsgSetIndustry(p *ParamMsgSetIndustry) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/template/api_set_industry").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("设置行业失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgSetIndustry(p)
		}
		return nil, fmt.Errorf("设置行业失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResGetIndustry struct {
	wx.Response
	PrimaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"primary_industry"`
	SecondaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"secondary_industry"`
}

// MsgGetIndustry
// @Description: 获取设置的行业信息
// @receiver ctx
// @return *ResGetIndustry
// @return error
func (ctx *Context) MsgGetIndustry() (*ResGetIndustry, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res ResGetIndustry
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiMp + "/template/get_industry").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取行业失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgGetIndustry()
		}
		return nil, fmt.Errorf("获取行业失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamMsgGetTemplateId struct {
	TemplateIdShort string `json:"template_id_short"`
}
type ResMsgGetTemplateId struct {
	wx.Response
	TemplateId string `json:"template_id"`
}

// MsgGetTemplateId
// @Description: 获得模板ID
// @receiver ctx
// @param shortId
// @return *ResMsgGetTemplateId
// @return error
func (ctx *Context) MsgGetTemplateId(shortId string) (*ResMsgGetTemplateId, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res ResMsgGetTemplateId
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/template/api_add_template").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamMsgGetTemplateId{TemplateIdShort: shortId}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获得模板ID失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgGetTemplateId(shortId)
		}
		return nil, fmt.Errorf("获得模板ID失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResMsgGetTemplateList struct {
	wx.Response
	TemplateList []struct {
		TemplateId      string `json:"template_id"`
		Title           string `json:"title"`
		PrimaryIndustry string `json:"primary_industry"`
		DeputyIndustry  string `json:"deputy_industry"`
		Content         string `json:"content"`
		Example         string `json:"example"`
	} `json:"template_list"`
}

// MsgGetTemplateList
// @Description: 获得模板列表
// @receiver ctx
// @return *ResMsgGetTemplateList
// @return error
func (ctx *Context) MsgGetTemplateList() (*ResMsgGetTemplateList, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res ResMsgGetTemplateList
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiMp + "/template/get_all_private_template").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获得模板列表失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgGetTemplateList()
		}
		return nil, fmt.Errorf("获得模板列表失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamMsgDelTemplate struct {
	TemplateId string `json:"template_id"`
}

// MsgDelTemplate
// @Description: 删除模板
// @receiver ctx
// @param templateId
// @return *wx.Response
// @return error
func (ctx *Context) MsgDelTemplate(templateId string) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/template/del_private_template").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamMsgDelTemplate{TemplateId: templateId}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("删除模板失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgDelTemplate(templateId)
		}
		return nil, fmt.Errorf("删除模板失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamMsgSendTemplate struct {
	Touser      string       `json:"touser"`
	TemplateId  string       `json:"template_id"`
	Url         string       `json:"url,omitempty"`
	Miniprogram *Miniprogram `json:"miniprogram,omitempty"`
	Data        map[string]struct {
		Value string `json:"value"`
		Color string `json:"color,omitempty"`
	} `json:"data"`
}
type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath,omitempty"`
}

// MsgSendTemplate
// @Description: 发送模板消息
// @receiver ctx
// @param p
// @return *wx.Response
// @return error
func (ctx *Context) MsgSendTemplate(p *ParamMsgSendTemplate) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/message/template/send").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("发送模板消息失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgSendTemplate(p)
		}
		return nil, fmt.Errorf("发送模板消息失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ParamMsgSubscribe struct {
	Touser      string       `json:"touser"`
	TemplateId  string       `json:"template_id"`
	Url         string       `json:"url,omitempty"`
	Miniprogram *Miniprogram `json:"miniprogram,omitempty"`
	Scene       string       `json:"scene"`
	Title       string       `json:"title"`
	Data        struct {
		Content struct {
			Value string `json:"value"`
			Color string `json:"color"`
		} `json:"content"`
	} `json:"data"`
}

// MsgSubscribe
// @Description: 推送订阅模板消息给到授权微信用户
// @receiver ctx
// @param p
// @return *wx.Response
// @return error
func (ctx *Context) MsgSubscribe(p *ParamMsgSubscribe) (*wx.Response, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiMp + "/message/template/subscribe").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("订阅模板消息失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgSubscribe(p)
		}
		return nil, fmt.Errorf("订阅模板消息失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}

type ResMsgGetAutoReply struct {
	wx.Response
	IsAddFriendReplyOpen   int `json:"is_add_friend_reply_open"`
	IsAutoreplyOpen        int `json:"is_autoreply_open"`
	AddFriendAutoreplyInfo struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"add_friend_autoreply_info"`
	MessageDefaultAutoreplyInfo struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"message_default_autoreply_info"`
	KeywordAutoreplyInfo struct {
		List []struct {
			RuleName        string `json:"rule_name"`
			CreateTime      int    `json:"create_time"`
			ReplyMode       string `json:"reply_mode"`
			KeywordListInfo []struct {
				Type      string `json:"type"`
				MatchMode string `json:"match_mode"`
				Content   string `json:"content"`
			} `json:"keyword_list_info"`
			ReplyListInfo []struct {
				Type     string `json:"type"`
				NewsInfo struct {
					List []struct {
						Title      string `json:"title"`
						Author     string `json:"author"`
						Digest     string `json:"digest"`
						ShowCover  int    `json:"show_cover"`
						CoverUrl   string `json:"cover_url"`
						ContentUrl string `json:"content_url"`
						SourceUrl  string `json:"source_url"`
					} `json:"list"`
				} `json:"news_info,omitempty"`
				Content string `json:"content,omitempty"`
			} `json:"reply_list_info"`
		} `json:"list"`
	} `json:"keyword_autoreply_info"`
}

// MsgGetAutoReply
// @Description: 获取公众号的自动回复规则
// @receiver ctx
// @return *ResMsgGetAutoReply
// @return error
func (ctx *Context) MsgGetAutoReply() (*ResMsgGetAutoReply, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	var res ResMsgGetAutoReply
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiMp + "/get_current_autoreply_info").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("获取自动回复信息失败 %s %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MsgGetAutoReply()
		}
		return nil, fmt.Errorf("获取自动回复信息失败 %s %d-%s", ctx.App.Appid, res.Errcode, res.Errmsg)
	}
	return &res, nil
}
