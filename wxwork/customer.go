package wxwork

import (
	"encoding/xml"
	"fmt"
	"github.com/hhcool/wx"
	"time"
)

type Customer struct {
	ExternalContact struct {
		ExternalUserid string `json:"external_userid" desc:"客户ID"`
		Name           string `json:"name" desc:"客户名称"`
		Position       string `json:"position" desc:"客户职位（如果也是企业客户）"`
		Avatar         string `json:"avatar" desc:"客户头像"`
		CorpName       string `json:"corp_name" desc:"客户企业简称"`
		CorpFullName   string `json:"corp_full_name" desc:"客户企业"`
		Type           int    `json:"type" desc:"客户类型，1微信2企业微信"`
		Gender         int    `json:"gender" desc:"性别0未知1男2女"`
		Unionid        string `json:"unionid" desc:"唯一标识"`
	} `json:"external_contact"`
	FollowInfo struct {
		Userid         string   `json:"userid" desc:"添加人"`
		Remark         string   `json:"remark" desc:"员工对顾客的备注"`
		Description    string   `json:"description" desc:"员工对顾客的描述"`
		Createtime     int64    `json:"createtime" desc:"添加时间"`
		TagId          []string `json:"tag_id"`
		RemarkCorpName string   `json:"remark_corp_name" desc:"员工给顾客备注的企业名称"`
		RemarkMobiles  []string `json:"remark_mobiles" desc:"员工给顾客备注的手机号"`
		OperUserid     string   `json:"oper_userid" desc:"添加人的userID，可能是员工，也可能是顾客"`
		AddWay         int64    `json:"add_way" desc:"添加来源"`
		AddWayInfo     string   `json:"add_way_info"`
		State          string   `json:"state" desc:"企业指定的，区分哪个联系我"`
	} `json:"follow_info"`
}
type NotifyEvent struct {
	XMLName        xml.Name      `xml:"xml" json:"-"`
	ToUserName     string        `xml:"ToUserName" json:"ToUserName"`
	FromUserName   string        `xml:"FromUserName" json:"FromUserName"`
	CreateTime     time.Duration `xml:"CreateTime" json:"CreateTime"`
	MsgType        string        `xml:"MsgType" json:"MsgType"`
	Event          string        `xml:"Event" json:"Event"`
	ChangeType     string        `xml:"ChangeType" json:"ChangeType,omitempty"`
	UserID         string        `xml:"UserID" json:"UserID,omitempty"`
	ExternalUserID string        `xml:"ExternalUserID" json:"ExternalUserID,omitempty"`
	State          string        `xml:"State" json:"State,omitempty"`
	WelcomeCode    string        `xml:"WelcomeCode" json:"WelcomeCode,omitempty"`
	Source         string        `xml:"Source" json:"Source,omitempty"`
	FailReason     string        `xml:"FailReason" json:"FailReason,omitempty"`
	ChatId         string        `xml:"ChatId" json:"ChatId,omitempty"`
	UpdateDetail   string        `xml:"UpdateDetail" json:"UpdateDetail,omitempty"`
	JoinScene      int64         `xml:"JoinScene" json:"JoinScene,omitempty"`
	QuitScene      int64         `xml:"QuitScene" json:"QuitScene,omitempty"`
	MemChangeCnt   int64         `xml:"MemChangeCnt" json:"MemChangeCnt,omitempty"`
	Id             string        `xml:"Id" json:"Id,omitempty"`
	TagType        string        `xml:"TagType" json:"TagType,omitempty"`
	StrategyId     string        `xml:"StrategyId" json:"StrategyId,omitempty"`
}

// 获取客户列表

type ParamCustomerList struct {
	wx.ParamAccessToken
	Userid string `json:"userid" query:"userid" desc:"员工ID"`
}
type ResponseCustomerList struct {
	Response
	ExternalUserid []string `json:"external_userid"`
}

func (ctx *Context) FindCustomerList(p ParamCustomerList) ([]string, error) {
	if !ctx.IsWork() {
		return nil, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.App.Appid)
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res ResponseCustomerList
	err := wechat.Get(wx.ApiWork + "/externalcontact/list").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return nil, fmt.Errorf("企业微信：查询客户列表失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return nil, fmt.Errorf("企业微信：查询客户列表失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res.ExternalUserid, nil
}

type ParamCustomerDetailWithExternalUserid struct {
	wx.ParamAccessToken
	ExternalUserid string `json:"external_userid"  query:"external_userid"`
}
type ResponseCustomerDetailWithExternalUserid struct {
	Errcode         int             `json:"errcode"`
	Errmsg          string          `json:"errmsg"`
	ExternalContact ExternalContact `json:"external_contact"`
	FollowUser      []FollowUser    `json:"follow_user"`
	NextCursor      string          `json:"next_cursor"`
}
type Text struct {
	Value string `json:"value"`
}
type Web struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
	Title    string `json:"title"`
}
type ExternalAttr struct {
	Type        int         `json:"type"`
	Name        string      `json:"name"`
	Text        Text        `json:"text,omitempty"`
	Web         Web         `json:"web,omitempty"`
	Miniprogram Miniprogram `json:"miniprogram,omitempty"`
}
type ExternalProfile struct {
	ExternalAttr []ExternalAttr `json:"external_attr"`
}
type ExternalContact struct {
	ExternalUserid  string          `json:"external_userid"`
	Name            string          `json:"name"`
	Position        string          `json:"position"`
	Avatar          string          `json:"avatar"`
	CorpName        string          `json:"corp_name"`
	CorpFullName    string          `json:"corp_full_name"`
	Type            int             `json:"type"`
	Gender          int             `json:"gender"`
	Unionid         string          `json:"unionid"`
	ExternalProfile ExternalProfile `json:"external_profile"`
}
type Tags struct {
	GroupName string `json:"group_name"`
	TagName   string `json:"tag_name"`
	TagID     string `json:"tag_id,omitempty"`
	Type      int    `json:"type"`
}
type FollowUser struct {
	Userid         string   `json:"userid"`
	Remark         string   `json:"remark"`
	Description    string   `json:"description"`
	Createtime     int      `json:"createtime"`
	Tags           []Tags   `json:"tags,omitempty"`
	RemarkCorpName string   `json:"remark_corp_name,omitempty"`
	RemarkMobiles  []string `json:"remark_mobiles,omitempty"`
	OperUserid     string   `json:"oper_userid"`
	AddWay         int64    `json:"add_way"`
	State          string   `json:"state,omitempty"`
}

// FindCustomerDetailWithExternalUserid
// @Description: 用客户userid查询客户
// @receiver ctx
// @param p
// @return Customer
// @return error
func (ctx *Context) FindCustomerDetailWithExternalUserid(p ParamCustomerDetailWithExternalUserid) (ResponseCustomerDetailWithExternalUserid, error) {
	var cus ResponseCustomerDetailWithExternalUserid
	if !ctx.IsWork() {
		return cus, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.App.Appid)
	}
	p.AccessToken = ctx.GetAccessToken()
	wechat := wx.NewWechat()
	var res ResponseCustomerDetailWithExternalUserid
	err := wechat.Get(wx.ApiWork + "/externalcontact/get").
		SetQuery(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return cus, fmt.Errorf("企业微信：查询客户详情失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return cus, fmt.Errorf("企业微信：查询客户详情失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	return res, nil
}

// 获取客户详情

type ParamCustomerDetail struct {
	UseridList []string `json:"userid_list" desc:"企业成员id列表，最多100个"`
	Cursor     string   `json:"cursor" desc:"分页游标"`
	Limit      int      `json:"limit"`
}
type ResponseCustomerDetail struct {
	Response
	ExternalContactList []Customer `json:"external_contact_list"`
	NextCursor          string     `json:"next_cursor"`
}

// FindCustomerDetail
// @Description: 根据员工获取客户详情
// @receiver ctx
// @param p
// @param list
// @return []Customer
// @return error
func (ctx *Context) FindCustomerDetail(p ParamCustomerDetail, list []Customer) ([]Customer, error) {
	if !ctx.IsWork() {
		return nil, fmt.Errorf("企业微信：应用 %s 非企业号", ctx.App.Appid)
	}
	if list == nil {
		list = []Customer{}
	}
	p.Limit = 100
	wechat := wx.NewWechat()
	var res ResponseCustomerDetail
	err := wechat.Post(wx.ApiWork + "/externalcontact/batch/get_by_user").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do()
	if err != nil {
		return list, fmt.Errorf("企业微信：查询客户详情失败（%s）", err.Error())
	}
	if res.Errcode != 0 {
		return list, fmt.Errorf("企业微信：查询客户详情失败（%d-%s）", res.Errcode, res.Errmsg)
	}
	for i := range res.ExternalContactList {
		info := res.ExternalContactList[i]
		info.FollowInfo.AddWayInfo = FormatCustomerSource(info.FollowInfo.AddWay)
		list = append(list, info)
	}
	if res.NextCursor != "" {
		p.Cursor = res.NextCursor
		return ctx.FindCustomerDetail(p, list)
	}
	return list, nil
}

func FormatCustomerSource(addWay int64) string {
	switch addWay {
	case 1:
		return "扫描二维码"
	case 2:
		return "搜索手机号"
	case 3:
		return "名片分享"
	case 4:
		return "群聊"
	case 5:
		return "手机通讯录"
	case 6:
		return "微信联系人"
	case 7:
		return "来自微信的添加好友"
	case 8:
		return "第三方客服"
	case 9:
		return "搜索邮箱"
	case 10:
		return "视频号主页添加"
	case 201:
		return "内部成员共享"
	case 202:
		return "管理员分配"
	default:
		return "未知来源"
	}
}
