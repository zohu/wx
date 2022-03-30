package test

import (
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/gtls/utils"
	"github.com/hhcool/wx"
	"github.com/hhcool/wx/wxcpt"
	"github.com/hhcool/wx/wxmp"
	"testing"
)

func TestWxmpMsg(t *testing.T) {
	app, err := wxmp.FindApp("wx8f4971af0e9f0c45")
	if err != nil {
		t.Error(err)
		return
	}
	msg, err1 := app.DecodeMessage(
		&wx.ParamNotify{
			MsgSignature: "4027de4f36110a8a910610f289cc197e60893ba5",
			Timestamp:    "1648610886",
			Nonce:        "1285034511",
		},
		&wxcpt.BizMsg4Recv{
			Tousername: "gh_1cd4920365d4",
			Encrypt:    "MZ0rcHt9GaEFlcrszwfIAwX6Nk/PG59nBQ9b1s13+u1rTaEDY4I/zOCV2xI97oaIR7YnopspaV4Bffnh56xmVyZV8Fu0HpLIqn57J1i8aaR+8kIcNyn+q1eXJqU0YOV7e0TigPN04z0xlQE2OTglBI2GUxsk/XX4upJ0AdpRICyKy7/7pCe0Ys9aC0NGed5GyLQbHdCwfDZnTVdIJc4Zy/xA5JDNtHU1R1JyUcVjj8FdMPHAHZ7UrbgFJDkrfqwwfAnlaRKDZkavP5gCuunStEzobiiTn8TKY+hvO6dd4F3yNNX07SufsHa7odqJhVKOFDSAZY6C6YAW0Nb2wsNHgrMvaqf2SQAYjS4wlcZpk7/02OmVYOrF/jhnbJK1I8afLE3T5eNT845fsEnsJ3sEhvTYJ8C1yPMbakll4BTVvyohRG8gIw+skMnjAc0+a1U7t0nNh9M7fRVt0mjRaofZtz0mFmmciUnBgfaFWQsWWj1OTOA0anQWlVLl/BBA966IbCRBNvWl2f/ri8vsz9hpgRCNL2R+9jsmcO+ztcdBekQ=",
		},
	)
	if err1 != nil {
		t.Error(err1)
		return
	}
	log.Info(utils.StructToString(msg))
}
