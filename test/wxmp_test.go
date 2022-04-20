package test

import (
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
	"github.com/hhcool/wx/wxmp"
	"github.com/hhcool/wx/wxnotify"
	"os"
	"testing"
)

func TestWxmpMsg(t *testing.T) {
	notify, err := wxnotify.NewNotify(mp.Appid())
	if err != nil {
		t.Error(err)
	}
	msg, err := notify.DecodeMessage(
		&wx.ParamNotify{
			MsgSignature: "8830140fe0589a59412ee68c4036179d02aab83d",
			Timestamp:    "1650350890",
			Nonce:        "1079017418",
		},
		&wxcpt.BizMsg4Recv{
			Tousername: "gh_1cd4920365d4",
			Encrypt:    "Hy0oiflo57b3EgRL38anjPG2EhLSanfsM53QXdtl3/MZg36tS2Gqr2K+5WK5D2kkgW8ZYESt0O2nfY8I36+xFqiLXxP6jpOmiqlyjZOhzpbX38Q5hPiCRZMmsOQDbK7YiCNBcI/F5Q8uFaXDRKrrAQZvYb5Q5zqVoIgVHmy6EA000VEf7beTf0HaOE7smpdBregaAnY5YZZeWDd1Fa6i8+18vN+XbdG40P0bmscG5X7EN9Y4M61Fkil5TnqNVIkbmCXFdOH9/5Ic3iUQ9+eHy+CqxzuKZ5xhP2aWxmvgTPrEWE9yRBjxGznZHJslowBSsbyEaRbT4WbZTK62I/FxmMJwJVxbDfzKIdkDnj5b8UNUIYo21jvWyes9wq7anshs+1vrSDkMJTYxkeTbRTWO8bPxc7v+bnJ5BVOTMpjFYp6XSdYHlmuhLZI3+UiDjPVLLKWdjg7nna+AXeQuPxuHEOGrnfJivUS+o6YTPUNMZoglF3atvTL6Idkb0EH9ERYN",
		},
	)
	if err != nil {
		t.Error(err)
	}
	if msg.FromUserName == "" {
		t.Error("解密失败")
	}
}
func TestH5GetOauth2URL(t *testing.T) {
	uri, err := mp.H5GetOauth2URL("https://beituyun.com?xx=1", wxmp.H5ScopeTypeSnsapiUserinfo, "0001")
	if err != nil {
		t.Error(err)
	}
	if uri != "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx8f4971af0e9f0c45&redirect_uri=https%3A%2F%2Fbeituyun.com%3Fxx%3D1&response_type=code&scope=snsapi_userinfo&state=0001#wechat_redirect" {
		t.Error("H5GetOauth2URL FAIL")
	}
}

func TestMediaTemporaryAdd(t *testing.T) {
	f, _ := os.Open("./11.JPG")
	_, err := mp.MediaTemporaryAdd(wxmp.MediaTypeImage, f, "11.JPG")
	if err != nil {
		t.Error(err)
	}
}
func TestDraftGetAll(t *testing.T) {
	_, err := mp.DraftGetAll(&wxmp.ParamDraftGetAll{
		Page:      1,
		Rows:      20,
		NoContent: 0,
	})
	if err != nil {
		t.Error(err)
	}
}
func TestMediaList(t *testing.T) {
	_, err := mp.MediaList(wxmp.MediaTypeImage, 1, 20)
	if err != nil {
		t.Error(err)
	}
}
func TestIssueList(t *testing.T) {
	_, err := mp.IssueList(0, 1, 20)
	if err != nil {
		t.Error(err)
	}
}
