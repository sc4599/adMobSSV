package models

import "errors"

type Room struct {
	MsgRoute map[int64]func()error
	currentBao int64
}

func NewRoom() *Room{
	r := Room{
		MsgRoute:make(map[int64]func()error),
	}
	r.MsgRoute[10002] = r.Join()
}


func (t *Room)Join()error{

	return nil
}

func (t *Room)SetBao(args...interface{})error{
	if bao,ok:= args[0].(int64);!ok{
		return errors.New("GetParams error")
	}else{
		t.currentBao = bao
	}
	return nil
}