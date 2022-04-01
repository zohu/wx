```shell
go test -v ./test
```
> 通过率 = %
```log
2022-04-01 20:13:16.792 INFO    rds/redis.go:49 初始化缓存库    [Single]
=== RUN   TestWxmpMsg
    wxmp_test.go:34: {"XMLName":{"Space":"","Local":""},"ToUserName":"gh_1cd4920365d4","FromUserName":"o5JV71OdYlLIZNA_4eG15_VpKyMg","CreateTime":1648607685,"MsgType":"text","MsgId":23602383620717229,"Content":"123"}
--- PASS: TestWxmpMsg (0.00s)
=== RUN   TestH5GetOauth2URL
--- PASS: TestH5GetOauth2URL (0.00s)
=== RUN   TestWxWorkMsg
    wxwork_test.go:31: {"XMLName":{"Space":"","Local":""},"ToUserName":"ww72ca60e7592549b5","FromUserName":"sys","CreateTime":1648612835,"MsgType":"event","Event":"change_external_contact","ChangeType":"edit_external_contact","UserID":"1002","ExternalUserID":"wmwaSCCwAABCFazWglrk8b-M3uSJAn3g"}
--- PASS: TestWxWorkMsg (0.00s)
PASS
ok      github.com/hhcool/wx/test       1.172s

```