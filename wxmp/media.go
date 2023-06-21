package wxmp

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/zohu/wx"
	"io"
	"mime/multipart"
)

/**
素材管理
*/

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVoice MediaType = "voice"
	MediaTypeVideo MediaType = "video"
	MediaTypeThumb MediaType = "thumb"
)

type ResMediaTemporaryAdd struct {
	wx.Response
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int    `json:"created_at"`
}

// MediaTemporaryAdd
// @Description: 新增临时素材
// @receiver ctx
// @param mediaType
// @param file
func (ctx *Context) MediaTemporaryAdd(mediaType MediaType, file io.Reader, fileName string) (*ResMediaTemporaryAdd, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	pr, pw := io.Pipe()
	bw := multipart.NewWriter(pw)
	go func() {
		fw, _ := bw.CreateFormFile("media", fileName)
		_, _ = io.Copy(fw, file)
		_ = bw.Close()
		_ = pw.Close()
	}()
	var q struct {
		wx.ParamAccessToken
		Type MediaType `query:"type"`
	}
	res := new(ResMediaTemporaryAdd)
	q.AccessToken = ctx.GetAccessToken()
	q.Type = mediaType
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiCgiBin + "/media/upload").
		SetQuery(&q).
		SetHeader(gout.H{"Content-Type": bw.FormDataContentType()}).
		SetBody(pr).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 上传临时素材失败 %s", ctx.App.Appid, err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MediaTemporaryAdd(mediaType, file, fileName)
		}
		return nil, fmt.Errorf("%s 上传临时素材失败 %s", ctx.App.Appid, res.Errmsg)
	}
	return res, nil
}

type ResMediaList struct {
	wx.Response
	List  []ResMediaListItem1 `json:"list"`
	Page  int64               `json:"page"`
	Rows  int64               `json:"rows"`
	Total int64               `json:"total"`
	Count int64               `json:"count"`
}
type ResMediaListItem struct {
	wx.Response
	TotalCount int                 `json:"total_count"`
	ItemCount  int                 `json:"item_count"`
	Item       []ResMediaListItem1 `json:"item"`
}
type ResMediaListItem1 struct {
	MediaId    string               `json:"media_id"`
	Name       string               `json:"name,omitempty"`
	Url        string               `json:"url,omitempty"`
	UpdateTime int                  `json:"update_time"`
	Content    *ResMediaListContent `json:"content,omitempty"`
}
type ResMediaListContent struct {
	NewsItem []struct {
		Title            string `json:"title"`
		ThumbMediaId     string `json:"thumb_media_id"`
		ShowCoverPic     int    `json:"show_cover_pic"`
		Author           string `json:"author"`
		Digest           string `json:"digest"`
		Content          string `json:"content"`
		Url              string `json:"url"`
		ContentSourceUrl string `json:"content_source_url"`
	} `json:"news_item"`
}
type ParamMediaList struct {
	Type   MediaType `json:"type"`
	Offset int64     `json:"offset"`
	Count  int64     `json:"count"`
}

// MediaList
// @Description: 查询永久素材
// @receiver ctx
// @param mediaType
// @param page
// @param rows
// @return *ResMediaList
// @return error
func (ctx *Context) MediaList(mediaType MediaType, page int64, rows int64) (*ResMediaList, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.App.Appid)
	}
	wechat := wx.NewWechat()
	var wxr ResMediaListItem
	if err := wechat.Post(wx.ApiCgiBin + "/material/batchget_material").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamMediaList{
			Offset: (page - 1) * rows,
			Count:  rows,
			Type:   mediaType,
		}).
		BindJSON(&wxr).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询永久素材失败 %s", ctx.App.Appid, err.Error())
	}
	if wxr.Errcode != 0 {
		if ctx.RetryAccessToken(wxr.Errcode) {
			return ctx.MediaList(mediaType, page, rows)
		}
		return nil, fmt.Errorf("%s 查询永久素材失败 %s", ctx.App.Appid, wxr.Errmsg)
	}
	res := new(ResMediaList)
	res.Page = page
	res.Rows = rows
	res.Count = int64(wxr.ItemCount)
	res.Total = int64(wxr.TotalCount)
	res.List = wxr.Item
	res.Errcode = wxr.Errcode
	res.Errmsg = wxr.Errmsg
	return res, nil
}
