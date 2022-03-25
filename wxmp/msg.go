package wxmp

import (
	"encoding/xml"
	"github.com/hhcool/wx/wxcpt"
)

// MessageType
// @Description: 消息类型
type MessageType string      // 基础消息类型
type MessageEventType string // 事件类型
type MessageOpenType string  // 开放平台消息类型

const (
	MessageTypeText            MessageType = "text"                      // 文本消息
	MessageTypeImage           MessageType = "image"                     // 图片消息
	MessageTypeVoice           MessageType = "voice"                     // 语音消息
	MessageTypeVideo           MessageType = "video"                     // 视频消息
	MessageTypeMusic           MessageType = "music"                     // 音乐消息
	MessageTypeNews            MessageType = "news"                      // 图文消息
	MessageTypeShortvideo      MessageType = "shortvideo"                // 小视频消息
	MessageTypeLocation        MessageType = "location"                  // 地理位置消息
	MessageTypeLink            MessageType = "link"                      // 链接消息
	MessageTypeMiniprogrampage MessageType = "miniprogrampage"           // 小程序卡片消息
	MessageTypeTransfer        MessageType = "transfer_customer_service" // 消息消息转发到客服
	MessageTypeEvent           MessageType = "event"                     // 事件消息
)
const (
	MessageEventSubscribe         MessageEventType = "subscribe"                 // 关注事件
	MessageEventUnsubscribe       MessageEventType = "unsubscribe"               // 取关事件
	MessageEventScan              MessageEventType = "SCAN"                      // 扫描二维码事件
	MessageLocation               MessageEventType = "LOCATION"                  // 上报地理位置事件
	MessageClick                  MessageEventType = "CLICK"                     // 点击菜单拉取消息时的事件
	MessageView                   MessageEventType = "VIEW"                      // 点击菜单跳转链接时的事件推送
	MessageScancodePush           MessageEventType = "scancode_push"             // 扫码推事件的事件推送
	MessageScancodeWaitmsg        MessageEventType = "scancode_waitmsg"          // 扫码推事件且弹出“消息接收中”提示框的事件推送
	MessagePicSysphoto            MessageEventType = "pic_sysphoto"              // 弹出系统拍照发图的事件推送
	MessagePicPhotoOrAlbum        MessageEventType = "pic_photo_or_album"        // 弹出拍照或者相册发图的事件推送
	MessagePicWeixin              MessageEventType = "pic_weixin"                // 弹出微信相册发图器的事件推送
	MessageLocationSelect         MessageEventType = "location_select"           // 弹出地理位置选择器的事件推送
	MessageTemplateSendJobFinish  MessageEventType = "TEMPLATESENDJOBFINISH"     // 发送模板消息推送通知
	MessageMassSendJobFinish      MessageEventType = "MASSSENDJOBFINISH"         // 群发消息推送通知
	MessageWxaMediaCheck          MessageEventType = "wxa_media_check"           // 异步校验图片/音频是否含有违法违规内容推送事件
	MessageSubscribeMsgPopupEvent MessageEventType = "subscribe_msg_popup_event" // 订阅通知事件推送
	MessagePublishJobFinish       MessageEventType = "PUBLISHJOBFINISH"          // 发布任务完成
)
const (
	OpenTypeVerifyTicket              MessageOpenType = "component_verify_ticket"    // 返回ticket
	OpenTypeAuthorized                MessageOpenType = "authorized"                 // 授权
	OpenTypeUnauthorized              MessageOpenType = "unauthorized"               // 取消授权
	OpenTypeUpdateAuthorized          MessageOpenType = "updateauthorized"           // 更新授权
	OpenTypeNotifyThirdFasterRegister MessageOpenType = "notify_third_fasteregister" // 注册审核事件推送
)

// CommonMessage
// @Description: 消息中通用的结构
type CommonMessage struct {
	XMLName      xml.Name    `xml:"xml"`
	ToUserName   wxcpt.CDATA `json:"ToUserName" xml:"ToUserName"`
	FromUserName wxcpt.CDATA `json:"FromUserName" xml:"FromUserName"`
	CreateTime   int64       `json:"CreateTime" xml:"CreateTime"`
	MsgType      MessageType `json:"MsgType" xml:"MsgType"`
}

// Message
// @Description: 微信推送的消息
type Message struct {
	CommonMessage
	// 普通消息
	MsgId         int64  `json:"MsgId,omitempty" xml:"MsgId,omitempty"`               // 普通消息的ID
	TemplateMsgID int64  `json:"MsgID,omitempty" xml:"MsgID,omitempty"`               // 模板消息的ID                                             // 模板消息推送成功的消息是MsgID
	Content       string `json:"Content,omitempty" xml:"Content,omitempty"`           // 文本消息
	PicUrl        string `json:"PicUrl,omitempty" xml:"PicUrl,omitempty"`             // 图片消息
	MediaId       string `json:"MediaId,omitempty" xml:"MediaId,omitempty"`           // 图片消息、语音消息、视频消息
	Format        string `json:"Format,omitempty" xml:"Format,omitempty"`             // 语音消息，语音格式，如amr，speex等
	Recognition   string `json:"Recognition,omitempty" xml:"Recognition,omitempty"`   // 语音消息，识别结果
	ThumbMediaId  string `json:"ThumbMediaId,omitempty" xml:"ThumbMediaId,omitempty"` // 视频消息，缩略图
	LocationX     string `json:"Location_X,omitempty" xml:"Location_X,omitempty"`     // 位置消息，纬度
	LocationY     string `json:"Location_Y,omitempty" xml:"Location_Y,omitempty"`     // 位置消息，经度
	Scale         int64  `json:"Scale,omitempty" xml:"Scale,omitempty"`               // 位置消息，地图缩放大小
	Label         string `json:"Label,omitempty" xml:"Label,omitempty"`               // 位置消息，地理位置信息
	Title         string `json:"Title,omitempty" xml:"Title,omitempty"`               // 链接消息，标题
	Description   string `json:"Description,omitempty" xml:"Description,omitempty"`   // 链接消息，描述
	Url           string `json:"Url,omitempty" xml:"Url,omitempty"`                   // 链接消息

	// 事件消息
	Event     string `json:"Event,omitempty" xml:"Event,omitempty"`         // 事件消息
	EventKey  string `json:"EventKey,omitempty" xml:"EventKey,omitempty"`   // 事件，二维码消息、关注、菜单
	Ticket    string `json:"Ticket,omitempty" xml:"Ticket,omitempty"`       // 事件，二维码消息，二维码ticket
	Latitude  string `json:"Latitude,omitempty" xml:"Latitude,omitempty"`   // 事件，地理位置，纬度
	Longitude string `json:"Longitude,omitempty" xml:"Longitude,omitempty"` // 事件，地理位置，经度
	Precision int64  `json:"Precision,omitempty" xml:"Precision,omitempty"` // 事件，地理位置，精度
}

// MessageText
// @Description: 回复文本消息
type MessageText struct {
	CommonMessage
	Content wxcpt.CDATA `xml:"Content"`
}

// MessageImage
// @Description: 回复图片消息
type MessageImage struct {
	CommonMessage
	Image struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Image"`
}

// MessageVoice
// @Description: 回复语音消息
type MessageVoice struct {
	CommonMessage
	Voice struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Voice"`
}

// MessageVideo
// @Description: 回复视频消息
type MessageVideo struct {
	CommonMessage
	Video struct {
		MediaID     string `xml:"MediaId"`
		Title       string `xml:"Title,omitempty"`
		Description string `xml:"Description,omitempty"`
	} `xml:"Video"`
}

// MessageMusic
// @Description: 回复音乐消息
type MessageMusic struct {
	CommonMessage
	Music struct {
		Title        string `xml:"Title"        `
		Description  string `xml:"Description"  `
		MusicURL     string `xml:"MusicUrl"     `
		HQMusicURL   string `xml:"HQMusicUrl"   `
		ThumbMediaID string `xml:"ThumbMediaId"`
	} `xml:"Music"`
}

// MessageNews
// @Description: 回复图文消息
type MessageNews struct {
	CommonMessage
	ArticleCount int       `xml:"ArticleCount"`
	Articles     []Article `xml:"Articles>item,omitempty"`
}

// Article
// @Description: 图文
type Article struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicURL      string `xml:"PicUrl,omitempty"`
	URL         string `xml:"Url,omitempty"`
}
