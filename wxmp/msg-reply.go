package wxmp

import (
	"encoding/xml"
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
	"strconv"
	"time"
)

func encryptMsg(app *wx.App, data []byte, timestamp int64, nonce string) (*wxcpt.XmlBizMsg4Send, error) {
	cpt := wxcpt.NewBizMsgCrypt(app.Token, app.EncodingAesKey, app.Appid)
	return cpt.EncryptXmlMsg(string(data), strconv.FormatInt(timestamp, 10), nonce)
}

// ReplyText
// @Description: 回复文本消息
// @receiver ctx
// @param content
// @return *MessageText
func (msg *Message) ReplyText(content string) *MessageText {
	text := new(MessageText)
	text.Nonce = msg.Nonce
	text.MsgType = MessageTypeText
	text.FromUserName = msg.ToUserName
	text.ToUserName = msg.FromUserName
	text.CreateTime = time.Now().Unix()
	text.Content = wxcpt.CDATA{Value: content}
	return text
}

func (msg *MessageText) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}

// ReplyImage
// @Description: 回复图片消息
// @receiver ctx
// @param mediaID
// @return *MessageImage
func (msg *Message) ReplyImage(mediaID string) *MessageImage {
	image := new(MessageImage)
	image.Nonce = msg.Nonce
	image.MsgType = MessageTypeImage
	image.FromUserName = msg.ToUserName
	image.ToUserName = msg.FromUserName
	image.CreateTime = time.Now().Unix()
	image.Image.MediaID = mediaID
	return image
}

func (msg *MessageImage) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}

// ReplyVoice
// @Description: 回复语音消息
// @receiver ctx
// @param mediaID
// @return *MessageVoice
func (msg *Message) ReplyVoice(mediaID string) *MessageVoice {
	voice := new(MessageVoice)
	voice.Nonce = msg.Nonce
	voice.MsgType = MessageTypeVoice
	voice.FromUserName = msg.ToUserName
	voice.ToUserName = msg.FromUserName
	voice.CreateTime = time.Now().Unix()
	voice.Voice.MediaID = mediaID
	return voice
}

func (msg *MessageVoice) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}

// ReplyVideo
// @Description: 回复视频消息
// @receiver ctx
// @param mediaID
// @param title
// @param description
// @return *MessageVideo
func (msg *Message) ReplyVideo(mediaID, title, description string) *MessageVideo {
	video := new(MessageVideo)
	video.Nonce = msg.Nonce
	video.MsgType = MessageTypeVideo
	video.FromUserName = msg.ToUserName
	video.ToUserName = msg.FromUserName
	video.CreateTime = time.Now().Unix()
	video.Video.MediaID = mediaID
	video.Video.Title = title
	video.Video.Description = description
	return video
}

func (msg *MessageVideo) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}

// ReplyMusic
// @Description: 回复音乐消息
// @receiver ctx
// @param title
// @param description
// @param musicURL
// @param hQMusicURL
// @param thumbMediaID
// @return *MessageMusic
func (msg *Message) ReplyMusic(title, description, musicURL, hQMusicURL, thumbMediaID string) *MessageMusic {
	music := new(MessageMusic)
	music.Nonce = msg.Nonce
	music.MsgType = MessageTypeMusic
	music.FromUserName = msg.ToUserName
	music.ToUserName = msg.FromUserName
	music.CreateTime = time.Now().Unix()
	music.Music.Title = title
	music.Music.Description = description
	music.Music.MusicURL = musicURL
	music.Music.HQMusicURL = hQMusicURL
	music.Music.ThumbMediaID = thumbMediaID
	return music
}

func (msg *MessageMusic) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}

// ReplyNews
// @Description: 回复图文消息
// @receiver ctx
// @param articles
// @return *MessageNews
func (msg *Message) ReplyNews(articles []Article) *MessageNews {
	news := new(MessageNews)
	news.Nonce = msg.Nonce
	news.MsgType = MessageTypeNews
	news.FromUserName = msg.ToUserName
	news.ToUserName = msg.FromUserName
	news.CreateTime = time.Now().Unix()
	news.ArticleCount = len(articles)
	news.Articles = articles
	return news
}

func (msg *MessageNews) Encrypted(app *wx.App) (*wxcpt.XmlBizMsg4Send, error) {
	str, _ := xml.Marshal(msg)
	return encryptMsg(app, str, msg.CreateTime, msg.Nonce)
}
