- [开始](#%E5%BC%80%E5%A7%8B)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [多账号管理](#%E5%A4%9A%E8%B4%A6%E5%8F%B7%E7%AE%A1%E7%90%86)
    - [获取实例](#%E8%8E%B7%E5%8F%96%E5%AE%9E%E4%BE%8B)
    - [消息加解密](#%E6%B6%88%E6%81%AF%E5%8A%A0%E8%A7%A3%E5%AF%86)
    - [AccessToken的共享](#accesstoken%E7%9A%84%E5%85%B1%E4%BA%AB)

## 开始
### 初始化
```go
// 本段代码全局唯一即可，可以放到main或者自定义的bootstrap
import "github.com/zohu/wx"
// 初始化微信服务
wx.Init(&wx.Option{
    Host:     Config.Redis.Host,// []string, 一个host列表，支持Client和Cluster；
    Password: Config.Redis.Password,
    Mode:     gin.Mode(), // 非必选，如果非gin框架，可以直接给字符串"debug"/"prod"
})
```
### 多账号管理
下面的例子是用数据库管理多账号，如果只有单账号，也可以只用配置实现。
```go
// 定义数据库表存储多账号
type WxApp struct {
    Id                int    `json:"id" gorm:"primarykey;autoIncrement"`
    Appid             string `json:"appid" gorm:"uniqueIndex;comment:唯一标识"`
    Appsecret         string `json:"appsecret" gorm:"comment:秘钥"`
    AppName           string `json:"app_name" gorm:"comment:app名称"`
    AppParent         string `json:"app_parent" gorm:"comment:绑定的服务号"`
    AppWork           string `json:"app_work" gorm:"comment:绑定的企业号"`
    AppToken          string `json:"app_token" gorm:"comment:消息token"`
    AppEncodingAesKey string `json:"app_encoding_aes_key" gorm:"comment:消息秘钥"`
    AppStatus         int    `json:"app_status" gorm:"default:1;comment:状态，1启用2停用"`
    AppType           string `json:"app_type" gorm:"comment:APP类型：1服务号、2订阅号、3企业号、4app、5小程序、6H5"`
    CreateTime        *Time  `json:"create_time" gorm:"type:datetime;autoCreateTime;comment:创建时间"`
    UpdateTime        *Time  `json:"update_time" gorm:"type:datetime;autoUpdateTime;comment:更新时间"`
}

// 服务启动时，遍历库表进行初始化，业务接口动态的增删参照循环体内逻辑实现
// 如果需要强制覆盖更新，可以不用判断FindApp，直接PutApp即可；
func AppInit() {
    var apps []repo.WxApp
    db.Find(&apps)
    for i := range apps {
        wp := apps[i]
        if wp.AppStatus == 1 {
            if ctx, err := wx.FindApp(wp.Appid); err != nil {
                _ = wx.PutApp(wx.App{
                    Appid:          wp.Appid,
                    AppSecret:      wp.Appsecret,
                    Token:          wp.AppToken,
                    EncodingAesKey: wp.AppEncodingAesKey,
                    AppType:        wp.AppType,
                })
                log.Infof("初始化应用（%s）", wp.Appid)
            } else {
                if ctx.App.ExpireTime.Before(time.Now()) {
                    log.Infof("应用Token过期，刷新Token（%s）", ctx.App.Appid)
                    _ = ctx.GetAccessToken()
                } else {
                    log.Infof("应用正常 %s", ctx.App.Appid)
                }
            }
        } else {
            // 状态非启用的，取消托管
            wx.DelApp(wp.Appid)
        }
    }
}
```

### 获取实例
调用api之前，需要获取对应的app实例
```go
import "github.com/zohu/wx/wxmp"
import "github.com/zohu/wx/wxwork"
// 获取公众号实例
app, err := wxmp.FindApp(appid)
if err != nil {
    return
}
// 获取企业微信实例
app, err := wxwork.FindApp(appid)
if err != nil {
    return
}
// 其他类似……
```
### AccessToken的共享
可以获取对应app的token，用于外部程序共享或自定义接口
```
token := app.GetAccessToken()
```