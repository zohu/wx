package wxmp

import (
	"fmt"
	"github.com/hhcool/wx"
)

/**
订阅通知
*/

type ParamSubAddTemplate struct {
	Tid       string `json:"tid"`
	KidList   []int  `json:"kidList"`
	SceneDesc string `json:"sceneDesc"`
}
type ResSubAddTemplate struct {
	wx.Response
	PriTmplId string `json:"priTmplId"`
}

// SubAddTemplate
// @Description: 从公共模板库中选用模板，到私有模板库中
// @receiver ctx
// @param p
// @return *ResSubAddTemplate
// @return error
func (ctx *Context) SubAddTemplate(p *ParamSubAddTemplate) (*ResSubAddTemplate, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResSubAddTemplate
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiWxa + "/newtmpl/addtemplate").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&p).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("选用模板失败 %s %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.SubAddTemplate(p)
		}
		return nil, fmt.Errorf("选用模板失败 %s %d-%s", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}
