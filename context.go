package wx

import (
	"github.com/zohu/zlog"
	"github.com/zohu/zstructs"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

/**
应用实体
*/

const (
	TypeMpServe     = "1"
	TypeMpSubscribe = "2"
	TypeWork        = "3"
	TypeApp         = "4"
	TypeMiniApp     = "5"
)

type App struct {
	Appid          string    `json:"appid"`
	AppSecret      string    `json:"app_secret"`
	Token          string    `json:"token"`
	EncodingAesKey string    `json:"encoding_aes_key"`
	AppType        string    `json:"app_type"`
	AccessToken    string    `json:"access_token"`
	Ticket         string    `json:"ticket"`
	ExpireTime     time.Time `json:"expire_time"`
	Retry          string    `json:"retry"`
}

type Context struct {
	App *App
	sync.Mutex
}

// NewAccessToken
// @Description: 生成当前应用新token
// @receiver c
// @return string
func (c *Context) NewAccessToken() string {
	c.Lock()
	defer c.Unlock()
	if c.App.ExpireTime.Before(time.Now()) {
		c.App.AccessToken = ""
	}
	var token, tk string
	switch c.App.AppType {
	case TypeMpServe:
		token, tk = c.newAccessTokenForMp()
	case TypeMpSubscribe:
		token, tk = c.newAccessTokenForMp()
	case TypeWork:
		token, tk = c.newAccessTokenForWork()
	case TypeApp:
		token, _ = c.newAccessTokenForMp()
	case TypeMiniApp:
		token, _ = c.newAccessTokenForMp()
	default:
		zlog.Warn("NewAccessToken 应用类型错误", zap.String("type", c.App.AppType))
	}
	if token != "" {
		c.App.Retry = "0"
		c.App.AccessToken = token
		c.App.ExpireTime = time.Now().Add(time.Second * 7000)
		c.App.Ticket = tk
		wechat.HSet(RdsAppPrefix+c.Appid(), zstructs.Map(c.App))
	} else {
		wechat.HIncrBy(RdsAppPrefix+c.Appid(), "retry", 1)
	}
	return c.App.AccessToken
}

// GetAccessToken
// @Description: 获取当前应用的token
// @receiver c
// @return string
func (c *Context) GetAccessToken() string {
	newCtx, _ := FindApp(c.Appid())
	c.App = newCtx.App
	if !c.IsExists() {
		return ""
	}
	if c.App.ExpireTime.Before(time.Now()) {
		c.App.AccessToken = ""
	}
	if c.App.AccessToken == "" {
		_ = c.NewAccessToken()
	}
	return c.App.AccessToken
}

func (c *Context) GetTicket() string {
	newCtx, _ := FindApp(c.Appid())
	c.App = newCtx.App
	if !c.IsExists() {
		return ""
	}
	return c.App.Ticket
}

// IsExists
// @Description: 是否存在
// @receiver c
// @return bool
func (c *Context) IsExists() bool {
	if c.App == nil || zstructs.IsZero(c.App) || c.Appid() == "" {
		zlog.Warn("应用不存在", zap.String("appid", c.App.Appid))
		return false
	}
	return true
}

// IsMpServe
// @Description: 是否服务号
// @receiver c
// @return bool
func (c *Context) IsMpServe() bool {
	if c.App.AppType == TypeMpServe {
		return true
	}
	return false
}

// IsMpSubscribe
// @Description: 是否订阅号
// @receiver c
// @return bool
func (c *Context) IsMpSubscribe() bool {
	if c.App.AppType == TypeMpSubscribe {
		return true
	}
	return false
}

// IsWork
// @Description: 是否企业号
// @receiver c
// @return bool
func (c *Context) IsWork() bool {
	if c.App.AppType == TypeWork {
		return true
	}
	return false
}

// IsApp
// @Description: 是否app
// @receiver c
// @return bool
func (c *Context) IsApp() bool {
	if c.App.AppType == TypeApp {
		return true
	}
	return false
}

// IsMiniApp
// @Description: 是否小程序
// @receiver c
// @return bool
func (c *Context) IsMiniApp() bool {
	if c.App.AppType == TypeMiniApp {
		return true
	}
	return false
}

// RetryAccessToken
// @Description: 是否可以刷新token并重试(每个app每2分钟只能重试一次)
// @receiver c
// @param errcode
// @return bool
func (c *Context) RetryAccessToken(errcode int64) bool {
	switch errcode {
	case 40014, 41001, 42001, 42007:
		if wechat.SetNX(RdsAppRetryPrefix+c.Appid(), 1, time.Minute*2).Val() {
			if t := c.NewAccessToken(); t != "" {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Appid
// @Description: 获取标志appid（如果是企微应用的话，返回的是应用ID，不能直接用于接口）
// @receiver c
// @return string
func (c *Context) Appid() string {
	if strings.Contains(c.App.Appid, ":") {
		return strings.Split(c.App.Appid, ":")[1]
	}
	return c.App.Appid
}

// AppidMain
// @Description: 获取主appid（如果是企微应用的话，返回的是企微APPID）
// @receiver c
// @return string
func (c *Context) AppidMain() string {
	if strings.Contains(c.App.Appid, ":") {
		return strings.Split(c.App.Appid, ":")[0]
	}
	return c.App.Appid
}

func (c *Context) AppSecret() string {
	return c.App.AppSecret
}
