package wxwork

import (
	"fmt"
	"github.com/zohu/wx"
)

type ScheduleAttendeesType string

const (
	AttendeesAdd ScheduleAttendeesType = "add_attendees" //新增
	AttendeesDel ScheduleAttendeesType = "del_attendees" //删除
)

type Schedule struct {
	Organizer               string      `json:"organizer" desc:"组织者"`
	StartTime               int64       `json:"start_time" desc:"日程开始时间"`
	EndTime                 int64       `json:"end_time" desc:"日程结束时间"`
	Location                string      `json:"location" desc:"日程地址"`
	AllowActiveJoin         bool        `json:"allow_active_join" desc:"是否允许非参与人主动加入日程，默认为开启"`
	OnlyOrganizerCreateChat int         `json:"only_organizer_create_chat" desc:"是否只允许组织者发起群聊"`
	CalID                   string      `json:"cal_id" desc:"日程所属日历ID"`
	Summary                 string      `json:"summary" desc:"日程标题。0 ~ 128 字符。不填会默认显示为“新建事件“"`
	Description             string      `json:"description" desc:"日程描述 不多于512个字符"`
	Attendees               []Attendees `json:"attendees" desc:"是否只允许组织者发起群聊"`
	Reminders               Reminders   `json:"reminders" desc:"提醒相关信息"`
}

type Attendees struct {
	Userid string `json:"userid" desc:"日程参与者ID"`
}
type Reminders struct {
	IsRemind              int   `json:"is_remind" desc:"是否需要提醒。0-否；1-是"`
	RemindBeforeEventSecs int   `json:"remind_before_event_secs" desc:"日程开始（start_time）前多少秒提醒，当is_remind为1时有效"`
	IsRepeat              int   `json:"is_repeat" desc:"是否重复日程"`
	RepeatType            int   `json:"repeat_type" desc:"重复类型，当is_repeat为1时有效"`
	RepeatUntil           int   `json:"repeat_until" desc:"重复结束时刻，Unix时间戳，当is_repeat为1时有效"`
	IsCustomRepeat        int   `json:"is_custom_repeat" desc:"是否自定义重复"`
	RepeatInterval        int   `json:"repeat_interval" desc:"重复间隔 仅当指定为自定义重复时有效"`
	RepeatDayOfWeek       []int `json:"repeat_day_of_week" desc:"每周周几重复"`
	RepeatDayOfMonth      []int `json:"repeat_day_of_month" desc:"每月哪几天重复"`
	Timezone              int   `json:"timezone" desc:"时区"`
}

type ScheduleResponse struct {
	wx.Response
	ScheduleID string `json:"schedule_id" desc:"日程ID"`
}

// ScheduleCreateInfo 创建日程实体
type ScheduleCreateInfo struct {
	Schedule Schedule `json:"schedule"`
	Agentid  int      `json:"agentid"`
}

// ScheduleCreate 创建日程
func (ctx *Context) ScheduleCreate(p ScheduleCreateInfo) (string, *wx.Err) {
	if !ctx.IsWork() {
		return "", &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	//response实体
	var res ScheduleResponse

	if err := wechat.Post(wx.ApiWorkCgiBin + "/oa/schedule/add").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return "", &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.ScheduleCreate(p)
		}
		return "", &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to add schedule",
			Desc:    "创建日程失败",
		}
	}
	return res.ScheduleID, nil
}

// ScheduleUpdateInfo 更新日程实体
type ScheduleUpdateInfo struct {
	SkipAttendees int      `json:"skip_attendees" desc:"是否不更新参与人。0-否；1-是。默认为0"`
	OpMode        int      `json:"op_mode" desc:"操作模式。是重复日程时有效。0-默认全部修改；1-仅修改此日程；2-修改将来的所有日程"` //
	OpStartTime   int      `json:"op_start_time" desc:"操作起始时间。仅当操作模式是1或2时有效"`
	Schedule      Schedule `json:"schedule" desc:"日程信息"`
}

// ScheduleUpdate 日程更新
func (ctx *Context) ScheduleUpdate(p ScheduleUpdateInfo) (string, *wx.Err) {
	if !ctx.IsWork() {
		return "", &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	//response实体
	var res ScheduleResponse

	if err := wechat.Post(wx.ApiWorkCgiBin + "/oa/schedule/update").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return "", &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.ScheduleUpdate(p)
		}
		return "", &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to update schedule",
			Desc:    "更新日程失败",
		}
	}
	return res.ScheduleID, nil
}

// AttendeesUpdateInfo 更新参与人
type AttendeesUpdateInfo struct {
	ScheduleID string      `json:"schedule_id"`
	Attendees  []Attendees `json:"attendees"`
}

// ScheduleAttendeesUpdate 更新参与人
func (ctx *Context) ScheduleAttendeesUpdate(p AttendeesUpdateInfo, attendeesType ScheduleAttendeesType) *wx.Err {
	if !ctx.IsWork() {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	//response实体
	var res wx.Response

	if err := wechat.Post(wx.ApiWorkCgiBin + fmt.Sprintf("oa/schedule/%s", attendeesType)).
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.ScheduleAttendeesUpdate(p, attendeesType)
		}
		return &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to update schedule attendees",
			Desc:    "更新日程参与者失败",
		}
	}
	return nil
}

type GetScheduleRequest struct {
	ScheduleIDList []string `json:"schedule_id_list" desc:"日程id集合"`
}

// ScheduleDetailResponse 日程详情实体
type ScheduleDetailResponse struct {
	wx.Response
	ScheduleList []ScheduleList `json:"schedule_list" desc:"日程列表"`
}
type AttendeesDetail struct {
	Userid         string `json:"userid"`
	ResponseStatus int    `json:"response_status" desc:"日程参与者的接受状态"`
}
type ExcludeTimeList struct {
	StartTime int `json:"start_time" desc:"不包含的日期时间戳。"`
}
type RemindersDetail struct {
	IsRemind              int               `json:"is_remind" desc:"是否需要提醒。0-否；1-是"`
	IsRepeat              int               `json:"is_repeat" desc:"是否重复日程。0-否；1-是"`
	RemindBeforeEventSecs int               `json:"remind_before_event_secs" desc:"日程开始（start_time）前多少秒提醒，当is_remind为1时有效。"`
	RemindTimeDiffs       []int             `json:"remind_time_diffs" desc:"日程开始（start_time）与提醒时间的差值，当is_remind为1时有效"`
	RepeatType            int               `json:"repeat_type" desc:"重复类型，当is_repeat为1时有效"`
	RepeatUntil           int               `json:"repeat_until" desc:"重复结束时刻"`
	IsCustomRepeat        int               `json:"is_custom_repeat" desc:"是否自定义重复。0-否；1-是"`
	RepeatInterval        int               `json:"repeat_interval" desc:"重复间隔"`
	RepeatDayOfWeek       []int             `json:"repeat_day_of_week" desc:"每周周几重复"`
	RepeatDayOfMonth      []int             `json:"repeat_day_of_month" desc:"每月哪几天重复"`
	Timezone              int               `json:"timezone" desc:"时区"`
	ExcludeTimeList       []ExcludeTimeList `json:"exclude_time_list" desc:"重复日程不包含的日期列表"`
}
type ScheduleList struct {
	ScheduleID              string            `json:"schedule_id"`
	Organizer               string            `json:"organizer" desc:"组织者"`
	Attendees               []AttendeesDetail `json:"attendees" desc:"日程参与者列表。最多支持2000人"`
	Summary                 string            `json:"summary" desc:"日程标题"`
	Description             string            `json:"description" desc:"日程描述"`
	Reminders               RemindersDetail   `json:"reminders" desc:"提醒相关信息"`
	Location                string            `json:"location" desc:"日程地址"`
	CalID                   string            `json:"cal_id" desc:"日程所属日历ID。该日历必须是access_token所对应应用所创建的日历。"`
	StartTime               int               `json:"start_time"`
	EndTime                 int               `json:"end_time"`
	AllowActiveJoin         int               `json:"allow_active_join" desc:"是否允许非参与人主动加入日程"`
	OnlyOrganizerCreateChat int               `json:"only_organizer_create_chat" desc:"是否只允许组织者发起群聊。0-否；1-是"`
	Status                  int               `json:"status" desc:"日程状态。0-正常；1-已取消"`
}

// GetScheduleDetail 获取日程详情
func (ctx *Context) GetScheduleDetail(p GetScheduleRequest) (ScheduleDetailResponse, *wx.Err) {
	//response实体
	res := ScheduleDetailResponse{}
	if !ctx.IsWork() {
		return res, &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()

	if err := wechat.Post(wx.ApiWorkCgiBin + "/oa/schedule/get").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return res, &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.GetScheduleDetail(p)
		}
		return res, &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to get schedule detail",
			Desc:    "获取日程详情失败",
		}
	}
	return res, nil
}

// CancelSchedule 取消日程实体
type CancelSchedule struct {
	ScheduleID  string `json:"schedule_id"`
	OpMode      int    `json:"op_mode"`
	OpStartTime int    `json:"op_start_time"`
}

// CancelSchedule 取消日程
func (ctx *Context) CancelSchedule(p CancelSchedule) *wx.Err {
	if !ctx.IsWork() {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   "not work",
			Desc:  "非企业微信应用",
		}
	}
	wechat := wx.NewWechat()
	//response实体
	var res wx.Response

	if err := wechat.Post(wx.ApiWorkCgiBin + "/oa/schedule/del").
		SetQuery(wx.ParamAccessToken{AccessToken: ctx.GetAccessToken()}).
		SetJSON(p).
		BindJSON(&res).
		Do(); err != nil {
		return &wx.Err{
			Appid: ctx.Appid(),
			Err:   err.Error(),
			Desc:  "请求失败",
		}
	}
	if res.Errcode != 0 {
		if ctx.RetryAccessToken(res.Errcode) {
			return ctx.CancelSchedule(p)
		}
		return &wx.Err{
			Appid:   ctx.Appid(),
			Errcode: res.Errcode,
			Errmsg:  res.Errmsg,
			Err:     "failed to delete schedule",
			Desc:    "取消日程失败",
		}
	}
	return nil
}
