package test

import (
	"github.com/hhcool/gtls/log"
	"github.com/hhcool/wx/wxmp"
	"testing"
)

func TestWxmpMsg(t *testing.T) {
	_, err := wxmp.FindApp("wx8f4971af0e9f0c45")
	if err != nil {
		log.Error(err.Error())
		return
	}
}
