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

func (ctx *Context) MediaForeverAdd()    {}
func (ctx *Context) MediaForeverQuery()  {}
func (ctx *Context) MediaForeverDelete() {}
func (ctx *Context) MediaForeverUpdate() {}

func (ctx *Context) MediaCount() {}
func (ctx *Context) MediaList()  {}
