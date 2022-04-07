package test

import (
	"github.com/hhcool/gtls/utils"
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
	"github.com/hhcool/wx/wxmp"
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
	t.Log(utils.StructToString(msg))

	msg, err = notify.DecodeMessage(
		&wx.ParamNotify{
			MsgSignature: "128e95f07f42bf3ead8ac6bf9b7c51b8ea79f19e",
			Timestamp:    "1649346655",
			Nonce:        "872915664",
		},
		&wxcpt.BizMsg4Recv{
			Tousername: "gh_1cd4920365d4",
			Encrypt:    "/NgOahn80RF551L3HEv/GxRzA2VLr6vyXLjSQMN6A5VyfR0hDFXYYFxsrWRvAn4NprNRh34yh8FBI1lsM7O+ETt3XFqoWb+JrzFEtKOzkMj3ImgL8nxqX5sj48XT2S3ZFve42evrdw/ViT303EKhAulHP7B49L+k6YGneuMjBCHrHtWd5sGVV5s8eI+bXdS8uI86VSpHvEEGkALRydu6JtP5glFta10HCSZdLeZgDcKtZOR2bqX8mBTC6vfdK/07rlFZn8x0MKOHwI+52d7dlp+ttA4f5Rw5mc9XX5FMdJOmEtigDHX7g9PqLOwzaKb7sNmM7SnCbHPpv06iiCynvWC+t28dlR4w2de4tn9GcqEZJOEKLn40oitKWKkyWeJ3",
		},
	)
	if err != nil {
		t.Error(err)
	}
	if msg.MsgId == 0 {
		t.Error("消息ID为空")
	}
	t.Log(utils.StructToString(msg))
}
func TestH5GetOauth2URL(t *testing.T) {
	uri, err := mp.H5GetOauth2URL("https://beituyun.com?xx=1", wxmp.H5ScopeTypeSnsapiUserinfo, "0001")
	if err != nil {
		t.Error(err)
	}
	if uri != "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx8f4971af0e9f0c45&redirect_uri=https%3A%2F%2Fbeituyun.com%3Fxx%3D1&response_type=code&scope=snsapi_userinfo&state=0001#wechat_redirect" {
		t.Error("H5GetOauth2URL FAIL")
	}
}
