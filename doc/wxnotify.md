## 微信回调消息
 - 回调消息不再使用传统的中间件形式，将接收与回复结合为msg.reply.encrypt形式，使使用更加灵活，任何框架都可以方便使用。
### 支持的推送回调
- [x] 公众号普通消息推送
- [x] 公众号事件消息推送
- [x] 公众号模板消息送达事件推送
- [x] 公众号订阅通知事件推送
- [x] 公众号H5用户授权信息变更事件推送
- [x] 企业微信客户联系推送
### 接收消息并回复(通用)
```go
func main() {
    r := gin.New()
    r.POST("msg/:appid", RecvHandle)
}
// 以gin举例，其他框架类似
func RecvHandle(c *gin.Context) {
    appid := c.Param("appid")
    p := new(wx.ParamNotify)
    p.MsgSignature, _ = c.GetQuery("msg_signature")
    p.Timestamp, _ = c.GetQuery("timestamp")
    p.Nonce, _ = c.GetQuery("nonce")
    echostr, _ := c.GetQuery("echostr")
    p.Echostr, _ = url.PathUnescape(echostr)
    recv := new(wxnotify.Message)
    if err := c.ShouldBindXML(recv); err != nil {
        log.Error(err)
    }
	// 解密、解析消息
    replyMsg := NotifyHandle(appid,p,recv)
	// 处理返回值
    if res == "" {
        c.String(http.StatusOK, "")
        c.Abort()
        return
    }
    output, _ := xml.MarshalIndent(replyMsg, "  ", "    ")
    _, _ = c.Writer.WriteString(xml.Header)
    _, _ = c.Writer.Write(output)
    c.Status(http.StatusOK)
    c.Abort()
}

func NotifyHandle(appid string,param *wx.ParamNotify,recv *wxcpt.BizMsg4Recv) string {
    notify, _ := wxnotify.NewNotify(appid)
    
    // 接收消息
    msg,err := notify.DecodeMessage(param, recv)
    if err != nil {
        return ""
    }
    
    // 回复文本，非安全模式
    reply := msg.ReplyText("我是文本消息")
    // 安全模式，回复密文
    reply := msg.ReplyText("我是文本消息").Encrypted()
    
    b,_ := xml.Marshal(reply)
    return string(b)
}
```
### 消息加解密，如果需要额外使用的话
参考官方java-sdk改写的go版本，支持xml和json
```go
import "github.com/hhcool/wx/wxcpt"

// 微信公众号
// p *wx.ParamNotify, encpt *wxcpt.BizMsg4Recv
cpt := wxcpt.NewBizMsgCrypt(ctx.App.Token, ctx.App.EncodingAesKey, ctx.Appid())
if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, encpt); err != nil {
    return nil, err
} else {
    event := new(wxmp.Message)
    if err := xml.Unmarshal(cptByte, event); err != nil {
        log.Error(err)
        return nil, err
    }
    return event, nil
}

// 企业微信
if wp, err := wxwork.FindApp(appid); err == nil {
    cpt := wxcpt.NewBizMsgCrypt(wp.App.Token, wp.App.EncodingAesKey, appid)
    if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, encpt); err != nil {
        log.Error(err)
    } else {
        event := new(wxwork.NotifyEvent)
        if err := xml.Unmarshal(cptByte, event); err != nil {
            log.Error(err)
            return ""
        }
        switch event.Event {
        case "change_external_contact": // customer
            w.changeExternalContact(event, wp)
        case "change_external_chat": // 客户群
        case "change_external_tag": // 标签
        }
    }
}
return "ok"
```