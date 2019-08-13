package main

import (
	"awesomeProject/robot/protocol"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)
func readFile(path string) ([]byte, error) {
	parentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(parentPath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return read(file)
}

func proecss() {
	path := "file_test.go"
	ba, err := readFile(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("The content of '%s':\n%s\n", path, ba)
}

func read(file *os.File)([]byte, error){
	return ioutil.ReadAll(file)
}
func TestFile(t *testing.T) {
	proecss()

}

type Handler struct {
	Method reflect.Method // 方法
	Type   reflect.Type   // 第二个参数类型
}
type SitHandler struct {

}
func (t *SitHandler)Execute(user, req *protocol.MsgHead){
	fmt.Println("req no =", req.AppId)
}

