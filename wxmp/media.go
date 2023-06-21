package wxmp

import (
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
func (ctx *Context) MediaTemporaryAdd(mediaType MediaType, file io.Reader, fileName string) (*ResMediaTemporaryAdd, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
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
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MediaTemporaryAdd(mediaType, file, fileName)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to upload temporary media",
			Desc:    "上传临时素材失败",
		}
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
func (ctx *Context) MediaList(mediaType MediaType, page int64, rows int64) (*ResMediaList, *wx.Err) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not public app",
			Desc:  "非公众号应用",
		}
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
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if wxr.Errcode != 0 {
		if ctx.RetryAccessToken(wxr.Errcode) {
			return ctx.MediaList(mediaType, page, rows)
		}
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: wxr.Errcode,
			Errmsg:  wxr.Errmsg,
			Err:     "failed to get material",
			Desc:    "查询永久素材失败",
		}
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
