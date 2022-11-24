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
}
