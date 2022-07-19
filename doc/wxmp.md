- [微信公众号](#%E5%BE%AE%E4%BF%A1%E5%85%AC%E4%BC%97%E5%8F%B7)
    - [自定义菜单](#%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8F%9C%E5%8D%95)
    - [基础消息能力](#%E5%9F%BA%E7%A1%80%E6%B6%88%E6%81%AF%E8%83%BD%E5%8A%9B)
    - [订阅通知](#%E8%AE%A2%E9%98%85%E9%80%9A%E7%9F%A5)
    - [客服消息](#%E5%AE%A2%E6%9C%8D%E6%B6%88%E6%81%AF)
    - [微信网页](#%E5%BE%AE%E4%BF%A1%E7%BD%91%E9%A1%B5)
    - [素材管理](#%E7%B4%A0%E6%9D%90%E7%AE%A1%E7%90%86)
    - [草稿箱](#%E8%8D%89%E7%A8%BF%E7%AE%B1)
    - [发布能力](#%E5%8F%91%E5%B8%83%E8%83%BD%E5%8A%9B)
    - [图文消息留言管理](#%E5%9B%BE%E6%96%87%E6%B6%88%E6%81%AF%E7%95%99%E8%A8%80%E7%AE%A1%E7%90%86)
    - [用户管理](#%E7%94%A8%E6%88%B7%E7%AE%A1%E7%90%86)
    - [账号管理](#%E8%B4%A6%E5%8F%B7%E7%AE%A1%E7%90%86)
    - [数据统计](#%E6%95%B0%E6%8D%AE%E7%BB%9F%E8%AE%A1)
    - [微信卡券](#%E5%BE%AE%E4%BF%A1%E5%8D%A1%E5%88%B8)
    - [微信门店](#%E5%BE%AE%E4%BF%A1%E9%97%A8%E5%BA%97)
    - [微信一物一码](#%E5%BE%AE%E4%BF%A1%E4%B8%80%E7%89%A9%E4%B8%80%E7%A0%81)
    - [微信发票](#%E5%BE%AE%E4%BF%A1%E5%8F%91%E7%A5%A8)
    - [扫服务号二维码打开小程序](#%E6%89%AB%E6%9C%8D%E5%8A%A1%E5%8F%B7%E4%BA%8C%E7%BB%B4%E7%A0%81%E6%89%93%E5%BC%80%E5%B0%8F%E7%A8%8B%E5%BA%8F)


## 微信公众号
### 自定义菜单
- [x] 创建菜单
```
err := app.MenuAdd(&wxmp.Menu{})
```
- [x] 查询菜单
```
menu,err := app.MenuQuery()
```
- [x] 删除菜单
```
err := app.MenuDelete()
```
- [x] 新增个性化菜单
```
_, err := app.MenuDiyAdd(&wxmp.MenuDiy{})
```
- [x] 删除个性化菜单
```
_, err := app.MenuDiyDelete(menuid)
```
- [x] 测试个性化菜单匹配结果
```
// userID 可以是粉丝的OpenID，也可以是粉丝的微信号
_, err := app.MenuDiyTest(userID)
```
- [x] 获取自定义菜单配置
```
// 查询所有菜单，包含个性化菜单
menu, err := app.MenuQueryAll()
```
### 基础消息能力
- [x] 接收普通消息
> 见【[微信回调消息](./wxnotify.md)】
- [x] 接收事件消息
> 见【[微信回调消息](./wxnotify.md)】
- [ ] 被动回复用户消息
> 见【[微信回调消息](./wxnotify.md)】
- [x] 模板消息
```
// 设置所属行业
_, err := app.MsgSetIndustry(&wxmp.ParamMsgSetIndustry{})

// 获取设置的行业信息
res, err := app.MsgGetIndustry()

// 获取模板ID
res, err := app.MsgGetTemplateId(string)

// 获得模板列表
res, err := app.MsgGetTemplateList()

// 删除模板
_, err := app.MsgDelTemplate(string)

// 发送模板消息
_, err := app.MsgSendTemplate(&wxmp.ParamMsgSendTemplate{})

// 是否送达成功事件，参考【[微信回调消息](./wxnotify.md)】
```

- [x] 消息解密
> 见【[微信回调消息](./wxnotify.md)】
```
// 见【接收普通消息】
```
- [x] 公众号一次性订阅消息
```
// 推送订阅模板消息给到授权微信用户
_, err := app.MsgSubscribe(&wxmp.ParamMsgSubscribe{})
```
- [ ] 群发和原创校验
```

```
- [x] 获取公众号的自动回复规则
```
res, err := app.MsgGetAutoReply()
```
### 订阅通知
- [x] 选用模板
```
// 从公共模板库中选用模板，到私有模板库中
res, err := app.SubAddTemplate(&wxmp.ParamSubAddTemplate{})
```
- [x] 删除模板
```
_, err := app.SubDelTemplate(priTmplId string)
```
- [x] 获取公众号类目
```
res, err := app.SubGetCategory()
```
- [ ] 获取模板中的关键词
```
res, err := app.SubGetTemplateKeywords(tid string)
```
- [x] 获取所属类目的公共模板
```
res, err := app.SubGetTemplateTitle(ids string, start int, limit int)
```
- [x] 获取私有模板列表
```
res, err := app.SubGetTemplates()
```
- [x] 发送订阅通知
```
_, err := app.SubBizSend(&wxmp.ParamSubBizSend{})
```
- [x] 订阅通知事件推送
> 见【[微信回调消息](./wxnotify.md)】
### 客服消息
- [x] 获取客服基本信息
```
// 获取客服基本信息
res, err := app.KfList()

// 获取在线客服基本信息
res, err := app.KfOnlineList()
```
- [x] 添加客服账号
```
_, err := app.KfAdd(account string, name string)
```
- [x] 邀请绑定客服账号
```
_, err := app.KfInvite(account string, inviteWx string)
```
- [x] 设置客服信息
```
_, err := app.KfUpdate(account string, name string)
```
- [ ] 上传客服头像
- [ ] 删除客服账号
- [ ] 创建会话
- [ ] 关闭会话
- [ ] 获取客户会话状态
- [ ] 获取客服会话列表
- [ ] 获取未接入会话列表
- [ ] 获取聊天记录
- [ ] 添加顾问
- [ ] 获取顾问信息
- [ ] 修改顾问信息
- [ ] 删除顾问
- [ ] 获取服务号顾问列表
- [ ] 生成顾问二维码
- [ ] 扫顾问二维码后的事件推送
- [ ] 获取顾问聊天记录
- [ ] 设置快捷回复与关注自动回复
- [ ] 获取快捷回复与关注自动回复
- [ ] 设置离线自动回复与敏感词
- [ ] 获取离线自动回复与敏感词
- [ ] 允许微信用户复制小程序页面路径
- [ ] 新建顾问分组
- [ ] 获取顾问分组列表
- [ ] 获取顾问分组信息
- [ ] 分组内添加顾问
- [ ] 分组内删除顾问
- [ ] 获取顾问所在分组
- [ ] 删除顾问分组
- [ ] 为顾问分配客户
- [ ] 为顾问移除客户
- [ ] 获取顾问的客户列表
- [ ] 为客户更换顾问
- [ ] 修改客户昵称
- [ ] 查询客户所属顾问
- [ ] 查询指定顾问和客户的关系
- [ ] 新建标签类型
- [ ] 删除标签类型
- [ ] 为标签添加可选值
- [ ] 获取标签和可选值
- [ ] 为客户设置标签
- [ ] 查询客户标签
- [ ] 根据标签值筛选客户
- [ ] 删除客户标签
- [ ] 设置自定义客户信息
- [ ] 获取自定义客户信息
- [ ] 添加小程序卡片素材
- [ ] 查询小程序卡片素材
- [ ] 删除小程序卡片素材
- [ ] 添加图片素材
- [ ] 查询图片素材
- [ ] 删除图片素材
- [ ] 添加文字素材
- [ ] 查询文字素材
- [ ] 删除文字素材
- [ ] 添加群发任务
- [ ] 获取群发任务列表
- [ ] 获取指定群发任务信息
- [ ] 修改群发任务
- [ ] 取消群发任务
### 微信网页
- [x] 网页授权
```
// ① 获取授权链接
uri,err := app.H5GetOauth2URL(redirectUri string, scope H5ScopeType, state string)

// ② code换用户信息，scope需要和第①步的一致
user, err := app.H5GetUserinfo(code string, scope H5ScopeType)

// 获取jssdk配置
config,err := H5GetJsSdkConfig(uri string)
```
- [x] 用户授权信息变更事件推送
> 见【[微信回调消息](./wxnotify.md)】
### 素材管理
- [x] 新增临时素材
```
// 如果是本地文件
file,_ := os.Open("xxx.jpg")
// 如果是multipart
var mf *multipart.FileHeader
file,_ := mf.Open()

res, err := app.MediaTemporaryAdd(mediaType MediaType, file io.Reader, fileName string)
```
- [ ] 获取临时素材
- [ ] 新增永久素材
- [ ] 获取永久素材
- [ ] 删除永久素材
- [ ] 修改永久图文素材
- [ ] 获取素材总数
- [x] 获取素材列表
- [ ] 上传素材文件
### 草稿箱
- [x] 查询草稿箱列表
### 发布能力
- [x] 获取成功发布列表
### 图文消息留言管理
### 用户管理
- [x] 用户标签管理
```
// 创建标签
res,err := app.UserTagCreate(name string)

// 获取公众号已创建的标签
res,err := app.UserTagQuery()

// 编辑标签
_, err := app.UserTagEdit(id int64, name string)

// 获取标签下粉丝列表
res, err := app.UserTagGetUser(id int64, nextOpenid string)

// 批量为用户打标签
_, err := app.UserTagBatch(openid []string, tagid int64)

// 批量为用户取消标签
_, err := app.UserTagUnBatch(openid []string, tagid int64)

// 获取用户身上的标签列表
res, err := app.UserTagGetFromUser(openid string)
```
- [x] 设置用户备注名
```
_, err := app.UserRemarkUpdate(openid string, remark string)
```
- [x] 获取用户基本信息（含unionID）
``` 
userinfo, err := app.UserFromOpenid(openID)
```
- [x] 获取用户列表
``` 
res, err := app.QueryUserList(nextOpenID)
```
- [x] 获取用户地理位置
> 见【[微信回调消息](./wxnotify.md)】
- [x] 黑名单管理
```
// 获取公众号的黑名单列表
res, err := app.UserGetBlackList(beginOpenid string)

// 拉黑用户
_, err := app.UserBlackListPush(openidList []string)

// 取消拉黑
_, err := app.UserBlackListUnPush(openidList []string)
```
### 账号管理
- [x] 生成带参数的二维码
```
res, err := app.Qrcode(&wxmp.ParamNewQrcode)
```
- [x] 长链接转短链接
```
// 官方已废弃
```
- [x] 短key托管
```
// 获取短Key，ex有效期，可以不传，默认2592000
res, err := app.ShortKey(data string, ex ...int)

// 还原短key
res, err := app.FetchGenShorten(shortKey string)
```
### 数据统计
### 微信卡券
### 微信门店
### 微信小店
### 智能接口
### 微信设备功能
### 微信一物一码
### 微信发票
### 微信非税缴纳
### 扫服务号二维码打开小程序