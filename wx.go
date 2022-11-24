package wx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zohu/zch"
	"github.com/zohu/zstructs"
	"time"
)

const (
	RdsAppPrefix      = "WX_APP:"
	RdsAppListPrefix  = "WX_APP_LIST"
	RdsAppRetryPrefix = "WX_APP_RETRY:"
)

type Option struct {
	Host     []string
	Password string
	Mode     string
}
type Wechat struct {
	debug    bool
	ctx      context.Context
	RSet     func(k string, v interface{}, t time.Duration) *redis.StatusCmd
	RGet     func(k string) *redis.StringCmd
	RDel     func(k string) *redis.IntCmd
	SetNX    func(k string, v interface{}, t time.Duration) *redis.BoolCmd
	SAdd     func(key string, member ...interface{}) *redis.IntCmd
	SRem     func(key string, member ...interface{}) *redis.IntCmd
	SMembers func(key string) *redis.StringSliceCmd
	HSet     func(key string, value ...interface{}) *redis.IntCmd
	HGetAll  func(key string) *redis.StringStringMapCmd
	HIncrBy  func(key string, field string, integer int64) *redis.IntCmd
	Cancel   context.CancelFunc
}

var wechat = new(Wechat)

func Init(op *Option) {
	rds := zch.NewRds(&redis.UniversalOptions{
		Addrs:    op.Host,
		Password: op.Password,
		DB:       0,
	})
	wechat.RSet = func(k string, v interface{}, t time.Duration) *redis.StatusCmd {
		return rds.Set(context.TODO(), k, v, t)
	}
	wechat.RGet = func(k string) *redis.StringCmd {
		return rds.Get(context.TODO(), k)
	}
	wechat.RDel = func(k string) *redis.IntCmd {
		return rds.Del(context.TODO(), k)
	}
	wechat.SetNX = func(k string, v interface{}, t time.Duration) *redis.BoolCmd {
		return rds.SetNX(context.TODO(), k, v, t)
	}
	wechat.SAdd = func(key string, member ...interface{}) *redis.IntCmd {
		return rds.SAdd(context.TODO(), key, member...)
	}
	wechat.SRem = func(key string, member ...interface{}) *redis.IntCmd {
		return rds.SRem(context.TODO(), key, member...)
	}
	wechat.SMembers = func(key string) *redis.StringSliceCmd {
		return rds.SMembers(context.TODO(), key)
	}
	wechat.HSet = func(key string, value ...interface{}) *redis.IntCmd {
		return rds.HSet(context.TODO(), key, value...)
	}
	wechat.HGetAll = func(key string) *redis.StringStringMapCmd {
		return rds.HGetAll(context.TODO(), key)
	}
	wechat.HIncrBy = func(key string, field string, integer int64) *redis.IntCmd {
		return rds.HIncrBy(context.TODO(), key, field, integer)
	}
	wechat.debug = op.Mode == "debug"
	go wechat.refreshAccessToken()
}

func NewWechat() *Wechat {
	c, cc := context.WithCancel(context.Background())
	wechat.ctx = c
	wechat.Cancel = cc
	return wechat
}

// FindApp
// @Description: 获取APP实例
// @param appid
// @return *Context
// @return error
func FindApp(appid string) (*Context, error) {
	if m := wechat.HGetAll(RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, fmt.Errorf("[wechat:FindApp] 应用不存在 %s", appid)
	} else {
		app := new(App)
		d, _ := json.Marshal(m)
		_ = json.Unmarshal(d, app)
		return &Context{App: app}, nil
	}
}

// PutApp
// @Description: 托管APP
// @param app
// @return error
func PutApp(app App) error {
	app.Retry = "0"
	app.ExpireTime = time.Now()
	wechat.SAdd(RdsAppListPrefix, app.Appid)
	if err := wechat.HSet(RdsAppPrefix+app.Appid, zstructs.Map(app)).Err(); err != nil {
		return fmt.Errorf("PutApp: %s", err.Error())
	}
	if ctx, err := FindApp(app.Appid); err != nil {
		return fmt.Errorf("PutApp find: %s", err.Error())
	} else {
		ctx.NewAccessToken()
	}
	return nil
}

// DelApp
// @Description: 停止APP
// @param appid
func DelApp(appid string) {
	_ = wechat.SRem(RdsAppListPrefix, appid).Err()
	_ = wechat.RDel(RdsAppPrefix + appid).Err()
}

// FindAllApp
// @Description: 查询所有托管的APP
// @return []string
func FindAllApp() []string {
	return wechat.SMembers(RdsAppListPrefix).Val()
}
