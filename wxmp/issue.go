package wxmp

import (
	"github.com/zohu/wx"
)

/**
发布
*/

type ResIssueList struct {
	wx.Response
	List  []ResIssueListItem1 `json:"list"`
	Page  int64               `json:"page"`
	Rows  int64               `json:"rows"`
	Total int64               `json:"total"`
	Count int64               `json:"count"`
}
type ResIssueListItem struct {
	wx.Response
	TotalCount int64               `json:"total_count"`
	ItemCount  int64               `json:"item_count"`
	Item       []ResIssueListItem1 `json:"item"`
}
type ResIssueListItem1 struct {
	ArticleId  string `json:"article_id"`
	UpdateTime int64  `json:"update_time"`
	Content    struct {
		NewsItem []struct {
			Title              string `json:"title"`
			Author             string `json:"author"`
			Digest             string `json:"digest"`
			Content            string `json:"content"`
			ContentSourceUrl   string `json:"content_source_url"`
			ThumbMediaId       string `json:"thumb_media_id"`
			ShowCoverPic       int    `json:"show_cover_pic"`
			NeedOpenComment    int    `json:"need_open_comment"`
			OnlyFansCanComment int    `json:"only_fans_can_comment"`
			Url                string `json:"url"`
			IsDeleted          bool   `json:"is_deleted"`
		} `json:"news_item"`
	} `json:"content"`
}
type ParamIssueList struct {
	Offset    int64 `json:"offset"`
	Count     int64 `json:"count"`
	NoContent int   `json:"no_content"`
}

// IssueList
// @Description: 查询发布文章列表
// @receiver ctx
// @param noContent
// @param page
// @param rows
// @return *ResIssueList
// @return error
func (ctx *Context) IssueList(noContent int, page int64, rows int64) (*ResIssueList, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
	}
	wechat := wx.NewWechat()
	var wxr ResIssueListItem
	if err := wechat.Post(wx.ApiCgiBin + "/freepublish/batchget").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamIssueList{
			Offset:    (page - 1) * rows,
			Count:     rows,
			NoContent: noContent,
		}).
		BindJSON(&wxr).
		Do(); err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if wxr.Errcode != 0 {
		if ctx.RetryAccessToken(wxr.Errcode) {
			return ctx.IssueList(noContent, page, rows)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: wxr.Errcode,
			Errmsg:  wxr.Errmsg,
			Err:     "failed to get publish list",
			Desc:    "查询发布记录失败",
		}
	}
	res := new(ResIssueList)
	res.Page = page
	res.Rows = rows
	res.Count = wxr.ItemCount
	res.Total = wxr.TotalCount
	res.List = wxr.Item
	res.Errcode = wxr.Errcode
	res.Errmsg = wxr.Errmsg
	return res, nil
}
