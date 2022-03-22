package wxwork

import (
	"fmt"
	"github.com/hhcool/wx"
)

type User struct {
	Userid           string  `json:"userid" desc:"成员ID"`
	Name             string  `json:"name" desc:"成员名称"`
	Department       []int64 `json:"department" desc:"所属部门"`
	Order            []int64 `json:"order" desc:"部门内的排序"`
	Position         string  `json:"position" desc:"职务信息"`
	Mobile           string  `json:"mobile" desc:"手机号"`
	Gender           string  `json:"gender" desc:"性别，0未定义1男2女"`
	Email            string  `json:"email" desc:"邮箱"`
	IsLeaderInDept   []int64 `json:"is_leader_in_dept" desc:"是否是上级，0否1是"`
	Avatar           string  `json:"avatar" desc:"头像url"`
	ThumbAvatar      string  `json:"thumb_avatar" desc:"头像缩略图url"`
	Telephone        string  `json:"telephone" desc:"座机"`
	Alias            string  `json:"alias" desc:"别名"`
	Status           int     `json:"status" desc:"状态，1已激活2已禁用4未激活5退出企业"`
	Address          string  `json:"address" desc:"地址"`
	HideMobile       int     `json:"hide_mobile" desc:"是否隐藏电话，0否1是"`
	EnglishName      string  `json:"english_name" desc:"英文名"`
	MainDepartment   int64   `json:"main_department" desc:"主部门"`
	QrCode           string  `json:"qr_code" desc:"员工个人二维码"`
	ExternalPosition string  `json:"external_position" desc:"成员对外职务"`
}

// 获取员工列表

type ParamUserList struct {
	wx.ParamAccessToken
	DepartmentId int64 `json:"department_id" query:"department_id"`
	FetchChild   int   `json:"fetch_child" query:"fetch_child" desc:"0只获取本部门，1递归获取"`
}
type ResponseUserList struct {
	Response
	Userlist []User `json:"userlist"`
}

func (ctx *Context) UserList(p ParamUserList) ([]User, error) {
	if !ctx.IsWork() {
		return nil, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.Appid())
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res ResponseUserList
	err := wechat.Get(wx.ApiWork + "/user/list").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, fmt.Errorf("企业微信：查询用户列表失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return nil, fmt.Errorf("企业微信：查询用户列表失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res.Userlist, nil
}
