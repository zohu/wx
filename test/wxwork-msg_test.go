package test

import (
	"github.com/hhcool/wx/wxcpt"
	"github.com/hhcool/wx/wxwork"
	"testing"
)

func TestWxWorkMsg(t *testing.T) {
	app, err := wxwork.FindApp("ww72ca60e7592549b5")
	if err != nil {
		t.Error(err)
		return
	}
	cpt := wxcpt.NewBizMsgCrypt(app.App.Token, app.App.EncodingAesKey, app.Appid())
	if cptByte, err := cpt.DecryptMsg(
		"ccb6850e23a060d1805cd93a66c115150f2781bd",
		"1648612835",
		"1649324577",
		&wxcpt.BizMsg4Recv{
			Tousername: "ww72ca60e7592549b5",
			Encrypt:    "eaLQ0i7SC6pPod9tj1S9nEgmxj+1eHFTFlM/3FZPyT+94MaWWgBRKafVXKlwAv1q32be72ojy3aVwRYb3ulYCGUuSKEDooZyAxblXl8Thzz3JfmqU0Ss2JewoC2CvhffWqlB5IN+RNO4+/4Vc0oeY9/fopByb2mvoXKA1dUTodnzS4GKF0MjMRlUv0AgQldt+8btNDQ3vLrlTwGuOc36zSzbWZCUsLJGjhjqmE9cwKUozwGs1kd5CspbvMn4hmZgK6iWpPD+yEICB2i1j/qo+IwTyZhk3QXTh47R8sCzdI1PpEkqDaLXe6gYziJAihpoykJEmHK/PWmDrjS1qBMGM/ZJbbgJXu3httwdfsqpYFoAECBHWD3TiZZ45XTRkytmYZ5qQyEksn11crZu3MX9St+NjfaWkLEDNg4dpZF1QOgPLxd5b2te9Glcsc2lL19sQ8dtaJ/maFf7iAXBmGC9/w2BY8HHzDXJrT2Uh1OeC9skzYUPLtV8P4mYK2IYUZfniqCDL85ASyklnWSxeEljFMwALORWhCUicEhEo0PmhrNvuplA+R1PpNDatt/hbTV7VaaSkfLEA825BA6Lks1tIA==",
			Agentid:    "2000003",
		}); err != nil {
		t.Error(err)
	} else {
		if len(cptByte) == 0 {
			t.Error("消息为空")
		}
	}
}
