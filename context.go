package wx

import (
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/gtls/structs"
	"go.uber.org/zap"
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
	token := ""
	switch c.App.AppType {
	case TypeMpServe:
		token = c.newAccessTokenForMp()
	case TypeMpSubscribe:
		token = c.newAccessTokenForMp()
	case TypeWork:
		token = c.newAccessTokenForWork()
	case TypeApp:
		token = c.newAccessTokenForMp()
	case TypeMiniApp:
		token = c.newAccessTokenForMp()
	default:
		log.Error("NewAccessToken 应用类型错误", zap.String("type", c.App.AppType))
	}
	if token != "" {
		c.App.Retry = "0"
		c.App.AccessToken = token
		c.App.ExpireTime = time.Now().Add(time.Second * 7000)
		wechat.HSet(RdsAppPrefix+c.App.Appid, StructToMap(c.App))
	} else {
		wechat.HIncrBy(RdsAppPrefix+c.App.Appid, "retry", 1)
	}
	return c.App.AccessToken
}

// GetAccessToken
// @Description: 获取当前应用的token
// @receiver c
// @return string
func (c *Context) GetAccessToken() string {
	newCtx, _ := FindApp(c.App.Appid)
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

// IsExists
// @Description: 是否存在
// @receiver c
// @return bool
func (c *Context) IsExists() bool {
	if c.App == nil || structs.IsZero(c.App) || c.App.Appid == "" {
		log.Error("应用不存在", zap.String("appid", c.App.Appid))
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
		if wechat.SetNX(RdsAppRetryPrefix+c.App.Appid, 1, time.Minute*2).Val() {
			if t := c.NewAccessToken(); t != "" {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func (c *Context) Appid() string {
	return c.App.Appid
}
