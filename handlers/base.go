package handlers

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

// 节点服务器数据库使用连接
var (
	Conf *toml.Tree
)

func init(){
	Conf = NewConf()

}

/**
 * 返回单例实例
 * @method New
 */
func NewConf() *toml.Tree {
	fmt.Println("初始化配置文件！！！！！！！！！！！！！！！")
	//config, err := toml.LoadFile("./extension/config/config.toml")
	filePath := "conf/config.toml"
	config, err := toml.LoadFile(filePath)

	if err != nil {
		fmt.Println("TomlError ", err.Error())
	}
	if config == nil {
		fmt.Println("初始化配置文件，config=nil")
	}

	return config
}

type ApiJson struct {
	Status int64        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}


func ApiResource(status int64, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: objects, Msg: msg}
	return
}
