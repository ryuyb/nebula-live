package bilibili

import (
	"fmt"
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
