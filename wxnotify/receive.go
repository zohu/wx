package wxnotify

import (
	"encoding/xml"
	"github.com/zohu/wx"
	"github.com/zohu/wx/wxcpt"
)

func (ctx *NotifyContext) DecodeMessage(p *wx.ParamNotify, encpt *wxcpt.BizMsg4Recv) (*Message, error) {
	cpt := wxcpt.NewBizMsgCrypt(ctx.App.Token, ctx.App.EncodingAesKey, ctx.Appid())
	if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, encpt); err != nil {
		return nil, err
	} else {
		msg := new(Message)
		msg.Nonce = p.Nonce
		msg.ctx = ctx
		if err := xml.Unmarshal(cptByte, msg); err != nil {
			return nil, err
		}
		return msg, nil
	}
}
