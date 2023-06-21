package test

import (
	"github.com/zohu/wx"
	"github.com/zohu/wx/wxmp"
	"github.com/zohu/wx/wxwork"
)

var (
	mp   *wxmp.Context
	work *wxwork.Context
)

func init() {
	wx.Init(&wx.Option{
		Host:     []string{"localhost:8011"},
		Password: "JCFkQYex4f",
		Mode:     "prod",
	})
	mp, _ = wxmp.FindApp("wx8f4971af0e9f0c45")
	work, _ = wxwork.FindApp("ww72ca60e7592549b5")
	if mp == nil {
		_ = wx.PutApp(wx.App{
			Appid:     "wx8f4971af0e9f0c45",
			AppSecret: "12c289dd81a8a31a7d6f16eb1bf18587",
			AppType:   "1",
		})
		mp, _ = wxmp.FindApp("wx8f4971af0e9f0c45")
	}
	if work == nil {
		_ = wx.PutApp(wx.App{
			Appid:     "ww72ca60e7592549b5",
			AppSecret: "qYpFkZI-p9_i_pWzqF0J5VoJJIxhyYySoW_MYrR6934",
			AppType:   "3",
		})
		work, _ = wxwork.FindApp("ww72ca60e7592549b5")
	}
	if mp == nil || work == nil {
		panic("init test fail")
	}
}
