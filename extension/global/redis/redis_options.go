package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sc4599/fireflygo/logger"
	"awesomeProject/extension/global"
	"awesomeProject/define"
)

/**
 * 设置手机验证码
 */
func SetSmsCode(phone string, code int) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	ok, err := rc.Do("SET", "SMS:"+phone, code, "EX", 60*5)
	if ok != "OK" {
		logger.Info("SmsCode set fail err=", err)
	}
	return
}

/**
 * 获取手机验证码
 */
func GetSmsCode(phone string) (code string) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	code, err := redis.String(rc.Do("GET", "SMS:"+phone))
	if err != nil {
		logger.Info("GetSmsCode Get fail err=", err)
		return
	}
	return
}

// 获取用户错误登录次数
func GetLoginErrorCount(phoneNum string) int {
	rc := global.RedisPool.Get()
	defer rc.Close()
	count, _ := redis.Int(rc.Do("GET", "LoginErrorCount:"+phoneNum))
	return count
}

// 设置用户错误登录次数
func SetLoginErrorCount(phoneNum string, count int) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	ok, err := rc.Do("SET", "LoginErrorCount:"+phoneNum, count, "EX", 3600)
	if ok != "OK" {
		logger.Info("LoginErrorCount set fail err=", err)
	}
}

// 获取用户信息
//func GetUserInfo(userId uint) (str string) {
//	rc := global.RedisPool.Get()
//	defer rc.Close()
//	str, errors := redis.String(rc.Do("GET", "userInfo:"+strconv.Itoa(int(userId))))
//	if str == "" || errors != nil {
//		return
//	}
//	return
//}
//
//// 设置用户信息
//func SetUserInfo(userId string, userInfo []byte) error {
//	rc := global.RedisPool.Get()
//	defer rc.Close()
//	ok, err := rc.Do("SET", "userInfo:"+userId, userInfo, "EX", 3600*12)
//	if ok != "OK" {
//		logger.Info("Set User Info Error", err)
//		return err
//	}
//	return nil
//}

// 获取用户已读公告ID
func GetNoticeRead(userId string) (n string) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	n, _ = redis.String(rc.Do("GET", "NoticeRead:"+userId))
	return
}

// 设置用户已读公告ID
func SetNoticeRead(userId string, r []byte) error {
	rc := global.RedisPool.Get()
	defer rc.Close()
	ok, err := rc.Do("SET", "NoticeRead:"+userId, r)
	if ok != "OK" {
		logger.Error("SetNoticeRead err=", err)
		return err
	}
	return nil
}

// 获取用户当前状态
func GetUserStateByKey(userId string, state string) string {
	rc := global.RedisPool.Get()
	defer rc.Close()
	r, err := redis.String(rc.Do("HGET", "UserState:"+userId, state))
	if err != nil {
		logger.Error(err)
	}
	return r
}

// 获取用户全部状态信息
func GetUserState(userId string) (map[string]string, error) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	options, err := redis.StringMap(rc.Do("hgetall",
		redis.Args{"UserState:" + userId}...))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	return options, nil
}

// 用户状态
func SetUserState(userId string, states map[string]string) error {
	rc := global.RedisPool.Get()
	defer rc.Close()
	args := redis.Args{"UserState:" + userId}
	for k, v := range states {
		args = args.Add(k).AddFlat(v)
	}
	ok, err := rc.Do("hmset", args...)
	if ok != "OK" {
		logger.Error("SetUserState err=", err)
		return err
	}
	return nil
}

// 获取游戏列表
func GetGameList() []string {
	rc := global.RedisPool.Get()
	defer rc.Close()
	r, err := redis.Strings(rc.Do("ZRANGE", "GameList", 0, -1))
	if err != nil {
		logger.Error(err)
	}
	return r
}

// 查询用户session
func GetSession(uid uint64) ([]byte, error) {
	rc := global.RedisPool.Get()
	defer rc.Close()
	session, err := redis.Bytes(rc.Do("get",
		redis.Args{fmt.Sprintf("%s:%d", define.KEY_TOKEN, uid)}...))
	if err != nil {
		return nil, err
	}
	return session, nil
}

func SetSession(userId string, r []byte) error {
	rc := global.RedisPool.Get()
	defer rc.Close()
	ok, err := rc.Do("SET", USER_INFO+userId, r, "EX", 3600*12)
	if ok != "OK" {
		logger.Error("SetSession err=", err)
		return err
	}
	return nil
}

// 获取用户状态 没有数据不返回error
func GetUserStateMap(uid uint64) (map[string]string, error) {
	rc := global.RedisPool.Get()
	options, err := redis.StringMap(rc.Do("hgetall",
		redis.Args{fmt.Sprintf("%s:%d", define.KEY_USER_STATE, uid)}...))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	return options, nil
}


// 风险控制  获取所有风控开关
func GetRiskSwitch(appId int64)(map[string]int64, error){
	rc := global.RedisPool.Get()
	args := redis.Args{fmt.Sprintf("%s:%d", define.KEY_RISK_CONTROL, appId)}
	args = args.Add("Switch")
	options, err := redis.Bytes(rc.Do("HGET",args...))
	if err!=nil && err!= redis.ErrNil {
		return nil, err
	}
	r := make(map[string]int64)
	_ = json.Unmarshal(options, &r)
	return r, nil
}

// 风险控制  获取胜率控制 的胜率
func GetWinRateConf(appId int64)(map[string]interface{}, error){
	rc := global.RedisPool.Get()
	args := redis.Args{fmt.Sprintf("%s:%d", define.KEY_RISK_CONTROL, appId)}
	args = args.Add("WinRate")
	options, err := redis.Bytes(rc.Do("HGET",args...))
	if err!=nil && err!= redis.ErrNil {
		return nil, err
	}
	r := make(map[string]interface{})
	_ = json.Unmarshal(options, &r)
	return r, nil
}
// 风险控制  获取投注盈利比 配置
func GetBetWinRateConf(appId int64)(map[string]interface{}, error){
	rc := global.RedisPool.Get()
	args := redis.Args{fmt.Sprintf("%s:%d", define.KEY_RISK_CONTROL, appId)}
	args = args.Add("BetWinRate")
	options, err := redis.Bytes(rc.Do("HGET",args...))
	fmt.Println("GetBetWinRateConf",string(options))
	if err!=nil && err!= redis.ErrNil {
		return nil, err
	}
	r := make(map[string]interface{})
	err = json.Unmarshal(options, &r)
	if err!=nil{
		return nil, err
	}
	return r, nil
}
// 风险控制  获取投注盈利比 配置
func GetGameStock(appId int64)(map[string]interface{}, error){
	rc := global.RedisPool.Get()
	args := redis.Args{fmt.Sprintf("%s:%d", define.KEY_RISK_CONTROL, appId)}
	args = args.Add("BetWinRate")
	options, err := redis.Bytes(rc.Do("HGET",args...))
	if err!=nil && err!= redis.ErrNil {
		return nil, err
	}
	r := make(map[string]interface{})
	_ = json.Unmarshal(options, &r)
	return r, nil
}