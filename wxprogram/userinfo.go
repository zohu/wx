package wxprogram

import "github.com/zohu/wx"

type ParamGetPhoneNumber struct {
	Code string `json:"code"`
}
type ResGetPhoneNumber struct {
	wx.Response
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int64  `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

// GetPhoneNumber
// @Description: 手机号快速验证
// @receiver ctx
// @param code
// @return *ResGetPhoneNumber
// @return *wx.Err
func (ctx *Context) GetPhoneNumber(code string) (*ResGetPhoneNumber, *wx.Err) {
	if !ctx.IsMiniApp() {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not miniprogram",
			Desc:  "非小程序",
		}
	}
	wechat := wx.NewWechat()
	var res ResGetPhoneNumber
	if err := wechat.Post(wx.ApiWxa + "/business/getuserphonenumber").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&ParamGetPhoneNumber{Code: code}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		return nil, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get phone number",
			Desc:    "获取手机号失败",
		}
	}
	return &res, nil
}
