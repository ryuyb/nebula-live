package bilibili

import (
	"fmt"
	"nebulaLive/pkg/utils"
	"testing"
)

func TestName(t *testing.T) {
	l := &Bilibili{}

	realRoomId, err := l.parseRealRoomId("6")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(realRoomId)
}

func TestName2(t *testing.T) {
	l := &Bilibili{}

	info, err := l.GetInfo("6")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", info)
}

func TestName3(t *testing.T) {
	l := &Bilibili{}

	cookies := utils.ParseCookieStr("")
	stream, err := l.GetStreams("6", 10000, cookies)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", stream)
}
