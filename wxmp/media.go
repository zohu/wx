package wxmp

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/hhcool/wx"
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
	if err := wechat.Post(wx.ApiMp + "/media/upload").
		SetQuery(&q).
		SetHeader(gout.H{"Content-Type": bw.FormDataContentType()}).
		SetBody(pr).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 上传临时素材失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MediaTemporaryAdd(mediaType, file, fileName)
		}
		return nil, fmt.Errorf("%s 上传临时素材失败 %s", ctx.Appid(), res.Errmsg)
	}
	return res, nil
}
