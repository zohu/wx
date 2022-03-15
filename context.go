package wx

import (
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/gtls/structs"
	"sync"
	"time"
)

/**
应用实体
*/

const (
	TypeMpServe = iota + 1
	TypeMpSubscribe
	TypeWork
	TypeApp
	TypeMiniApp
	TypeH5
)

type App struct {
	Appid          string    `json:"appid"`
	AppSecret      string    `json:"app_secret"`
	Token          string    `json:"token"`
	EncodingAesKey string    `json:"encoding_aes_key"`
	AppType        int       `json:"app_type"`
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
	case TypeH5:
		token = c.newAccessTokenForMp()
	default:
		log.Error("NewAccessToken 应用类型错误")
	}
	if token != "" {
		c.App.Retry = "0"
		c.App.AccessToken = token
		c.App.ExpireTime = time.Now().Add(time.Second * 7200)
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
	newCtx := FindApp(c.App.Appid)
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
		log.Error("应用不存在")
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

// IsH5
// @Description: 是否H5
// @receiver c
// @return bool
func (c *Context) IsH5() bool {
	if c.App.AppType == TypeH5 {
		return true
	}
	return false
}
