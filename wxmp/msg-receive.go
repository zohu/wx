package wxmp

import (
	"encoding/xml"
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
)

// DecodeMessage
// @Description: 接收并解析消息
// @receiver ctx
// @param data
// @return error
func (ctx *Context) DecodeMessage(p *wx.ParamNotify, encpt *wxcpt.BizMsg4Recv) (*Message, error) {
	cpt := wxcpt.NewBizMsgCrypt(ctx.App.Token, ctx.App.EncodingAesKey, ctx.Appid())
	if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, encpt); err != nil {
		return nil, err
	} else {
		msg := new(Message)
		if err := xml.Unmarshal(cptByte, msg); err != nil {
			log.Error(err)
			return nil, err
		}
		return msg, nil
	}
}
