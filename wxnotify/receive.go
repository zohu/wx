package wxnotify

import (
	"encoding/xml"
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
)

func (ctx *NotifyContext) DecodeMessage(p *wx.ParamNotify, encpt *wxcpt.BizMsg4Recv) (*Message, error) {
	cpt := wxcpt.NewBizMsgCrypt(ctx.App.Token, ctx.App.EncodingAesKey, ctx.Appid())
	if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, encpt); err != nil {
		return nil, err
	} else {
		msg := new(Message)
		msg.Nonce = p.Nonce
		if err := xml.Unmarshal(cptByte, msg); err != nil {
			return nil, err
		}
		return msg, nil
	}
}
