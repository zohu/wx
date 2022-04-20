package test

import (
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxmp"
	"github.com/hhcool/wx/wxwork"
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
