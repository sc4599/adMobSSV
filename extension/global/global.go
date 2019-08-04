package global

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"os"
	"time"
)


//节点服务器配置
// 节点服务器数据库使用连接
var (
	Conf *toml.Tree
	RedisCon redis.Conn
	RedisPool *redis.Pool
	//DB *gorm.DB 		// 初始化gorm
	Env = getAppEnv()	// 获取当前环境
)

func init(){
	Conf = NewConf()
	RedisCon = InitRedis()
	RedisPool =InitRedisPool()
	//DB = NewDB()		// 初始化gorm
}


/**
 * 返回单例实例
 * @method New
 */
func NewConf() *toml.Tree {
	fmt.Println("初始化配置文件！！！！！！！！！！！！！！！")
	//config, err := toml.LoadFile("./extension/config/config.toml")
	var filePath string
	env:= getAppEnv()
	switch env {
	case "product":
		filePath = "./extension/config/config.toml"
	case "developer":
		filePath = "conf/config.toml"
	case "-test.run":
		filePath = "E:/gopath/src/awesomeProject/conf/developer_config.toml"
	}
	//if getAppEnv() == "product" {
	//	filePath = "./extension/config/config.toml"
	//}else{
	//	//dir, _ := os.Getwd()
	//	// 由于统一模块，所以测试时可能导致配置路径引用错误问题，所以请根据自己开发环境更改配置文件路径
	//	filePath =  "E:/gopath/src/awesomeProject/conf/developer_config.toml"
	//
	//}
	config, err := toml.LoadFile(filePath)

	if err != nil {
		fmt.Println("TomlError ", err.Error())
	}
	if config == nil {
		fmt.Println("初始化配置文件，config=nil")
	}

	return config
}

//获取程序运行环境
// 根据程序运行路径后缀判断
//如果是 test 就是测试环境
func getAppEnv() string {
	var file string
	if len(os.Args) >2 {
		file = os.Args[2]
	}else{
		file = "developer"
	}
	fmt.Println("当前环境=", file)
	return file
}
func getNodeName() string {
	var noneName string
	if len(os.Args) >1 {
		noneName = os.Args[1]
	}else{
		noneName = ""
	}
	return noneName
}



func InitRedis() redis.Conn {
	fmt.Println("初始化Redis！！！！！！！！！！！！！！！")

	var err error
	conn, err := redis.Dial("tcp",
		Conf.Get("redis.Addr").(string),
		redis.DialConnectTimeout(5*time.Second),
		redis.DialPassword(Conf.Get("redis.Password").(string)),
		redis.DialKeepAlive(3*time.Second),
	)
	if err!=nil{
		checkErr(err)
	}
	dbNum := Conf.Get("redis.DB").(string)
	ok, err:=conn.Do("select", dbNum)
	if ok != "OK"{
		checkErr(err)
	}
	return conn
}

func InitRedisPool() *redis.Pool{
	fmt.Println("初始化Redis连接池！！！！！！！！！！！！！！！")

	// 建立连接池
	connPool := &redis.Pool{
		MaxIdle:     10,//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive:   1000,//最大的激活连接数，表示同时最多有N个连接
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Conf.Get("redis.Addr").(string))
			if err != nil {
				return nil, err
			}
			// 验证密码
			c.Do("auth",Conf.Get("redis.Password"))
			// 选择db
			c.Do("SELECT", Conf.Get("redis.DB").(string))
			return c, nil
		},
	}
	return connPool
}

func checkErr(err error){
	if err != nil{
		panic(err)
	}
}

/**
*设置数据库连接
*@param diver string
 */
func NewDB() *gorm.DB {
	fmt.Println("初始化Mysql连接池！！！！！！！！！！！！！！！")

	driver := Conf.Get("database.dirver").(string)
	configTree := Conf.Get(driver).(*toml.Tree)
	userName := configTree.Get("databaseUserName").(string)
	password := configTree.Get("databasePassword").(string)
	databaseName := configTree.Get("databaseName").(string)
	databaseHost := configTree.Get("databaseHost").(string)
	connect := userName + ":" + password + "@tcp("+databaseHost+")/" + databaseName + "?charset=utf8&parseTime=True&loc=Local"

	fmt.Println(connect)

	DB, err := gorm.Open(driver, connect)

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}

	return DB
}

