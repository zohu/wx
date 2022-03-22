package wxmp

import "mime/multipart"

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
	Response
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int    `json:"created_at"`
}

// MediaTemporaryAdd
// @Description: 新增临时素材
// @receiver ctx
// @param mediaType
// @param file
func (ctx *Context) MediaTemporaryAdd(mediaType MediaType, file *multipart.FileHeader) (*ResMediaTemporaryAdd, error) {
	return nil, nil
}

// MediaTemporaryQuery
// @Description: 查询临时素材
// @receiver ctx
// @param mediaId
// @return {}
func (ctx *Context) MediaTemporaryQuery(mediaId string) {}

type ParamMediaForeverAdd struct {
	Title              string `json:"title"`
	ThumbMediaId       string `json:"thumb_media_id"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	Author             string `json:"author,omitempty"`
	Digest             string `json:"digest,omitempty"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}
type ResMediaForeverAdd struct {
	Response
	MediaID string `json:"media_id"`
}

// MediaForeverAdd
// @Description: 新增永久素材
// @receiver ctx
// @param param
// @return *ResMediaForeverAdd
// @return error
func (ctx *Context) MediaForeverAdd(param *ParamMediaForeverAdd) (*ResMediaForeverAdd, error) {
	return nil, nil
}

func (ctx *Context) MediaForeverQuery(mediaID string) {

}

func (ctx *Context) MediaForeverDelete(mediaID string) {

}
func (ctx *Context) MediaForeverUpdate() {}

func (ctx *Context) MediaCount() {}
func (ctx *Context) MediaList()  {}

func (ctx *Context) MediaFileUpload() {

}
