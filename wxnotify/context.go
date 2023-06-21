package wxnotify

import "github.com/zohu/wx"

type NotifyContext struct {
	*wx.Context
}

func NewNotify(appid string) (*NotifyContext, *wx.Err) {
	app, err := wx.FindApp(appid)
	if err != nil {
		return nil, err
	}
	return &NotifyContext{app}, nil
}
