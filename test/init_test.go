package test

import "github.com/hhcool/wx"

func init() {
	wx.Init(&wx.Option{
		Host:     []string{"localhost:8011"},
		Password: "JCFkQYex4f",
		Mode:     "debug",
	})
}
