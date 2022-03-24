package wxwork

import (
	"fmt"
	"github.com/hhcool/wx"
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

func (ctx *Context) TagList(p ParamTagList) ([]ResponseTagListItem, error) {
	if !ctx.IsWork() {
		return nil, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	wechat := wx.NewWechat()
	var res ResponseTagList
	err := wechat.Post(wx.ApiWork + "/externalcontact/get_corp_tag_list").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, fmt.Errorf("企业微信：查询标签失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return nil, fmt.Errorf("企业微信：查询标签失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res.TagGroup, nil
}
