package wxwork

import (
	"github.com/zohu/wx"
)

type Tag struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CreateTime int    `json:"create_time"`
	Order      int    `json:"order"`
}

type ParamTagList struct {
	TagId   string `json:"tag_id,omitempty"`
	GroupId string `json:"group_id,omitempty"`
}
type ResponseTagList struct {
	wx.Response
	TagGroup []ResponseTagListItem `json:"tag_group"`
}
type ResponseTagListItem struct {
	GroupId    string `json:"group_id"`
	GroupName  string `json:"group_name"`
	CreateTime int    `json:"create_time"`
	Tag        []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		CreateTime int    `json:"create_time"`
		Order      int64  `json:"order"`
	} `json:"tag"`
	Order int `json:"order"`
}

func (ctx *Context) TagList(p ParamTagList) ([]ResponseTagListItem, *wx.Err) {
	if !ctx.IsWork() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	var res ResponseTagList
	err := wechat.Post(wx.ApiWorkCgiBin + "/externalcontact/get_corp_tag_list").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
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
			return ctx.TagList(p)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get tag list",
			Desc:    "查询标签失败",
		}
	}
	return res.TagGroup, nil
}
