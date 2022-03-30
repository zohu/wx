package test

import (
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
	"github.com/hhcool/wx/wxnotify"
	"testing"
)

func TestWxmpMsg(t *testing.T) {
	notify, err := wxnotify.NewNotify(mp.Appid())
	if err != nil {
		t.Error(err)
	}
	msg, err := notify.DecodeMessage(
		&wx.ParamNotify{
			MsgSignature: "b8a00acdc89448f48914379489c86bbe4221550d",
			Timestamp:    "1648635794",
			Nonce:        "129791658",
		},
		&wxcpt.BizMsg4Recv{
			Tousername: "gh_1cd4920365d4",
			Encrypt:    "EFJNdlg9e5qHnTOuEOzcxDMc60U0SFgzI/HDHPUkmWUyibF2tJWSNOOG2QSLl2uVdLX191/FobK+H0IDSECychRt47WKrde+OFcOURrjmDUbFpXAzg19i3H0lQdA+AqgdSzsueySJa688BWRwrOCEtuR3NM08WszCcYuy6yVvZaoVNbtV4CJ3/gtU3SYkyxYLlbUO0Mhii65AytBuWWLDtWl19sK3WB8/8sXWRVSIYjaryF/kxt/Yz75fA90yl+Ds7BCfIWKe90c3IIkvaq75oqLM5MQwlmFiwVHK+4fUBA8Ay01sTqmFIqu92lHHXlTrFdAubR2yMiC1l7Jo6F42JXoeW7Sb4MFK0DwDrjpYo1ldl2gb7La0tTztG2XVFxXJrU2zztadaGFfa/aVN2CdyehhNq8O/lrVGAVNj7mZs1/WI3Sp12rJh20tJkzEAMmLBetPMVokYL+1tSvrsldHc+MZnUV62xE1Mvu1MrDN0oKSWmMQ6M6PoBK+VVybM2J0qzTsM1CX9f7zW0oqMD9at2Ku//MQ+pHXoyS7ynhXAw=",
		},
	)
	if err != nil {
		t.Error(err)
	}
	if msg.MsgId == 0 {
		t.Error("消息ID为空")
	}
}
