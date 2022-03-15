package wx

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/gtls/rds"
	"time"
)

const (
	RdsAppPrefix     = "WX_APP:"
	RdsAppListPrefix = "WX_APP_LIST"
)

type Option struct {
	Host     []string
	Password string
	Debug    bool
}
type Wechat struct {
	debug    bool
	ctx      context.Context
	SAdd     func(key string, member ...interface{}) *redis.IntCmd
	SMembers func(key string) *redis.StringSliceCmd
	HSet     func(key string, value ...interface{}) *redis.IntCmd
	HGetAll  func(key string) *redis.StringStringMapCmd
	HIncrBy  func(key string, field string, integer int64) *redis.IntCmd
	Cancel   context.CancelFunc
}

var wechat = new(Wechat)

func Init(op *Option) {
	rds.NewRedis(&rds.Option{Host: op.Host, Password: op.Password})
	wechat.SAdd = rds.Client.SAdd
	wechat.SMembers = rds.Client.SMembers
	wechat.HSet = rds.Client.HSet
	wechat.HGetAll = rds.Client.HGetAll
	wechat.HIncrBy = rds.Client.HIncrBy
	wechat.debug = op.Debug
	go wechat.refreshAccessToken()
}

func NewWechat() *Wechat {
	c, cc := context.WithCancel(context.Background())
	wechat.ctx = c
	wechat.Cancel = cc
	return wechat
}
func FindApp(appid string) *Context {
	if m := wechat.HGetAll(RdsAppPrefix + appid).Val(); len(m) == 0 {
		log.Errorf("[wechat:FindApp] 应用不存在 %s", appid)
		return nil
	} else {
		app := new(App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{App: app}
	}
}
func PutApp(app App) {
	app.Retry = "0"
	app.ExpireTime = time.Now()
	wechat.SAdd(RdsAppListPrefix, app.Appid)
	if err := wechat.HSet(RdsAppPrefix+app.Appid, StructToMap(app)).Err(); err != nil {
		log.Errorf("[wechat:PutApp] %s", err.Error())
	}
	ctx := FindApp(app.Appid)
	ctx.NewAccessToken()
}
