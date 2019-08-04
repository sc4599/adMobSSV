package demo

import (
	"reflect"
	"testing"
)

func Test1(t *testing.T) {
	s:= Student{Address{},"heilj",13}
	ty:= reflect.TypeOf(s)
	v:=reflect.ValueOf(s)
	for i:=0; i<ty.NumMethod(); i++{
		mName:= ty.Method(i)
		m1:= v.MethodByName(mName.Name)
		m1.Call([]reflect.Value{})
	}
}