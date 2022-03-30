package wxnotify

import "github.com/hhcool/wx"

type NotifyContext struct {
	*wx.Context
}

func NewNotify(appid string) (*NotifyContext, error) {
	app, err := wx.FindApp(appid)
	if err != nil {
		return nil, err
	}
	return &NotifyContext{app}, nil
}
