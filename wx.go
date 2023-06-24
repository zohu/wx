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
	RdsAppPrefix      = "wx:01:"
	RdsAppListPrefix  = "wx:02:"
	RdsAppRetryPrefix = "wx:03:"
)

type Option struct {
	Host     []string
	Password string
	Mode     string
	Prefix   string
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
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.Set(context.TODO(), k, v, t)
	}
	wechat.RGet = func(k string) *redis.StringCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.Get(context.TODO(), k)
	}
	wechat.RDel = func(k string) *redis.IntCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.Del(context.TODO(), k)
	}
	wechat.SetNX = func(k string, v interface{}, t time.Duration) *redis.BoolCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.SetNX(context.TODO(), k, v, t)
	}
	wechat.SAdd = func(k string, member ...interface{}) *redis.IntCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.SAdd(context.TODO(), k, member...)
	}
	wechat.SRem = func(k string, member ...interface{}) *redis.IntCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.SRem(context.TODO(), k, member...)
	}
	wechat.SMembers = func(k string) *redis.StringSliceCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.SMembers(context.TODO(), k)
	}
	wechat.HSet = func(k string, value ...interface{}) *redis.IntCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.HSet(context.TODO(), k, value...)
	}
	wechat.HGetAll = func(k string) *redis.StringStringMapCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.HGetAll(context.TODO(), k)
	}
	wechat.HIncrBy = func(k string, field string, integer int64) *redis.IntCmd {
		if op.Prefix != "" {
			k = fmt.Sprintf("%s:%s", op.Prefix, k)
		}
		return rds.HIncrBy(context.TODO(), k, field, integer)
	}
	wechat.debug = op.Mode == "debug"
	go wechat.refreshAccessToken()
}

func NewWechat() *Wechat {
	wechat.ctx, wechat.Cancel = context.WithCancel(context.Background())
	return wechat
}

// FindApp
// @Description: 获取APP实例
// @param appid
// @return *Context
// @return error
func FindApp(appid string) (*Context, *Err) {
	if m := wechat.HGetAll(RdsAppPrefix + appid).Val(); len(m) == 0 {
		return nil, &Err{Appid: appid, Err: "not exits", Desc: "APP不存在"}
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
		return fmt.Errorf("PutApp find: %s", err.Desc)
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
