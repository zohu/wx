package wxmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zohu/wx"
	"strings"
)

/**
自定义菜单
*/

type MenuType string // 按钮类型
type MenuKey string  // 按钮功能

const (
	MenuTypeView        MenuType = "view"        // 视图按钮
	MenuTypeClick       MenuType = "click"       // 点击按钮
	MenuTypeMiniprogram MenuType = "miniprogram" // 小程序按钮
)
const (
	MenuKeyClick              MenuKey = "click"
	MenuKeyView               MenuKey = "view"
	MenuKeyScancodePush       MenuKey = "scancode_push"        // 扫码推事件
	MenuKeyScancodeWaitmsg    MenuKey = "scancode_waitmsg"     // 扫码推事件且弹出“消息接收中”提示框
	MenuKeyPicSysphoto        MenuKey = "pic_sysphoto"         // 弹出系统拍照发图
	MenuKeyPicPhotoOrAlbum    MenuKey = "pic_photo_or_album"   // 弹出拍照或者相册发图
	MenuKeyPicWeixin          MenuKey = "pic_weixin"           // 弹出微信相册发图器
	MenuKeyLocationSelect     MenuKey = "location_select"      // 弹出地理位置选择器
	MenuKeyMediaId            MenuKey = "media_id"             // 下发消息（除文本消息）
	MenuKeyArticleId          MenuKey = "article_id"           // 微信客户端将会以卡片形式，下发开发者在按钮中填写的图文消息
	MenuKeyArticleViewLimited MenuKey = "article_view_limited" // 类似 view_limited，但不使用 media_id 而使用 article_id
)

type MenuButtonItem struct {
	Type      MenuType         `json:"type,omitempty"`
	Name      string           `json:"name"`
	Key       MenuKey          `json:"key,omitempty"`
	Url       string           `json:"url,omitempty"`
	Appid     string           `json:"appid,omitempty"`
	MediaId   string           `json:"media_id,omitempty"`
	ArticleId string           `json:"article_id,omitempty"`
	Pagepath  string           `json:"pagepath,omitempty"`
	SubButton []MenuButtonItem `json:"sub_button,omitempty"`
}
type Menu struct {
	Button []MenuButtonItem `json:"button"`
}

// MenuDiyMatchRule
// @Description: 性别、国家、省市区、语言 官方已经废除，不再支持
type MenuDiyMatchRule struct {
	TagId              string `json:"tag_id,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
}

type MenuDiy struct {
	Button    []MenuButtonItem `json:"button"`
	Matchrule MenuDiyMatchRule `json:"matchrule"`
}

// MenuAdd
// @Description: 新增菜单
// @receiver ctx
// @param button
// @return *wx.Response
// @return error
func (ctx *Context) MenuAdd(menu *Menu) error {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res wx.Response
	wechat := wx.NewWechat()
	b, _ := json.Marshal(menu)
	b1 := string(b)
	b1 = strings.ReplaceAll(b1, "\\u0026", "&")
	if err := wechat.Post(wx.ApiCgiBin + "/menu/create").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetBody([]byte(b1)).
		BindJSON(&res).
		Do(); err != nil {
		return fmt.Errorf("%s 创建菜单失败（%s）", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuAdd(menu)
		}
		return fmt.Errorf("%s 创建菜单失败（%d-%s）", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return nil
}

type ResMenuQuery struct {
	wx.Response
	IsMenuOpen   int `json:"is_menu_open"`
	SelfmenuInfo struct {
		Button []struct {
			Type      string `json:"type,omitempty"`
			Name      string `json:"name"`
			Key       string `json:"key,omitempty"`
			Url       string `json:"url,omitempty"`
			SubButton struct {
				List []MenuButtonItem `json:"list"`
			} `json:"sub_button,omitempty"`
		} `json:"button"`
	} `json:"selfmenu_info"`
}

// MenuQuery
// @Description: 查询公众号菜单
// @receiver ctx
// @return *ResMenuQuery
// @return error
func (ctx *Context) MenuQuery() (*ResMenuQuery, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResMenuQuery
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiCgiBin + "/get_current_selfmenu_info").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuQuery()
		}
		return nil, fmt.Errorf("%s 查询菜单失败（%d-%s）", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return &res, nil
}

// MenuDelete
// @Description: 删除公众号菜单
// @receiver ctx
// @return error
func (ctx *Context) MenuDelete() error {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiCgiBin + "/menu/delete").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return fmt.Errorf("%s 删除菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuDelete()
		}
		return fmt.Errorf("%s 删除菜单失败（%d-%s）", ctx.Appid(), res.Errcode, res.Errmsg)
	}
	return nil
}

type ResMenuDiyAdd struct {
	wx.Response
	Menuid string `json:"menuid"`
}

// MenuDiyAdd
// @Description: 新增个性化菜单
// @receiver ctx
// @param menu
// @return *ResMenuDiyAdd
// @return error
func (ctx *Context) MenuDiyAdd(menu MenuDiy) (*ResMenuDiyAdd, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResMenuDiyAdd
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiCgiBin + "/menu/addconditional").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&menu).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 新增个性化菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuDiyAdd(menu)
		}
		return nil, fmt.Errorf("%s 新增个性化菜单失败 %s", ctx.Appid(), res.Errmsg)
	}
	return &res, nil
}

// MenuDiyDelete
// @Description: 删除个性化菜单
// @receiver ctx
// @param menuid
// @return *wx.Response
// @return error
func (ctx *Context) MenuDiyDelete(menuid string) (*wx.Response, error) {
	if menuid == "" {
		return nil, errors.New("menuid不能为空")
	}
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res wx.Response
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiCgiBin + "/menu/delconditional").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&struct {
			Menuid string `json:"menuid"`
		}{Menuid: menuid}).
		BindJSON(&res).Do(); err != nil {
		return nil, fmt.Errorf("%s 删除个性化菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuDiyDelete(menuid)
		}
		return nil, fmt.Errorf("%s 删除个性化菜单失败 %s", ctx.Appid(), res.Errmsg)
	}
	return &res, nil
}

type ResMenuDiyTest struct {
	wx.Response
	Button []MenuButtonItem `json:"button"`
}

// MenuDiyTest
// @Description: 测试个性化菜单匹配结果
// @receiver ctx
// @param userID 可以是粉丝的OpenID，也可以是粉丝的微信号
// @return *ResMenuDiyTest
// @return error
func (ctx *Context) MenuDiyTest(userID string) (*ResMenuDiyTest, error) {
	if userID == "" {
		return nil, errors.New("userID不能为空")
	}
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResMenuDiyTest
	wechat := wx.NewWechat()
	if err := wechat.Post(wx.ApiCgiBin + "/menu/trymatch").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(&struct {
			UserID string `json:"user_id"`
		}{UserID: userID}).
		BindJSON(&res).Do(); err != nil {
		return nil, fmt.Errorf("%s 测试个性化菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuDiyTest(userID)
		}
		return nil, fmt.Errorf("%s 测试个性化菜单失败 %s", ctx.Appid(), res.Errmsg)
	}
	return &res, nil
}

type ResMenuQueryAll struct {
	wx.Response
	Menu struct {
		Button []MenuButtonItem `json:"button"`
		Menuid string           `json:"menuid"`
	} `json:"menu"`
	Conditionalmenu struct {
		MenuDiy
		Menuid string `json:"menuid"`
	} `json:"conditionalmenu"`
}

// MenuQueryAll
// @Description: 查询所有菜单（包含个性化菜单）
// @receiver ctx
// @return *ResMenuQueryAll
// @return error
func (ctx *Context) MenuQueryAll() (*ResMenuQueryAll, error) {
	if !ctx.IsMpServe() && !ctx.IsMpSubscribe() {
		return nil, fmt.Errorf("%s 非公众号", ctx.Appid())
	}
	var res ResMenuQueryAll
	wechat := wx.NewWechat()
	if err := wechat.Get(wx.ApiCgiBin + "/menu/get").
		SetQuery(&wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		BindJSON(&res).
		Do(); err != nil {
		return nil, fmt.Errorf("%s 查询菜单失败 %s", ctx.Appid(), err.Error())
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.MenuQueryAll()
		}
		return nil, fmt.Errorf("%s 查询菜单失败 %s", ctx.Appid(), res.Errmsg)
	}
	return &res, nil
}
